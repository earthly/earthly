# Go example

This page will walk you through an example of how to build a hello-world application using Earthly.

First, let's assume that you have written a Hello World app in `./main.go`:

```go
// main.go
package main

import "fmt"

func main() {
	fmt.Println("hello world")
}
```

In order to build it, you would normally run something like

```bash
go build -o build/go-example main.go
```

In our case, we will use a `Earthfile` to build it. We will create a target called `build`, and we will copy the necessary files within it (in this case, just `main.go`) and then execute the `go build` command. We will also need a base Docker image that has go preinstalled: `golang:1.15-alpine3.13`.

```Dockerfile
# Earthfile
VERSION 0.7
FROM golang:1.15-alpine3.13

WORKDIR /go-example

build:
    COPY main.go .
    RUN go build -o build/go-example main.go
    SAVE ARTIFACT build/go-example /go-example AS LOCAL build/go-example
```

The `SAVE ARTIFACT` line is necessary to inform Earthly that the resulting file `go-example` is important to us. This will output the file in a local directory at `build/go-example`.

To execute the build, we can run `earthly +build`:

```
~/workspace/earthly/examples/go ❯ earthly +build
buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
context | --> local context .
+base | --> FROM golang:1.15-alpine3.13
+base | resolve docker.io/library/golang:1.15-alpine3.13@sha256:7d45a6fc9cde63c3bf41651736996fe94a8347e726fe581926fd8c26e244e3b2 0%
+base | resolve docker.io/library/golang:1.15-alpine3.13@sha256:7d45a6fc9cde63c3bf41651736996fe94a8347e726fe581926fd8c26e244e3b2 100%
+base | --> WORKDIR /go-example
context | transferring .: 0%
context | transferring .: 0%
context | transferring .: 100%
+build | --> FROM ([]) +base
+build | --> COPY [main.go] .
+build | --> RUN [go build -o build/go-example main.go]
+build | Target github.com/earthly/earthly/examples/go:docs-vlad-examples+build built successfully
=========================== SUCCESS ===========================
+build | Artifact github.com/earthly/earthly/examples/go:docs-vlad-examples+build/go-example as local build/go-example
```

And then, we can execute the hello world program:

```
~/workspace/earthly/examples/go ❯ ./build/go-example
hello world
```

Let's say that we want to build a Docker image for this program. For this, we can add another target, which depends on `build` and uses the built program.

```Dockerfile
# Earthfile

# ...

docker:
    COPY +build/go-example .
    ENTRYPOINT ["/go-example/go-example"]
    SAVE IMAGE go-example:latest
```

We can then run `earthly +docker` to build this target:

```
~/workspace/earthly/examples/go ❯ earthly +docker     
buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
context | --> local context .
+base | --> FROM golang:1.15-alpine3.13
+base | resolve docker.io/library/golang:1.15-alpine3.13@sha256:7d45a6fc9cde63c3bf41651736996fe94a8347e726fe581926fd8c26e244e3b2 0%
+base | resolve docker.io/library/golang:1.15-alpine3.13@sha256:7d45a6fc9cde63c3bf41651736996fe94a8347e726fe581926fd8c26e244e3b2 100%
+base | *cached* --> WORKDIR /go-example
+build | *cached* --> FROM ([]) +base
context | transferring .: 0%
context | transferring .: 0%
context | transferring .: 100%
+build | *cached* --> COPY [main.go] .
+build | *cached* --> RUN [go build -o build/go-example main.go]
+build | --> SAVE ARTIFACT build/go-example +build/go-example
+docker | --> COPY ([]) +build/go-example .
+docker | Target github.com/earthly/earthly/examples/go:docs-vlad-examples+docker built successfully
=========================== SUCCESS ===========================
2f7c4d7718b7: Loading layer [==================================================>]  301.3kB/301.3kB
1e3681958479: Loading layer [==================================================>]     153B/153B
468765414030: Loading layer [==================================================>]  126.9MB/126.9MB
88bb86abec1d: Loading layer [==================================================>]     127B/127B
6309ef8dcbb1: Loading layer [==================================================>]     898B/898B
66d3e336d237: Loading layer [==================================================>]     914B/914B
160251c27a86: Loading layer [==================================================>]  1.048MB/1.048MB
Loaded image: go-example:latest
+docker | Image github.com/earthly/earthly/examples/go:docs-vlad-examples+docker as go-example:latest
+build | Artifact github.com/earthly/earthly/examples/go:docs-vlad-examples+build/go-example as local build/go-example
```

