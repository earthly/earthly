To copy the files for [this example ( Part 2 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/go/part2) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/go:main+part2/part2 ./part2
```
Examples in [Python](#more-examples), [JavaScript](#more-examples) and [Java](#more-examples) are at the bottom of this page.

## Not All Targets Produce Output
Targets have the ability to produce output outside of the build environment. You can save files and docker images to your local machine or push them to remote repositories. Targets can also run commands that affect the local environment outside of the build, such as running database migrations, but not all targets produce output. Let's take a look at which commands produce output and how to use them.

## Saving Files
We've already seen how the command [SAVE ARTIFACT](https://docs.earthly.dev/docs/earthfile#save-artifact) copies a file or directory from the build environment into the target's artifact environment.

This gives us the ability to copy files between targets, **but it does not allow us to save any files to our local machine.**

```Dockerfile
build:
    COPY main.go .
    RUN go build -o output/example main.go
    SAVE ARTIFACT output/example

docker:
    #  COPY command copies files from the +build target
    COPY +build/example .
    ENTRYPOINT ["/go-workdir/example"]
    SAVE IMAGE go-example:latest
```
In order to **save the file locally** , we need to add `AS LOCAL` to the command.

```Dockerfile
build:
    COPY main.go .
    RUN go build -o output/example main.go
    SAVE ARTIFACT output/example AS LOCAL local-output/go-example
```

If we run this example with `earthly +build`, we'll see a `local-output` directory show up locally with a `go-example` file inside of it.

## Saving Docker Images
Saving Docker images to your local machine is easy with the `SAVE IMAGE` command.

```Dockerfile
build:
    COPY main.go .
    RUN go build -o output/example main.go
    SAVE ARTIFACT output/example

docker:
    COPY +build/example .
    ENTRYPOINT ["/go-workdir/example"]
    SAVE IMAGE go-example:latest
```
In this example, running `earthly +docker` will save an image named `go-example` with the tag `latest`.

```bash
~$ earthly +docker
...
~$ docker image ls
REPOSITORY          TAG       IMAGE ID       CREATED          SIZE
go-example          latest    08b9f749023d   19 seconds ago   297MB

# or podman
~$ podman image ls
REPOSITORY          TAG       IMAGE ID       CREATED          SIZE
go-example          latest    08b9f749023d   19 seconds ago   297MB
```
**NOTE**

If we run a target as a reference in `FROM` or `COPY`, **outputs will not be produced**. Take this Earthfile for example.

```Dockerfile
build:
    COPY main.go .
    RUN go build -o output/example main.go
    SAVE ARTIFACT output/example AS LOCAL local-output/go-example

docker:
    COPY +build/example .
    ENTRYPOINT ["/go-workdir/example"]
    SAVE IMAGE go-example:latest
```
In this case, running `earthly +docker` will not produce any output. In other words, you will not have a `local-output/go-example` written locally, but running `earthly +build` will still produce output as expected.

The exception to this rule is the `BUILD` command. If you want to use `COPY` or `FROM` and still have Earthly create `local-output/go-example` locally, you'll need to use the `BUILD` command to do so.

```Dockerfile
build:
    COPY main.go .
    RUN go build -o output/example main.go
    SAVE ARTIFACT output/example AS LOCAL local-output/go-example

docker:
    BUILD +build
    COPY +build/example .
    ENTRYPOINT ["/go-workdir/example"]
    SAVE IMAGE go-example:latest
```
Running `earthly +docker` in this case will now output `local-output/go-example` locally.

## The Push Flag

### Docker Images

In addition to saving files and images locally, we can also push them to remote repositories.

```Dockerfile
docker:
    COPY +build/example .
    ENTRYPOINT ["/go-workdir/example"]
    SAVE IMAGE --push go-example:latest
```
Note that adding the `--push` flag to `SAVE IMAGE` is not enough, we'll also need to invoke push when we run earthly. `earthly --push +docker`.

#### External Changes
You can also use `--push` as part of a `RUN` command to define commands that have an effect external to the build. These kinds of effects are only allowed to take place if the entire build succeeds.

This allows you to push to remote repositories. 

```Dockerfile
release:
    RUN --push --secret GITHUB_TOKEN=GH_TOKEN github-release upload
```
```bash
earthly --push +release
```
But also allows you to do things like run database migrations.

```Dockerfile
migrate:
    FROM +build
    RUN --push bundle exec rails db:migrate
```
```bash
earthly --push +migrate
```
Or apply terraform changes

```Dockerfile
apply:
    RUN --push terraform apply -auto-approve
```
```bash
earthly --push +apply
```
**NOTE**

Just like saving files, any command that uses `--push` **will only produce output if called directly**, `earthly --push +target-with-push` **or via a** `BUILD` command. Calling a target via `FROM` or `COPY` will not invoke `--push`.

### More Examples
<details open>
<summary>JavaScript</summary>

To copy the files for [this example ( Part 2 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/js/part2) run

```bash
mkdir tutorial
cd tutorial
earthly --artifact github.com/earthly/earthly/examples/tutorial/js:main+part2/part2 ./part2
```

`./Earthfile`

```Dockerfile
VERSION 0.8
FROM node:13.10.1-alpine3.11
WORKDIR /js-example

build:
    # In JS, there's nothing to build in this simple form.
    # The source is also the artifact used in production.
    COPY src/index.js .
    SAVE ARTIFACT index.js /dist/index.js AS LOCAL ./dist/index.js

docker:
    COPY +build/dist dist
    ENTRYPOINT ["node", "./dist/index.js"]
    SAVE IMAGE js-example:latest
```

The code of the app might look like this

`./src/index.js`

```js
console.log("hello world");
```

</details>


<details open>
<summary>Java</summary>

To copy the files for [this example ( Part 2 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/java/part2) run

```bash
mkdir tutorial
cd tutorial
earthly --artifact github.com/earthly/earthly/examples/tutorial/java:main+part2/part2 ./part2
```

`./Earthfile`

```Dockerfile
VERSION 0.8
FROM openjdk:8-jdk-alpine
RUN apk add --update --no-cache gradle
WORKDIR /java-example

build:
    COPY build.gradle ./
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

The code of the app might look like this

`./src/main/java/hello/HelloWorld.java`

```java

package hello;

public class HelloWorld {
    public static void main(String[] args) {
        System.out.println("hello world");
    }
}
```

`./build.gradle`

```groovy
apply plugin: 'java'
apply plugin: 'application'

mainClassName = 'hello.HelloWorld'

jar {
    baseName = 'hello-world'
    version = '0.0.1'
}

sourceCompatibility = 1.8
targetCompatibility = 1.8
```
</details>


<details open>
<summary>Python</summary>

To copy the files for [this example ( Part 2 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/python/part2) run

```bash
mkdir tutorial
cd tutorial
earthly --artifact github.com/earthly/earthly/examples/tutorial/python:main+part2/part2 ./part2
```

`./Earthfile`

```Dockerfile
VERSION 0.8
FROM python:3
WORKDIR /code

build:
     # In Python, there's nothing to build.
    COPY src src
    SAVE ARTIFACT src /src

docker:
    COPY +build/src src
    ENTRYPOINT ["python3", "./src/hello.py"]
    SAVE IMAGE --push python-example:latest
```

The code of the app might look like this

`./src/hello.py`

```python
print("hello world")
```

</details>
