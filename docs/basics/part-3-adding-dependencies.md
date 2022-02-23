
# Adding dependencies

Now Let's imagine that we want to add some dependancies to our app. Here's how our build might look as a result.

- [Go](#go) 
- [JavaScript](#javascript) 
- [Java](#java)
- [Python](#python)

### Go

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

The build then might become

`./Earthfile`

```Dockerfile
VERSION 0.6
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

{% hint style='info' %}
##### Note

To copy the files for [this example ( Part 3 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/go/part3) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/go:main+part3/part3 ./part3
```
{% endhint %}

### Javascript

`./package.json`

```json
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

`./package-lock.json` (empty)

```json
```

The code of the app might look like this

`./src/index.js`

```js
function component() {
    const element = document.createElement('div');
    element.innerHTML = "hello world"
    return element;
}

document.body.appendChild(component());
```

`./src/index.html`

```html
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

`./Earthfile`

```Dockerfile
VERSION 0.6
FROM node:13.10.1-alpine3.11
WORKDIR /js-example

build:
    COPY package.json package-lock.json ./
    COPY src src
    RUN mkdir -p ./dist && cp ./src/index.html ./dist/
    RUN npm install
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

To copy the files for [this example ( Part 3 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/js/part3) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/js:main+part3/part3 ./part3
```
{% endhint %}

### Java

`./build.gradle`

```groovy
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

`./src/main/java/hello/HelloWorld.java`

```java
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

`./Earthfile`

```Dockerfile
VERSION 0.6
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

{% hint style='info' %}
##### Note

To copy the files for [this example ( Part 3 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/java/part3) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/java:main+part3/part3 ./part3
```
{% endhint %}

### Python

`./requirements.txt`

```
Markdown==3.2.2
```
The code of the app would now look like this

`./src/hello.py`

```python
from markdown import markdown

def hello():
    return markdown("## hello world")

print(hello())
```

The build might then become as follows.

`./Earthfile`

```Dockerfile
VERSION 0.6
FROM python:3
WORKDIR /code

build:
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

{% hint style='info' %}
##### Note

To copy the files for [this example ( Part 3 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/python/part3) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/python:main+part3/part3 ./part3
```
{% endhint %}

However, as we build this new setup and make changes to the main source code, we notice that the dependencies are downloaded every single time we change the source code. While the build is not necessarily incorrect, it is inefficient for proper development speed.

To improve the speed we will make some changes in part 4 of the tutorial.
