# Reduce code duplication

In some cases, the dependencies might be used in more than one build target. For this use case, we might want to separate dependency downloading and reuse it. For this reason, let's consider breaking this out into a separate build target, called `deps`. We can then inherit from `deps` by using the command `FROM +deps`.

Note that in our case, only the JavaScript version has an example where `FROM +deps` is used in more than one place: both in `build` and in `docker`. Nevertheless, all versions show how dependencies may be separated.

{% method %}
{% sample lang="Go" %}
`./Earthfile`

```Dockerfile
FROM golang:1.15-alpine3.13
WORKDIR /go-example

deps:
    COPY go.mod go.sum ./
	RUN go mod download
	SAVE ARTIFACT go.mod AS LOCAL go.mod
	SAVE ARTIFACT go.sum AS LOCAL go.sum

build:
    FROM +deps
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

To copy these files locally run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/go:main+part5/part5 ./part5
```
{% endhint %}

{% sample lang="JavaScript" %}
`./Earthfile`

```Dockerfile
FROM node:13.10.1-alpine3.11
WORKDIR /js-example

deps:
    COPY package.json ./
    COPY package-lock.json ./
    RUN npm install
    SAVE ARTIFACT package.json AS LOCAL ./package.json
    SAVE ARTIFACT package-lock.json AS LOCAL ./package-lock.json

build:
    FROM +deps
    COPY src src
    COPY dist dist
    RUN npx webpack
    SAVE ARTIFACT dist /dist AS LOCAL dist

docker:
    FROM +deps
    COPY +build/dist ./dist
    EXPOSE 8080
    ENTRYPOINT ["/js-example/node_modules/http-server/bin/http-server", "./dist"]
    SAVE IMAGE js-example:latest
```

{% hint style='info' %}
##### Note

To copy these files locally run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/js:main+part5/part5 ./part5
```
{% endhint %}

{% sample lang="Java" %}
`./Earthfile`

```Dockerfile
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

{% hint style='info' %}
##### Note

To copy these files locally run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/java:main+part5/part5 ./part5
```
{% endhint %}

{% sample lang="Python" %}
`./Earthfile`

```Dockerfile
FROM python:3
WORKDIR /code

deps:
    RUN pip install wheel
    COPY requirements.txt ./
    RUN pip wheel -r requirements.txt --wheel-dir=wheels

build:
    FROM +deps
    COPY src src
    SAVE ARTIFACT src /src
    SAVE ARTIFACT wheels /wheels

docker:
    COPY +build/src src
    COPY +build/wheels wheels
    COPY requirements.txt ./
    RUN pip install --no-index --find-links=wheels -r requirements.txt
    ENTRYPOINT ["python3", "./src/hello.py"]
    SAVE IMAGE python-example:latest
```

{% hint style='info' %}
##### Note

To copy these files locally run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/python:main+part5/part5 ./part5
```
{% endhint %}

{% endmethod %}

## Continue tutorial

👉 [Final words](./final-words.md)
