To copy the files for [this example ( Part 4 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/go/part4) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/go:main+part4/part4 ./part4
```

Examples in [Python](#more-examples), [JavaScript](#more-examples) and [Java](#more-examples) are at the bottom of this page.

## Just Like Docker...Mostly

`ARG`s in Earthly work similar to `ARG`s in Dockerfiles, however there are a few differences when it comes to scope. Also, Earthly has a number of [built in `ARG`s](../earthfile/builtin-args.md) that are available to use.

Let's say we wanted to have the option to pass in a tag for our Docker image when we run `earthly +docker`.

```Dockerfile
docker:
    ARG tag='latest'
    COPY +build/example .
    ENTRYPOINT ["/go-workdir/example"]
    SAVE IMAGE go-example:$tag
```
In our `+docker` target we can create an `ARG` called tag. In this case, we give it a default value of `latest`. If we do not provide a default value the default will be an empty string.

Then, down in our `SAVE IMAGE` command, we are able to reference the `ARG` with `$` followed by the `ARG` name.

Now we can take advantage of this when we run Earthly.

```bash
earthly +docker --tag='my-new-image-tag'
```
In this case `my-new-image-tag` will override the default value and become the new tag for our docker image. If we hadn't passed in a value for tag, then the default `latest` would have been used. 

```bash
earthly +docker
# tag for image will be 'latest'
```

### Passing ARGs in FROM, BUILD, and COPY
We can also pass `ARG`s when referencing a target inside an Earthfile. Using the `FROM` and `BUILD` commands, this looks pretty similar to what we did above on the command line.

```Dockerfile
docker:
    ARG tag='latest'
    COPY +build/example .
    ENTRYPOINT ["/go-workdir/example"]
    SAVE IMAGE go-example:$tag

with-build:
    BUILD +docker --tag='my-new-image-tag'

with-from:
    FROM +docker --tag='my-new-image-tag'
```
We can also pass `ARG`s when using the `COPY` command, though the syntax is a little different.

```Dockerfile
build:
    ARG version
    COPY main.go .
    RUN go build -o output/example-$version main.go
    SAVE ARTIFACT output/example-$version AS LOCAL local-output/go-example

with-copy:
    COPY (+build/example --version='1.0') .
```

## Builtin `ARG`s
There are a number of builtin `ARG`s that Earthly offers. You can read about a [complete list of them](https://docs.earthly.dev/docs/earthfile/builtin-args), but for now, let's take a look at how they work.

**In order to use Earthly builtin `ARG`s they need to be pre-declared.** Once you do that, you can use them just like any other `ARG`.

```Dockerfile
ARG USERARCH
RUN echo $USERARCH
```
In this case we've declared the `ARG` `USERARCH` which is a builtin that holds the processor architecture the target is being built from.


## More Examples

<details open>
<summary>JavaScript</summary>

To copy the files for [this example ( Part 4 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/js/part4) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/js:main+part4/part4 ./part4
```

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
    ARG tag='latest'
    COPY +build/dist ./dist
    EXPOSE 8080
    ENTRYPOINT ["/js-example/node_modules/http-server/bin/http-server", "./dist"]
    SAVE IMAGE js-example:$tag
```

</details>


<details open>
<summary>Java</summary>

To copy the files for [this example ( Part 4 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/java/part4) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/java:main+part4/part4 ./part4
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
    ARG tag='latest'
    ENTRYPOINT ["/java-example/bin/java-example"]
    SAVE IMAGE java-example:$tag
```

</details>


<details open>
<summary>Python</summary>

To copy the files for [this example ( Part 4 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/python/part4) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/python:main+part4/part4 ./part4
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
    ARG tag='latest'
    RUN pip install --no-index --find-links=wheels -r requirements.txt
    ENTRYPOINT ["python3", "./src/hello.py"]
    SAVE IMAGE python-example:$tag
```

</details>
