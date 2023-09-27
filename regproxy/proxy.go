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

// NewRegistryProxy creates and returns a new RegistryProxy that streams image
// data from the BK embedded Docker registry.
func NewRegistryProxy(cl registry.RegistryClient) *RegistryProxy {
	return &RegistryProxy{cl: cl}
}

// RegistryProxy uses a gRPC stream to translate incoming Docker image requests
// into a gRPC byte stream and back out into a valid HTTP response. The data is
// streamed over the gRPC connection rather than buffered as the images can be
// rather large.
type RegistryProxy struct {
	cl registry.RegistryClient
}

// ServeHTTP implements http.Handler.
func (r *RegistryProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	err := r.serve(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// parseHeader parses an HTTP response header and extracts the response status &
// header values. This function does not close the reader as it will be used
// further on to stream the body data.
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

// serve the HTTP request by writing it to BK via a streaming gRPC request. The
// request will be reconstituted on the other end and forwarded to the embedded
// registry and back again.
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
