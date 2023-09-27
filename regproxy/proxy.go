package regproxy

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	registry "github.com/moby/buildkit/api/services/registry"
	"google.golang.org/grpc/metadata"
)

var statusInternal = http.StatusInternalServerError

func NewRegistryProxy(cl registry.RegistryClient) *RegistryProxy {
	return &RegistryProxy{cl: cl}
}

type RegistryProxy struct {
	cl registry.RegistryClient
}

func (r *RegistryProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := r.serve(w, req)
	if err != nil {
		http.Error(w, err.Error(), statusInternal)
	}
}

func parseHeader(r io.Reader) (status int, header http.Header, err error) {
	sc := bufio.NewScanner(r)
	header = http.Header{}

	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			break
		}
		if strings.Contains(line, "HTTP/") {
			parts := strings.Split(line, " ")
			if len(parts) < 2 {
				err = errors.New("invalid status line")
			}
			status, err = strconv.Atoi(parts[1])
			if err != nil {
				return
			}
			continue
		}
		parts := strings.Split(line, ": ")
		if len(parts) < 2 {
			err = errors.New("invalid header format")
			return
		}
		header.Add(parts[0], parts[1])
	}

	err = sc.Err()

	return
}

func (r *RegistryProxy) serve(w http.ResponseWriter, req *http.Request) error {

	stream, err := r.cl.Proxy(req.Context())
	if err != nil {
		return fmt.Errorf("failed to send proxy request: %w", err)
	}

	rw := registry.NewStreamRW(stream)

	err = req.WriteProxy(rw)
	if err != nil {
		return fmt.Errorf("failed to write request: %w", err)
	}

	err = stream.CloseSend()
	if err != nil {
		return fmt.Errorf("failed to send close: %w", err)
	}

	status, header, err := parseHeader(rw)
	if err != nil {
		return err
	}

	for key, vals := range header {
		for _, val := range vals {
			w.Header().Add(key, val)
		}
	}

	w.WriteHeader(status)

	_, err = io.Copy(w, rw)
	if err != nil {
		return fmt.Errorf("failed to write body: %w", err)
	}

	return nil
}

func mdFirst(md metadata.MD, key string) string {
	v := md.Get(key)
	if len(v) == 0 {
		return ""
	}
	return v[0]
}
