package regproxy

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strconv"

	"github.com/davecgh/go-spew/spew"
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
	err := r.serveDirect(w, req)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), statusInternal)
	}
}

func (r *RegistryProxy) serveDirect(w http.ResponseWriter, req *http.Request) error {

	fmt.Println("ORIG REQ")
	data, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(data))

	u := "http://localhost:8371" + req.URL.Path

	out, err := http.NewRequestWithContext(req.Context(), req.Method, u, req.Body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	for key := range req.Header {
		out.Header.Set(key, req.Header.Get(key))
	}

	fmt.Println("PROXY REQ")
	data, _ = httputil.DumpRequestOut(out, true)
	fmt.Println(string(data))

	res, err := http.DefaultClient.Do(out)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, _ = httputil.DumpResponse(res, true)
	fmt.Println(string(data))

	w.WriteHeader(res.StatusCode)

	for key := range res.Header {
		w.Header().Set(key, res.Header.Get(key))
	}

	_, err = io.Copy(w, res.Body)
	if err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}

	return nil
}

func (r *RegistryProxy) serve(w http.ResponseWriter, req *http.Request) error {

	data, _ := httputil.DumpRequest(req, true)
	fmt.Println(string(data))

	head := map[string]string{}
	for key := range req.Header {
		head[key] = req.Header.Get(key)
	}

	out := &registry.RegistryRequest{
		Path:    req.URL.Path,
		Headers: head,
		Method:  req.Method,
	}

	stream, err := r.cl.Proxy(req.Context(), out)
	if err != nil {
		return fmt.Errorf("failed to send proxy request: %w", err)
	}

	md, err := stream.Header()
	if err != nil {
		return fmt.Errorf("failed to receive stream header: %w", err)
	}

	if md == nil {
		return errors.New("invalid stream metadata")
	}

	spew.Dump(md)

	status, err := strconv.Atoi(mdFirst(md, "status"))
	if err != nil {
		return fmt.Errorf("invalid status: %w", err)
	}
	w.WriteHeader(status)

	md.Delete("status")

	for key := range md {
		w.Header().Set(key, mdFirst(md, key))
	}

	for {
		msg, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return fmt.Errorf("failed to receive stream message: %w", err)
		}
		buf := msg.GetData()
		fmt.Print(string(buf))
		n, err := w.Write(buf)
		if err != nil {
			return fmt.Errorf("failed to write response: %w", err)
		}
		if n != len(buf) {
			return fmt.Errorf("invalid buffer length", err)
		}
	}

	fmt.Println()

	return nil
}

func mdFirst(md metadata.MD, key string) string {
	v := md.Get(key)
	if len(v) == 0 {
		return ""
	}
	return v[0]
}
