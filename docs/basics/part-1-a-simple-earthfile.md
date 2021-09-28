# A simple Earthfile

Earthfiles are always named `Earthfile`, regardless of their location in the codebase.

{% method %}
{% sample lang="Go" %}
Here is a sample Earthfile of a Go app

`./Earthfile`

```Dockerfile
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

`./main.go`

```go
package main

import "fmt"

func main() {
	fmt.Println("hello world")
}
```

{% hint style='info' %}

##### Note

To copy the files for [this example ( Part 1 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/go/part1) run

```bash
mkdir tutorial
cd tutorial
earthly --artifact github.com/earthly/earthly/examples/tutorial/go:main+part1/part1 ./part1
```

{% endhint %}

{% sample lang="JavaScript" %}
Here is a sample Earthfile of a JS app

`./Earthfile`

```Dockerfile
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

{% hint style='info' %}

##### Note

To copy the files for [this example ( Part 1 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/js/part1) run

```bash
mkdir tutorial
cd tutorial
earthly --artifact github.com/earthly/earthly/examples/tutorial/js:main+part1/part1 ./part1
```

{% endhint %}

{% sample lang="Java" %}
Here is a sample Earthfile of a Java app

`./Earthfile`

```Dockerfile
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

{% hint style='info' %}

##### Note

To copy the files for [this example ( Part 1 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/java/part1) run

```bash
mkdir tutorial
cd tutorial
earthly --artifact github.com/earthly/earthly/examples/tutorial/java:main+part1/part1 ./part1
```

{% endhint %}

{% sample lang="Python" %}
Here is a sample Earthfile of a Python app

`./Earthfile`

```Dockerfile
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

`./src/hello.py`

```python
print("hello world")
```

{% hint style='info' %}

##### Note

To copy the files for [this example ( Part 1 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/python/part1) run

```bash
mkdir tutorial
cd tutorial
earthly --artifact github.com/earthly/earthly/examples/tutorial/python:main+part1/part1 ./part1
```

{% endhint %}

{% endmethod %}

From the example above, you may notice that an Earthfile is very similar to a Dockerfile. This is an intentional design decision. Existing Dockerfiles can easily be ported to Earthly by copying them to an Earthfile and tweaking them slightly.

<!-- Compared to Dockerfile syntax, there are some additional commands specific to Earthly (like `SAVE ARTIFACT`), some commands have additional semantics (like `COPY +target/some-artifact`) and other semantics have been removed (like `FROM ... AS ...` and `COPY --from`). -->

You may notice the command `COPY +build/... ...`, which has an unfamiliar form. This is a special type of `COPY`, which can be used to pass artifacts from one target to another. In this case, the target `build` (referenced as `+build`) produces an artifact, which has been declared with `SAVE ARTIFACT`, and the target `docker` copies that artifact in its build environment.

With Earthly you have the ability to pass such artifacts or images between targets within the same Earthfile, but also across different Earthfiles across directories or even across repositories. To read more about this, see the [target, artifact and image referencing guide](../guides/target-ref.md).
