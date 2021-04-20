
# Adding dependencies in the mix

Let's imagine now that in our simple app, we now want to add a programming language dependency. Here's how our build might look like as a result

{% method %}
{% sample lang="Go" %}
```go.mod
// go.mod

module github.com/earthly/earthly/examples/go

go 1.13

require github.com/sirupsen/logrus v1.5.0
```

The code of the app might look like this

```go
// main.go

package main

import "github.com/sirupsen/logrus"

func main() {
	logrus.Info("hello world")
}
```

The build then might become

```Dockerfile
# Earthfile

FROM golang:1.15-alpine3.13
WORKDIR /go-example

build:
    COPY go.mod go.sum .
    COPY main.go .
    RUN go build -o build/go-example main.go
    SAVE ARTIFACT build/go-example /go-example AS LOCAL build/go-example

docker:
    COPY +build/go-example .
    ENTRYPOINT ["/go-example/go-example"]
    SAVE IMAGE go-example:latest
```
{% sample lang="JavaScript" %}
```json
// package.json

{
  "name": "example-js",
  "version": "0.0.1",
  "description": "Hello world",
  "private": true,
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "author": "",
  "license": "MPL-2.0",
  "devDependencies": {
    "webpack": "^4.42.1",
    "webpack-cli": "^3.3.11"
  },
  "dependencies": {
    "http-server": "^0.12.1"
  }
}
```

The code of the app might look like this

```js
// src/index.js

function component() {
    const element = document.createElement('div');
    element.innerHTML = "hello world"
    return element;
}

document.body.appendChild(component());
```

```html
<!-- dist/index.html -->

<!doctype html>
<html>

<head>
    <title>Getting Started</title>
</head>

<body>
    <script src="./main.js"></script>
</body>

</html>
```

The build then might become

```Dockerfile
# Earthfile

FROM node:13.10.1-alpine3.11
WORKDIR /js-example

build:
    COPY package.json package-lock.json ./
    COPY src src
    COPY dist dist
    RUN npm install
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
```groovy
// build.gradle

apply plugin: 'java'
apply plugin: 'application'

mainClassName = 'hello.HelloWorld'

repositories {
    mavenCentral()
}

jar {
    baseName = 'hello-world'
    version = '0.0.1'
}

sourceCompatibility = 1.8
targetCompatibility = 1.8

dependencies {
    compile "joda-time:joda-time:2.2"
    testCompile "junit:junit:4.12"
}
```

The code of the app might look like this

```java
// src/main/java/hello/HelloWorld.java

package hello;

import org.joda.time.LocalTime;

public class HelloWorld {
    public static void main(String[] args) {
        LocalTime currentTime = new LocalTime();
        System.out.println(currentTime + " - hello world");
    }
}
```

The Earthfile file would not change

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
{% sample lang="Python" %}
```
// Requirements.txt

Markdown==3.2.2
```
The code of the app would now look like this
```python
# src/hello.py

from markdown import markdown

def hello():
    return markdown("Hello *Earthly*")

print(hello())
```
The build might then become as follows.  
```Docker
# EarthFile

FROM python:3
WORKDIR /code

build:
    # Use Python Wheels to produce package files into /wheels
    RUN pip install wheel
    COPY requirements.txt ./
    RUN pip wheel -r requirements.txt --wheel-dir=wheels
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
{% endmethod %}

However, as we build this new setup and make changes to the main source code, we notice that the dependencies are downloaded every single time we change the source code. While the build is not necessarily incorrect, it is inefficient for proper development speed.

To improve the speed we will make some changes in part 4 of the tutorial.

## Continue tutorial

👉 [Part 4: Efficient caching of dependencies](./part-4-efficient-caching-of-dependencies.md)