And then we can run the built image like so:

```
~/workspace/earthly/examples/go ❯ docker run --rm go-example:latest
hello world
```

Not only you can run your program with Earthly, but also your unit and integration tests. 

To execute the unit-tests, we can run `earthly -P +unit-test`:

```
           buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
golang:1.15-alpine3.13 | --> Load metadata linux/amd64
               +base | --> FROM golang:1.15-alpine3.13
             context | --> local context .
               +base | [██████████] resolve docker.io/library/golang:1.15-alpine3.13@sha256:330f31a4415d97bb64f244d5f4d838bea7a7ee1ab5a1a0bac49e7973c57cbb88 ... 100%
             context | transferred 3 file(s) for context . (2.4 MB, 9 file/dir stats)
               +base | --> WORKDIR /go-example
               +deps | --> COPY go.mod go.sum ./
               +deps | --> RUN go mod download
               +deps | --> SAVE ARTIFACT go.sum +deps/go.sum AS LOCAL go.sum
               +deps | --> SAVE ARTIFACT go.mod +deps/go.mod AS LOCAL go.mod
          +unit-test | --> COPY main.go .
          +unit-test | --> COPY main_test.go .
          +unit-test | --> RUN CGO_ENABLED=0 go test github.com/earthly/earthly/examples/go
              output | --> exporting outputs
              output | [██████████] copying files ... 100%
================================ SUCCESS [main] ================================
               +deps | Artifact github.com/earthly/earthly/examples/go:go-integration-test-example+deps/go.mod as local go.mod
               +deps | Artifact github.com/earthly/earthly/examples/go:go-integration-test-example+deps/go.sum as local go.sum

``` 

To execute the integration-tests, we can run `earthly -P +integration-test`:

```
           buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
golang:1.15-alpine3.13 | --> Load metadata linux/amd64
               +base | --> FROM golang:1.15-alpine3.13
             context | --> local context .
               +base | [██████████] resolve docker.io/library/golang:1.15-alpine3.13@sha256:330f31a4415d97bb64f244d5f4d838bea7a7ee1ab5a1a0bac49e7973c57cbb88 ... 100%
             context | transferred 1 file(s) for context . (6.8 kB, 9 file/dir stats)
               +base | --> WORKDIR /go-example
               +deps | --> COPY go.mod go.sum ./
               +deps | --> RUN go mod download
   +integration-test | --> COPY main.go .
   +integration-test | --> COPY main_integration_test.go .
   +integration-test | --> COPY docker-compose.yml ./
   +integration-test | --> WITH DOCKER (install deps)
   +integration-test | --> WITH DOCKER (docker-compose config)
    redis:6.0-alpine | --> Load metadata linux/amd64
    redis:6.0-alpine | --> DOCKER PULL redis:6.0-alpine
    redis:6.0-alpine | [██████████] resolve docker.io/library/redis:6.0-alpine@sha256:61f3e955fbef87ea07d7409a48a48b069579e32f37d2f310526017d68e9983b7 ... 100%
               +deps | --> SAVE ARTIFACT go.sum +deps/go.sum AS LOCAL go.sum
               +deps | --> SAVE ARTIFACT go.mod +deps/go.mod AS LOCAL go.mod
             context | transferred 1 file(s) for context /var/folders/5f/jkczhmh52g71v8_q34kt2wm80000gn/T/earthly-docker-load330575103 (10 MB, 1 file/dir stats)
   +integration-test | --> WITH DOCKER RUN --privileged CGO_ENABLED=0 go test github.com/earthly/earthly/examples/go
   +integration-test | Loading images...
   +integration-test | Loaded image: redis:6.0-alpine
   +integration-test | ...done
   +integration-test | Creating network "go-example_default" with the default driver
   +integration-test | Creating local-redis ... done
   +integration-test | Creating local-redis ... done
             ongoing | ok  egratgithub.com/earthly/earthly/examples/go  0.006s
   +integration-test | Stopping local-redis ... done
   +integration-test | Removing local-redis ... done
   +integration-test | Removing network go-example_default
              output | --> exporting outputs
              output | [██████████] copying files ... 100%
================================ SUCCESS [main] ================================
               +deps | Artifact github.com/earthly/earthly/examples/go:go-integration-test-example+deps/go.mod as local go.mod
               +deps | Artifact github.com/earthly/earthly/examples/go:go-integration-test-example+deps/go.sum as local go.sum
```

