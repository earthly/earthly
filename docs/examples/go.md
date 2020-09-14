# Go example

This page will walk you through an example of how to build a hello-world application using Earthly.

First, let's assume that you have written a Hello Wolrd app in `./main.go`:

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

In our case, we will use a `Earthfile` to build it. We will create a target called `build`, and we will copy the necessary files within it (in this case, just `main.go`) and then execute the `go build` command. We will also need a base Docker image that has go pre-installed: `golang:1.13-alpine3.11`.

```Dockerfile
# Earthfile

FROM golang:1.13-alpine3.11

WORKDIR /go-example

build:
    COPY main.go .
    RUN go build -o build/go-example main.go
    SAVE ARTIFACT build/go-example /go-example AS LOCAL build/go-example
```

The `SAVE ARTIFACT` line is necessary to inform Earthly that the resulting file `go-example` is important to us. This will output the file in a local directory at `build/go-example`.

To execute the build, we can run `earth +build`:

```
~/workspace/earthly/examples/go ❯ earth +build
buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
context | --> local context .
+base | --> FROM golang:1.13-alpine3.11
+base | resolve docker.io/library/golang:1.13-alpine3.11@sha256:7d45a6fc9cde63c3bf41651736996fe94a8347e726fe581926fd8c26e244e3b2 0%
+base | resolve docker.io/library/golang:1.13-alpine3.11@sha256:7d45a6fc9cde63c3bf41651736996fe94a8347e726fe581926fd8c26e244e3b2 100%
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

Finally, let's say that we want to build a Docker image for this program. For this, we can add another target, which depends on `build` and uses the built program.

```Dockerfile
# Earthfile

# ...

docker:
    COPY +build/go-example .
    ENTRYPOINT ["/go-example/go-example"]
    SAVE IMAGE go-example:latest
```

We can then run `earth +docker` to build this target:

```
~/workspace/earthly/examples/go ❯ earth +docker     
buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
context | --> local context .
+base | --> FROM golang:1.13-alpine3.11
+base | resolve docker.io/library/golang:1.13-alpine3.11@sha256:7d45a6fc9cde63c3bf41651736996fe94a8347e726fe581926fd8c26e244e3b2 0%
+base | resolve docker.io/library/golang:1.13-alpine3.11@sha256:7d45a6fc9cde63c3bf41651736996fe94a8347e726fe581926fd8c26e244e3b2 100%
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

[![asciicast](https://asciinema.org/a/314637.svg)](https://asciinema.org/a/314637)

## See also

* The [Earthly basics page](../guides/basics.md), which includes an extended Go example
* The [Earthfile reference](../earthfile/earthfile.md)
* The [Earth command reference](../earth-command/earth-command.md)
* More [examples](../examples/examples.md)
