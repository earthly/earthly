To copy the files for [this example ( Part 5 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/go/part5) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/go:main+part5/part5 ./part5
```

Examples in [Python](#more-examples), [JavaScript](#more-examples) and [Java](#more-examples) are at the bottom of this page.

## Calling on Targets From Other Earthfiles

So far we've seen how the `FROM` command in Earthly has the ability to reference another target's image as its base image, like in the case below where the `+build` target uses the image from the `+deps` target.

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

```

But `FROM` also has the ability to import targets from Earthfiles in different directories. Let's say we have a directory structure like this.

```
.
├── services
|   └── service-one
|       ├── Earthfile (containing +deps)
|       ├── go.mod
|       └── go.sum
├── main.go
└── Earthfile

```
We can use a target in the Earthfile in `/services/service-one` from inside the Earthfile in the root of our directory. NOTE: relative paths must use `./` or `../`.

`./services/service-one/Earthfile`

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
```

`./Earthfile`

```Dockerfile
build:
    FROM ./services/service-one+deps
    COPY main.go .
    RUN go build -o output/example main.go
    SAVE ARTIFACT output/example AS LOCAL local-output/go-example
```
This code tells `FROM` that there is another Earthfile in  the `services/service-one` directory and that the Earthfile  contains a target called `+deps`. In this case, if we were to run `+build` Earthly is smart enough to go into the subdirectory, run the  `+deps` target in that Earthfile, and then use it as the base image for `+build`.

We can also reference an Earthfile in another repo, which works in a similar way. If the reference does not begin with one of `/`, `./`, or `../`, then earthly treats it as a repository.  See [the reference](../earthfile/earthfile.md#from) for details.

```Dockerfile
build:
    FROM github.com/example/project+remote-target
    COPY main.go .
    RUN go build -o output/example main.go
    SAVE ARTIFACT output/example AS LOCAL local-output/go-example
```

## Importing Whole Projects
In addition to importing single targets from other files, we can also import an entire Earthfile with the `IMPORT` command. This is helpful if there are several targets in a separate Earthfile that you want access to in your current file. It also allows you to create an alias.

```Dockerfile
VERSION 0.8
IMPORT ./services/service-one AS my_service
FROM golang:1.15-alpine3.13
WORKDIR /go-workdir

build:
    FROM my_service+deps
    COPY main.go .
    RUN go build -o output/example main.go
    SAVE ARTIFACT output/example AS LOCAL local-output/go-example
```
In this example, we assume there is a `./services/service-one` directory that contains its own Earthfile. We import it and then use the `AS` keyword to give it an alias.

Then, in our `+build` target we can inherit from any target in the imported Earthfile by passing `alias+target-name`. In this case the Earthfile in the service directory has a target named `+deps`.

## More Examples

<details open>
<summary>JavaScript</summary>

To copy the files for [this example ( Part 5 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/js/part5) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/js:main+part5/part5 ./part5
```

`./Earthfile`

```Dockerfile
VERSION 0.8
FROM node:13.10.1-alpine3.11
WORKDIR /js-example

build:
    FROM ./services/service-one+deps
    COPY src src
    RUN mkdir -p ./dist && cp ./src/index.html ./dist/
    RUN npx webpack
    SAVE ARTIFACT dist /dist AS LOCAL dist

docker:
    FROM +deps
    ARG tag='latest'
    COPY +build/dist ./dist
    EXPOSE 8080
    ENTRYPOINT ["/js-example/node_modules/http-server/bin/http-server", "./dist"]
    SAVE IMAGE js-example:$tag
```

</details>

<details open>
<summary>Java</summary>

To copy the files for [this example ( Part 5 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/java/part5) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/java:main+part5/part5 ./part5
```


`./Earthfile`

```Dockerfile
VERSION 0.8
FROM openjdk:8-jdk-alpine
RUN apk add --update --no-cache gradle
WORKDIR /java-example

build:
    FROM ./services/service-one+deps
    COPY src src
    RUN gradle build
    RUN gradle install
    SAVE ARTIFACT build/install/java-example/bin AS LOCAL build/bin
    SAVE ARTIFACT build/install/java-example/lib AS LOCAL build/lib

docker:
    COPY +build/bin bin
    COPY +build/lib lib
    ARG tag='latest'
    ENTRYPOINT ["/java-example/bin/java-example"]
    SAVE IMAGE java-example:$tag
```

</details>

<details open>
<summary>Python</summary>

To copy the files for [this example ( Part 5 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/python/part5) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/python:main+part5/part5 ./part5
```

`./Earthfile`

```Dockerfile
VERSION 0.8
FROM python:3
WORKDIR /code

build:
    FROM ./services/service-one+deps
    COPY src src
    SAVE ARTIFACT src /src

docker:
    COPY +deps/wheels wheels
    COPY +build/src src
    COPY requirements.txt ./
    ARG tag='latest'
    RUN pip install --no-index --find-links=wheels -r requirements.txt
    ENTRYPOINT ["python3", "./src/hello.py"]
    SAVE IMAGE python-example:$tag
```

</details>