Finally, to run the build, unit test, integration test and docker image just run `earthly -P +all`:

```
          buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
golang:1.15-alpine3.13 | --> Load metadata linux/amd64
             context | --> local context .
               +base | --> FROM golang:1.15-alpine3.13
               +base | [██████████] resolve docker.io/library/golang:1.15-alpine3.13@sha256:330f31a4415d97bb64f244d5f4d838bea7a7ee1ab5a1a0bac49e7973c57cbb88 ... 100%
             context | transferred 1 file(s) for context . (9.6 kB, 9 file/dir stats)
               +base | --> WORKDIR /go-example
               +deps | --> COPY go.mod go.sum ./
               +deps | --> RUN go mod download
   +integration-test | --> COPY main.go .
   +integration-test | --> COPY main_integration_test.go .
   +integration-test | --> COPY docker-compose.yml ./
   +integration-test | --> WITH DOCKER (install deps)
   +integration-test | --> WITH DOCKER (docker-compose config)
    redis:6.0-alpine | --> Load metadata linux/amd64
    redis:6.0-alpine | --> DOCKER PULL redis:6.0-alpine
    redis:6.0-alpine | [██████████] resolve docker.io/library/redis:6.0-alpine@sha256:61f3e955fbef87ea07d7409a48a48b069579e32f37d2f310526017d68e9983b7 ... 100%
               +deps | --> SAVE ARTIFACT go.sum +deps/go.sum AS LOCAL go.sum
              +build | --> RUN go build -o build/go-example main.go
              +build | --> SAVE ARTIFACT build/go-example +build/go-example AS LOCAL build/go-example
             +docker | --> COPY +build/go-example ./
               +deps | --> SAVE ARTIFACT go.mod +deps/go.mod AS LOCAL go.mod
   +integration-test | --> WITH DOCKER RUN --privileged CGO_ENABLED=0 go test github.com/earthly/earthly/examples/go
          +unit-test | --> COPY main_test.go .
          +unit-test | --> RUN CGO_ENABLED=0 go test github.com/earthly/earthly/examples/go
              output | --> exporting outputs
              output | [██████████] exporting layers ... 100%
              output | [██████████] exporting manifest sha256:73cba3d853028a7fb74b936ce78b1adaf510b9d8ca57da67e5120bd38283b685 ... 100%
              output | [██████████] exporting config sha256:9fef849e6fa9477fed2f481b7ac07d419c5ddd63c8179cac4e3401be174f4025 ... 100%
              output | [██████████] copying files ... 100%
              output | [██████████] transferring docker.io/earthly/examples:go ... 100%
================================ SUCCESS [main] ================================
              +build | Artifact github.com/earthly/earthly/examples/go:go-integration-test-example+build/go-example as local build/go-example
               +deps | Artifact github.com/earthly/earthly/examples/go:go-integration-test-example+deps/go.mod as local go.mod
               +deps | Artifact github.com/earthly/earthly/examples/go:go-integration-test-example+deps/go.sum as local go.sum
             +docker | Image github.com/earthly/earthly/examples/go:go-integration-test-example+docker as earthly/examples:go
             +docker | Did not push earthly/examples:go. Use earthly --push to enable pushing
```

[![asciicast](https://asciinema.org/a/314637.svg)](https://asciinema.org/a/314637)
