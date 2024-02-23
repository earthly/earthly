
To copy the files for [this example ( Part 3 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/go/part3) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/go:main+part3/part3 ./part3
```
Examples in [Python](#more-examples), [JavaScript](#more-examples) and [Java](#more-examples) are at the bottom of this page.

## Dependencies
Now let's imagine that we want to add some dependencies to our app. In Go, that means adding `go.mod` and `go.sum`. 

`./go.mod`

```go.mod
module github.com/earthly/earthly/examples/go

go 1.13

require github.com/sirupsen/logrus v1.5.0
```

`./go.sum` (empty)

```go.sum
```

The code of the app might look like this

`./main.go`

```go
package main

import "github.com/sirupsen/logrus"

func main() {
	logrus.Info("hello world")
}
```

Now we can update our Earthfile to copy in the `go.mod` and `go.sum`.

`./Earthfile`

```Dockerfile
VERSION 0.8
FROM golang:1.15-alpine3.13
WORKDIR /go-workdir

build:
    COPY go.mod go.sum .
    COPY main.go .
    RUN go build -o output/example main.go
    SAVE ARTIFACT output/example AS LOCAL local-output/go-example

docker:
    COPY +build/example .
    ENTRYPOINT ["/go-workdir/example"]
    SAVE IMAGE go-example:latest
```
This works, but it is inefficient because we have not made proper use of caching. In the current setup, when a file changes, the corresponding `COPY` command is re-executed without cache, causing all commands after it to also re-execute without cache.

### Caching

If, however, we could first download the dependencies and only afterwards copy and build the code, then the cache would be reused every time we changed the code.

`./Earthfile`

```Dockerfile
VERSION 0.8
FROM golang:1.15-alpine3.13
WORKDIR /go-workdir

build:
    # Download deps before copying code.
    COPY go.mod go.sum .
    RUN go mod download
    # Copy and build code.
    COPY main.go .
    RUN go build -o output/example main.go
    SAVE ARTIFACT output/example AS LOCAL local-output/go-example

docker:
    COPY +build/example .
    ENTRYPOINT ["/go-workdir/example"]
    SAVE IMAGE go-example:latest
```

For a primer into Dockerfile layer caching see [this article](https://pythonspeed.com/articles/docker-caching-model/). The same principles apply to Earthfiles.

## Reusing Dependencies

In some cases, the dependencies might be used in more than one build target. For this use case, we might want to separate dependency downloading and reuse it. For this reason, let's consider breaking this out into a separate target called `+deps`. We can then inherit from `+deps` by using the command `FROM +deps`.

`./Earthfile`

```Dockerfile
VERSION 0.8
FROM golang:1.15-alpine3.13
WORKDIR /go-workdir

deps:
    COPY go.mod go.sum ./
    RUN go mod download
    # Output these back in case go mod download changes them.
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

build:
    FROM +deps
    COPY main.go .
    RUN go build -o output/example main.go
    SAVE ARTIFACT output/example AS LOCAL local-output/go-example

docker:
    COPY +build/example .
    ENTRYPOINT ["/go-workdir/example"]
    SAVE IMAGE go-example:latest
```

## More Examples

<details open>
<summary>JavaScript</summary>

To copy the files for [this example ( Part 3 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/js/part3) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/js:main+part3/part3 ./part3
```

Note that in our case, only the JavaScript version has an example where `FROM +deps` is used in more than one place: both in `build` and in `docker`. Nevertheless, all versions show how dependencies may be separated.

`./Earthfile`

```Dockerfile
VERSION 0.8
FROM node:13.10.1-alpine3.11
WORKDIR /js-example

deps:
    COPY package.json ./
    COPY package-lock.json ./
    RUN npm install
    # Output these back in case npm install changes them.
    SAVE ARTIFACT package.json AS LOCAL ./package.json
    SAVE ARTIFACT package-lock.json AS LOCAL ./package-lock.json

build:
    FROM +deps
    COPY src src
    RUN mkdir -p ./dist && cp ./src/index.html ./dist/
    RUN npx webpack
    SAVE ARTIFACT dist /dist AS LOCAL dist

docker:
    FROM +deps
    COPY +build/dist ./dist
    EXPOSE 8080
    ENTRYPOINT ["/js-example/node_modules/http-server/bin/http-server", "./dist"]
    SAVE IMAGE js-example:latest
```

</details>


<details open>
<summary>Java</summary>

To copy the files for [this example ( Part 3 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/java/part3) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/java:main+part3/part3 ./part3
```

`./Earthfile`

```Dockerfile
VERSION 0.8
FROM openjdk:8-jdk-alpine
RUN apk add --update --no-cache gradle
WORKDIR /java-example

deps:
    COPY build.gradle ./
    RUN gradle build

build:
    FROM +deps
    COPY src src
    RUN gradle build
    RUN gradle install
    SAVE ARTIFACT build/install/java-example/bin AS LOCAL build/bin
    SAVE ARTIFACT build/install/java-example/lib AS LOCAL build/lib

docker:
    COPY +build/bin bin
    COPY +build/lib lib
    ENTRYPOINT ["/java-example/bin/java-example"]
    SAVE IMAGE java-example:latest
```

</details>


<details open>
<summary>Python</summary>

To copy the files for [this example ( Part 3 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/python/part3) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/python:main+part3/part3 ./part3
```

`./Earthfile`

```Dockerfile
VERSION 0.8
FROM python:3
WORKDIR /code

deps:
    RUN pip install wheel
    COPY requirements.txt ./
    RUN pip wheel -r requirements.txt --wheel-dir=wheels
    SAVE ARTIFACT wheels /wheels

build:
    FROM +deps
    COPY src src
    SAVE ARTIFACT src /src

docker:
    COPY +deps/wheels wheels
    COPY +build/src src
    COPY requirements.txt ./
    RUN pip install --no-index --find-links=wheels -r requirements.txt
    ENTRYPOINT ["python3", "./src/hello.py"]
    SAVE IMAGE python-example:latest
```

</details>
