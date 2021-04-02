
# A simple Earthfile

Earthfiles are always named `Earthfile`, regardless of their location in the codebase.

{% method %}
{% sample lang="Go" %}
Here is a sample earthfile of a Go app

```Dockerfile
# Earthfile

FROM golang:1.15-alpine3.13
WORKDIR /go-example

build:
    COPY main.go .
    RUN go build -o build/go-example main.go
    SAVE ARTIFACT build/go-example /go-example AS LOCAL build/go-example

docker:
    COPY +build/go-example .
    ENTRYPOINT ["/go-example/go-example"]
    SAVE IMAGE go-example:latest
```

The code of the app might look like this

```go
// main.go

package main

import "fmt"

func main() {
	fmt.Println("hello world")
}
```

{% sample lang="JavaScript" %}
Here is a sample earthfile of a JS app

```Dockerfile
# Earthfile

FROM node:13.10.1-alpine3.11
WORKDIR /js-example

build:
    # In JS, there's nothing to build in this simple form.
    # The source is also the artifact used in production.
    COPY index.js .
    SAVE ARTIFACT index.js /dist/index.js AS LOCAL ./dist/index.js

docker:
    COPY +build/dist dist
    ENTRYPOINT ["node", "./dist/index.js"]
    SAVE IMAGE js-example:latest
```

The code of the app might look like this

```js
// index.js

console.log("hello world");
```

{% sample lang="Java" %}
Here is a sample earthfile of a Java app

```Dockerfile
# Earthfile

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

```java
// src/main/java/hello/HelloWorld.java

package hello;

public class HelloWorld {
    public static void main(String[] args) {
        System.out.println("hello world");
    }
}
```
{% sample lang="Python" %}
Here is a sample earthfile of a Python app

```Dockerfile
# Earthfile

FROM python:3
WORKDIR /code

build:
     # In Python, there's nothing to build.
    COPY src src
    SAVE ARTIFACT src /src

docker:
    COPY +build/src src
    ENTRYPOINT ["python3", "./src/hello.py"]
    SAVE IMAGE python-example:latest
```

The code of the app might look like this

```python
// src/hello.py

print("hello world")
```
{% endmethod %}

You will notice that the recipes look very much like Dockerfiles. This is an intentional design decision. Existing Dockerfiles can be ported to earthfiles by copy-pasting them over and then tweaking them slightly. Compared to Dockerfile syntax, some commands are new (like `SAVE ARTIFACT`), others have additional semantics (like `COPY +target/some-artifact`) and other semantics have been removed (like `FROM ... AS ...` and `COPY --from`).

You might notice the command `COPY +build/... ...`, which has an unfamiliar form. This is a special type of `COPY`, which can be used to pass artifacts from one target to another. In this case, the target `build` (referenced as `+build`) produces an artifact, which has been declared with `SAVE ARTIFACT`, and the target `docker` copies that artifact in its build environment.

With Earthly you have the ability to pass such artifacts or images between targets within the same Earthfile, but also across different Earthfiles across directories or even across repositories. To read more about this, see the [target, artifact and image referencing guide](./target-ref.md).

<img src="img/ref-infographic.png" alt="Target and artifact reference syntax" title="Reference targets using +" width="500px" />

## Executing a build

In this particular example, we can see two explicit targets: `build` and `docker`. In order to execute the build, we can run, for example:

```bash
earthly +build
```

or

```bash
earthly +docker
```

The output might look like this:

![Earthly build output](img/go-example.png)

Notice how to the left of `|`, within the output, we can see some targets like `+base`, `+build` and `+docker` . Notice how the output is interleaved between `+docker` and `+build`. This is because the system executes independent build steps in parallel. The reason this is possible effortlessly is because only very few things are shared between the builds of the recipes and those things are declared and obvious. The rest is completely isolated.

In addition, notice how even though the base is used as part of both `build` and `docker`, it is only executed once. This is because the system deduplicates execution, where possible.

Furthermore, the fact that the `docker` target depends on the `build` target is visible within the command `COPY +build/...`. Through this command, the system knows that it also needs to build the target `+build`, in order to satisfy the dependency on the artifact.

Finally, notice how the output of the build (the docker image and the files) are only written after the build is declared a success. This is due to another isolation principle of Earthly: a build either succeeds completely or it fails altogether.

Once the build has executed, we can run the resulting docker image to try it out:

{% method %}
{% sample lang="Go" %}
```
$ docker run --rm go-example:latest
hello world
```
{% sample lang="JavaScript" %}
```
$ docker run --rm js-example:latest
hello world
```
{% sample lang="Java" %}
```
$ docker run --rm java-example:latest
hello world
```
{% sample lang="Python" %}
```
$ docker run --rm python-example:latest
hello world
```
{% endmethod %}

{% hint style='info' %}
##### Note

Targets have a particular referencing convention which helps Earthly to identify which recipe to execute. In the simplest form, targets are referenced by `+<target-name>`.  For example, `+build`. For more details see the [target referencing page](./target-ref.md).
{% endhint %}

## Continue tutorial

👉 [Part 2: Detailed explanation](./part-2-detailed-explanation.md)
