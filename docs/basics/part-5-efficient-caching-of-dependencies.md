# Efficient caching of dependencies

The reason the build is inefficient is because we have not made proper use of caching. When a file changes, the corresponding `COPY` command is re-executed without cache, causing all commands after it to also re-execute without cache.

If, however, we could first download the dependencies and only afterwards copy and build the code, then the cache would be reused every time we changed the code.

{% method %}
{% sample lang="Go" %}
```Dockerfile
# Earthfile

FROM golang:1.15-alpine3.13
WORKDIR /go-example

build:
    # Download deps before copying code.
    COPY go.mod go.sum .
    RUN go mod download
    # Also save these back to host, in case go.sum changes.
    SAVE ARTIFACT go.mod AS LOCAL go.mod
	SAVE ARTIFACT go.sum AS LOCAL go.sum
    # Copy and build code.
    COPY main.go .
    RUN go build -o build/go-example main.go
    SAVE ARTIFACT build/go-example /go-example AS LOCAL build/go-example

docker:
    COPY +build/go-example .
    ENTRYPOINT ["/go-example/go-example"]
    SAVE IMAGE go-example:latest
```
{% sample lang="JavaScript" %}
```Dockerfile
# Earthfile

FROM node:13.10.1-alpine3.11
WORKDIR /js-example

build:
    # Download deps before copying code.
    COPY package.json package-lock.json ./
    RUN npm install
    # Also save these back to host, in case package-lock.json changes.
    SAVE ARTIFACT package.json AS LOCAL ./package.json
    SAVE ARTIFACT package-lock.json AS LOCAL ./package-lock.json
    # Copy and build code.
    COPY src src
    COPY dist dist
    RUN npx webpack
    SAVE ARTIFACT dist /dist AS LOCAL ./dist

docker:
    COPY package.json package-lock.json ./
    RUN npm install
    COPY +build/dist dist
    ENTRYPOINT ["node", "./dist/index.js"]
    SAVE IMAGE js-example:latest
```
{% sample lang="Java" %}
```Dockerfile
# Earthfile

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
{% sample lang="Python" %}
```Docker
# EarthFile
FROM python:3
WORKDIR /code

build:
    RUN pip install wheel
    COPY requirements.txt ./
    RUN pip wheel -r requirements.txt --wheel-dir=wheels

    #save wheels before copy source, for cache efficiency 
    SAVE ARTIFACT wheels /wheels

    COPY src src
    SAVE ARTIFACT src /src

docker:
    COPY +build/src src
    COPY +build/wheels wheels
    COPY requirements.txt ./
    RUN pip install --no-index --find-links=wheels -r requirements.txt
    ENTRYPOINT ["python3", "./src/hello.py"]
    SAVE IMAGE python-example:latest
```
{% endmethod %}

For a primer into Dockerfile layer caching see [this article](https://pythonspeed.com/articles/docker-caching-model/). The same principles apply to Earthfiles.

## Continue tutorial

👉 [Part 6: Reduce code duplication](./part-6-reduce-code-duplication.md)
