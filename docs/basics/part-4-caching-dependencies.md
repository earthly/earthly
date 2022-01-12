# Caching dependencies

The reason the build is inefficient is because we have not made proper use of caching. When a file changes, the corresponding `COPY` command is re-executed without cache, causing all commands after it to also re-execute without cache.

If, however, we could first download the dependencies and only afterwards copy and build the code, then the cache would be reused every time we changed the code.

{% method %}
{% sample lang="Go" %}
`./Earthfile`

```Dockerfile
VERSION 0.6
FROM golang:1.15-alpine3.13
WORKDIR /go-example

build:
    # Download deps before copying code.
    COPY go.mod go.sum .
    RUN go mod download
    # Copy and build code.
    COPY main.go .
    RUN go build -o build/go-example main.go
    SAVE ARTIFACT build/go-example /go-example AS LOCAL build/go-example

docker:
    COPY +build/go-example .
    ENTRYPOINT ["/go-example/go-example"]
    SAVE IMAGE go-example:latest
```

{% hint style='info' %}
##### Note

To copy the files for [this example ( Part 4 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/go/part4) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/go:main+part4/part4 ./part4
```
{% endhint %}

{% sample lang="JavaScript" %}
`./Earthfile`

```Dockerfile
VERSION 0.6
FROM node:13.10.1-alpine3.11
WORKDIR /js-example

build:
    # Download deps before copying code.
    COPY package.json package-lock.json ./
    RUN npm install
    # Copy and build code.
    COPY src src
    RUN mkdir -p ./dist && cp ./src/index.html ./dist/
    RUN npx webpack
    SAVE ARTIFACT dist /dist AS LOCAL ./dist

docker:
    COPY package.json package-lock.json ./
    RUN npm install
    COPY +build/dist dist
    EXPOSE 8080
    ENTRYPOINT ["/js-example/node_modules/http-server/bin/http-server", "./dist"]
    SAVE IMAGE js-example:latest
```

{% hint style='info' %}
##### Note

To copy the files for [this example ( Part 4 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/js/part4) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/js:main+part4/part4 ./part4
```
{% endhint %}

{% sample lang="Java" %}
`./Earthfile`

```Dockerfile
VERSION 0.6
FROM openjdk:8-jdk-alpine
RUN apk add --update --no-cache gradle
WORKDIR /java-example

build:
    # Download deps before copying code.
    COPY build.gradle ./
    RUN gradle build
    # Copy and build code.
    COPY src src
    RUN gradle build
    RUN gradle install
    SAVE ARTIFACT build/install/java-example/bin /bin AS LOCAL build/bin
    SAVE ARTIFACT build/install/java-example/lib /lib AS LOCAL build/lib

docker:
    COPY +build/bin bin
    COPY +build/lib lib
    ENTRYPOINT ["/java-example/bin/java-example"]
    SAVE IMAGE java-example:latest
```

{% hint style='info' %}
##### Note

To copy the files for [this example ( Part 4 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/java/part4) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/java:main+part4/part4 ./part4
```
{% endhint %}

{% sample lang="Python" %}

`./Earthfile`

```Dockerfile
VERSION 0.6
FROM python:3
WORKDIR /code

build:
    # Download deps before copying code.
    RUN pip install wheel
    COPY requirements.txt ./
    RUN pip wheel -r requirements.txt --wheel-dir=wheels
    SAVE ARTIFACT wheels /wheels
    # Copy and build code.
    COPY src src
    SAVE ARTIFACT src /src

docker:
    COPY +build/wheels wheels
    COPY +build/src src
    COPY requirements.txt ./
    RUN pip install --no-index --find-links=wheels -r requirements.txt
    ENTRYPOINT ["python3", "./src/hello.py"]
    SAVE IMAGE python-example:latest
```

{% hint style='info' %}
##### Note

To copy the files for [this example ( Part 4 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/python/part4) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/python:main+part4/part4 ./part4
```
{% endhint %}
{% endmethod %}

For a primer into Dockerfile layer caching see [this article](https://pythonspeed.com/articles/docker-caching-model/). The same principles apply to Earthfiles.
