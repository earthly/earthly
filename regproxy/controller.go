package regproxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	conslog "github.com/earthly/earthly/conslogging"
	"github.com/earthly/earthly/util/containerutil"
	"github.com/earthly/earthly/util/stringutil"
	registry "github.com/moby/buildkit/api/services/registry"
	"github.com/pkg/errors"
)

const (
	darwinContainerPrefix = "earthly-darwin-proxy"
	darwinContainerMaxAge = 5 * time.Hour
)

// Controller handles the management of the registry proxy. This may also
// include the Darwin proxy used to enable Docker Desktop setups.
type Controller struct {
	registryClient    registry.RegistryClient
	containerFrontend containerutil.ContainerFrontend
	darwinProxy       bool
	darwinProxyImage  string
	darwinProxyWait   time.Duration
	cons              conslog.ConsoleLogger
}

// NewController creates and returns a new registry proxy controller.
func NewController(
	registryClient registry.RegistryClient,
	containerFrontend containerutil.ContainerFrontend,
	darwinProxy bool,
	darwinProxyImage string,
	darwinProxyWait time.Duration,
	cons conslog.ConsoleLogger,
) *Controller {
	return &Controller{
		registryClient:    registryClient,
		containerFrontend: containerFrontend,
		darwinProxy:       darwinProxy,
		darwinProxyImage:  darwinProxyImage,
		darwinProxyWait:   darwinProxyWait,
		cons:              cons,
	}
}

// Start the proxy and create any support containers.
func (c *Controller) Start(ctx context.Context) (string, func(), error) {
	addr := "127.0.0.1:0"

	ln, err := (&net.ListenConfig{}).Listen(ctx, "tcp", addr)
	if err != nil {
		return "", nil, errors.Wrap(err, "failed to create proxy listener")
	}

	p := newRegistryProxy(ln, c.registryClient)
	go p.serve(ctx)

	// Find the assigned port.
	registryPort := ln.Addr().(*net.TCPAddr).Port
	addr = fmt.Sprintf("127.0.0.1:%d", registryPort)

	c.cons.VerbosePrintf("Starting registry proxy on %s", addr)

	doneCh := make(chan struct{})

	go func() {
		for err := range p.err() {
			if err != nil && !errors.Is(err, context.Canceled) {
				c.cons.VerbosePrintf("Failed to serve registry proxy: %v", err)
			}
		}
		doneCh <- struct{}{}
	}()

	closers := []func(ctx context.Context){
		func(ctx context.Context) {
			p.close()
			select {
			case <-ctx.Done():
			case <-doneCh:
			}
		},
	}

	if c.darwinProxy {
		containerName := fmt.Sprintf("%s-%s", darwinContainerPrefix, stringutil.RandomAlphanumeric(6))
		stopFn := func(ctx context.Context) {
			err := c.stopDarwinProxy(ctx, containerName, true)
			if err != nil {
				c.cons.VerbosePrintf("Failed to stop registry proxy support container: %v", err)
			}
		}
		port, err := c.startDarwinProxy(ctx, containerName, registryPort)
		if err != nil {
			stopFn(ctx)
			return "", nil, errors.Wrap(err, "failed to start Darwin support container")
		}
		addr = fmt.Sprintf("127.0.0.1:%d", port)
		c.cons.VerbosePrintf("Starting Darwin proxy on %s", addr)
		closers = append(closers, stopFn)
	}

	return addr, func() {
		for _, closer := range closers {
			closer(ctx)
		}
	}, nil
}

// startDarwinProxy: Since Docker Desktop (Mac) containers run in a VM, a
// special host name, host.docker.internal, is made available to access the host
// machine. Docker can only pull insecurely from localhost, so we use a socat
// container to proxy localhost:<port> request back out to the local registry
// proxy created above.
func (c *Controller) startDarwinProxy(ctx context.Context, containerName string, registryPort int) (int, error) {
	go func() {
		err := c.stopOldDarwinProxies(ctx)
		if err != nil {
			c.cons.VerbosePrintf("Failed to stop old Darwin proxy support container: %s", err)
		}
	}()

	containerPort, err := acquireFreePort(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "failed to acquire free port")
	}

	runCfg := containerutil.ContainerRun{
		NameOrID: containerName,
		ImageRef: c.darwinProxyImage,
		Ports: []containerutil.Port{
			{
				IP:            "127.0.0.1",
				HostPort:      containerPort, // Bind to available port
				ContainerPort: 80,
				Protocol:      containerutil.ProtocolTCP,
			},
		},
		ContainerArgs: []string{
			"tcp-listen:80,fork,reuseaddr",
			fmt.Sprintf("tcp:host.docker.internal:%d", registryPort),
		},
	}

	err = c.containerFrontend.ContainerRun(ctx, runCfg)
	if err != nil {
		return 0, errors.Wrap(err, "failed to start support container")
	}

	childCtx, cancel := context.WithTimeout(ctx, c.darwinProxyWait)
	defer cancel()

	// Wait for the proxy chain to resolve to the BK registry. The /v2/ path
	// will return a 200 when ready.
	for {
		res, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/v2/", containerPort))
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
		if err == nil && res != nil && res.StatusCode == http.StatusOK {
			break
		}
		select {
		case <-childCtx.Done():
			return 0, childCtx.Err()
		case <-time.After(time.Second):
			continue
		}
	}

	return containerPort, nil
}

func (c *Controller) stopOldDarwinProxies(ctx context.Context) error {
	containers, err := c.containerFrontend.ContainerList(ctx)
	if err != nil {
		return err
	}
	for _, container := range containers {
		if strings.HasPrefix(container.Name, darwinContainerPrefix) &&
			time.Since(container.Created) > darwinContainerMaxAge {
			err = c.stopDarwinProxy(ctx, container.Name, false)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Controller) stopDarwinProxy(_ context.Context, containerName string, checkExists bool) error {
	// Ignore parent context cancellations as to prevent orphaned containers.
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if checkExists {
		infos, err := c.containerFrontend.ContainerInfo(ctx, containerName)
		if err != nil {
			return err
		}
		if info, ok := infos[containerName]; !ok || info.Status == containerutil.StatusMissing {
			return nil
		}
	}
	err := c.containerFrontend.ContainerRemove(ctx, true, containerName)
	if err != nil {
		return errors.Wrap(err, "failed to stop support container")
	}
	return nil
}

func acquireFreePort(ctx context.Context) (int, error) {
	addr := "127.0.0.1:0"

	ln, err := (&net.ListenConfig{}).Listen(ctx, "tcp", addr)
	if err != nil {
		return 0, errors.Wrap(err, "failed to listen on open port")
	}
	defer ln.Close() // Immediately close the listener

	return ln.Addr().(*net.TCPAddr).Port, nil
}
