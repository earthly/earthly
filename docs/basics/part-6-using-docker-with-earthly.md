To copy the files for [this example ( Part 6 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/go/part6) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/go:main+part6/part6 ./part6
```

Examples in [Python](#more-examples), [JavaScript](#more-examples) and [Java](#more-examples) are at the bottom of this page.

## The `WITH DOCKER` Command

You may find that you need to run Docker commands inside a target. For those cases Earthly offers `WITH DOCKER`. `WITH DOCKER` will initialize a Docker daemon that can be used in the context of a `RUN` command.

Whenever you need to use `WITH DOCKER` we recommend (though it is not required) that you use Earthly's own Docker in Docker (dind) image: `earthly/dind:alpine-3.19-docker-25.0.5-r0`.

Notice `WITH DOCKER` creates a block of code that has an `END` keyword. Everything that happens within this block is going to take place within our `earthly/dind:alpine-3.19-docker-25.0.5-r0` container.

### Pulling an Image
```Dockerfile
hello:
    FROM earthly/dind:alpine-3.19-docker-25.0.5-r0
    WITH DOCKER --pull hello-world
        RUN docker run hello-world
    END

```
You can see in the command above that we can pass a flag to `WITH DOCKER` telling it to pull an image from Docker Hub. We can pass other flags to [load in artifacts built by other targets](#loading-an-image) `--load` or even images defined by [docker-compose](#a-real-world-example) `--compose`. These images will be available within the context of `WITH DOCKER`'s docker daemon.

### Loading an Image
We can load in an image created by another target with the `--load` flag.

```Dockerfile
my-hello-world:
    FROM ubuntu
    CMD echo "hello world"
    SAVE IMAGE my-hello:latest

hello:
    FROM earthly/dind:alpine-3.19-docker-25.0.5-r0
    WITH DOCKER --load hello:latest=+my-hello-world
        RUN docker run hello:latest
    END
```

## A Real World Example

One common use case for `WITH DOCKER` is running integration tests that require other services. In this case we need to set up a redis service for our tests. For this we can user a `docker-compose.yml`.

`docker-compose.yml`
```yml
version: "3"

services:
  redis:
    container_name: local-redis
    image: redis:6.0-alpine
    ports:
      - 127.0.0.1:6379:6379
    hostname: redis
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:6379"]
      interval: 1s
      timeout: 10s
      retries: 5
    networks:
      - go/part6_default

networks:
  go/part6_default:
```

`main.go`

```go
package main

import (
	"github.com/sirupsen/logrus"
)

var howCoolIsEarthly = "IceCool"

func main() {
	logrus.Info("hello world")
}
```
`main_integration_test.go`

```go
package main

import (
	"context"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "howCoolIsEarthly", howCoolIsEarthly, 0).Err()
	if err != nil {
		panic(err)
	}

	resultFromDB, err := rdb.Get(ctx, "howCoolIsEarthly").Result()
	if err != nil {
		panic(err)
	}
	require.Equal(t, howCoolIsEarthly, resultFromDB)
}
```

```Dockerfile
VERSION 0.8
FROM golang:1.15-alpine3.13
WORKDIR /go-workdir

deps:
    COPY go.mod go.sum ./
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

test-setup:
    FROM +deps
    COPY main.go .
    COPY main_integration_test.go .
    ENV CGO_ENABLED=0
    ENTRYPOINT ["go", "test", "github.com/earthly/earthly/examples/go"]
    SAVE IMAGE test:latest

integration-tests:
    FROM earthly/dind:alpine-3.19-docker-25.0.5-r0
    COPY docker-compose.yml ./
    WITH DOCKER --compose docker-compose.yml --load tests:latest=+test-setup
        RUN docker run --network=default_go/part6_default tests:latest
    END
```
When we use the `--compose` flag, Earthly will start up the services defined in the `docker-compose` file for us. In this case, we built a separate image that copies in our test files and uses the command to run the tests as its `ENTRYPOINT`. We can then load this image into our `WITH DOCKER` command. Note that loading an image will not run it by default, we need to explicitly run the image after we load it.

You'll need to use `--allow-privileged` (or `-P` for short) to run this example. 

```bash
earthly --allow-privileged +integration-tests
```


## More Examples

<details open>
<summary>JavaScript</summary>

To copy the files for [this example ( Part 6 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/js/part6) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/js:main+part6/part6 ./part6
```
In this example, we use `WITH DOCKER` to run a frontend app and backend api together using Earthly.

The App

`./app/package.json`

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

`./app/package-lock.json` (empty)

```json
```

The code of the app might look like this

`./app/src/index.js`

```js
async function getUsers() {

  const response = await fetch('http://0.0.0.0:3080/api/users');
  return await response.json();

}

function component() {
  const element = document.createElement('div');
  getUsers()
    .then( users => {
      element.innerHTML = `hello world <b>${users[0].first_name} ${users[0].last_name}</b>`
    })

	return element;
}

document.body.appendChild(component());
```

`./app/src/index.html`

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
And our api.

`./api/package.json`

```json
{
  "name": "api",
  "version": "1.0.0",
  "description": "",
  "main": "index.js",
  "scripts": {
    "test": "echo \"Error: no test specified\" && exit 1"
  },
  "author": "",
  "license": "ISC",
  "dependencies": {
    "cors": "^2.8.5",
    "express": "^4.17.1",
    "http-proxy-middleware": "^1.0.4",
    "pg": "^8.7.3"
  }
}
```
`./api/package-lock.json` (empty)

```json
```
`./api/server.js`

```js
const express = require('express');
const path = require('path');
const cors = require("cors");
const app = express(),
bodyParser = require("body-parser");
port = 3080;

app.use(bodyParser.json());
app.use(express.static(path.join(__dirname, '../my-app/build')));

app.use(cors());

const users = [
  {
    'first_name': 'Lee',
    'last_name' : 'Earth'
  }
]

app.get('/api/users', (req, res) => {
  console.log('api/users called!')
  res.json(users);
});

app.listen(port, '0.0.0.0', () => {
  console.log(`Server listening on the port::${port}`);
});
```

The `Earthfile` is at the root of the directory.

`./Earthfile`

```Dockerfile
VERSION 0.8
FROM node:13.10.1-alpine3.11
WORKDIR /js-example

app-deps:
    COPY ./app/package.json ./
    COPY ./app/package-lock.json ./
    RUN npm install
    # Output these back in case npm install changes them.
    SAVE ARTIFACT package.json AS LOCAL ./app/package.json
    SAVE ARTIFACT package-lock.json AS LOCAL ./app/package-lock.json

build-app:
    FROM +app-deps
    COPY ./app/src ./app/src
    RUN mkdir -p ./app/dist && cp ./app/src/index.html ./app/dist/
    RUN cd ./app && npx webpack
    SAVE ARTIFACT ./app/dist /dist AS LOCAL ./app/dist

app-docker:
    FROM +app-deps
    ARG tag='latest'
    COPY +build-app/dist ./app/dist
    EXPOSE 8080
    ENTRYPOINT ["/js-example/node_modules/http-server/bin/http-server", "./app/dist"]
    SAVE IMAGE js-example:$tag

api-deps:
    COPY ./api/package.json ./
    COPY ./api/package-lock.json ./
    RUN npm install
    # Output these back in case npm install changes them.
    SAVE ARTIFACT package.json AS LOCAL ./api/package.json
    SAVE ARTIFACT package-lock.json AS LOCAL ./api/package-lock.json

api-docker:
    FROM +api-deps
    ARG tag='latest'
    COPY ./api/server.js .
    RUN pwd
    RUN ls
    EXPOSE 3080
    ENTRYPOINT ["node", "server.js"]
    SAVE IMAGE js-api:$tag

# Run your app and api side by side
app-with-api:
    FROM earthly/dind:alpine-3.19-docker-25.0.5-r0
    RUN apk add curl
    WITH DOCKER \
        --load app:latest=+app-docker \
        --load api:latest=+api-docker
        RUN docker run -d -p 3080:3080 api && \
            docker run -d -p 8080:8080 app  && \
            sleep 5 && \
            curl 0.0.0.0:8080 | grep 'Getting Started' && \
            curl 0.0.0.0:3080/api/users | grep 'Earth'
    END

```
Now you can run `earthly -P +app-with-api` to run the app and api side-by-side.
</details>

<details open>
<summary>Java</summary>

To copy the files for [this example ( Part 6 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/java/part6) run

```bash
mkdir tutorial
cd tutorial
earthly --artifact github.com/earthly/earthly/examples/tutorial/java:main+part6/part6 ./part6
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

with-postgresql:
    FROM earthly/dind:alpine-3.19-docker-25.0.5-r0
    COPY ./docker-compose.yml .
    RUN apk update
    RUN apk add postgresql-client
    WITH DOCKER --compose docker-compose.yml --load app:latest=+docker
        RUN while ! pg_isready --host=localhost --port=5432; do sleep 1; done ;\
            docker run --network=default_java/part6_default app
    END

```

`docker-compose.yml`

```yml
version: "3.9"
   
services:
  db:
    image: postgres
    container_name: db
    hostname: postgres
    environment:
      - POSTGRES_DB=test_db
      - POSTGRES_USER=earthly
      - POSTGRES_PASSWORD=password
    ports:
      - 127.0.0.1:5432:5432
    networks:
      - java/part6_default

networks:
  java/part6_default:


```
The code of the app might look like this

`./src/main/java/hello/HelloWorld.java`

```java

package postgresclient;

import org.joda.time.LocalTime;
import java.sql.Connection;
import java.sql.DriverManager;


public class PostgreSQLJDBC {
   public static void main(String args[]) {
      Connection c = null;
      try {
         Class.forName("org.postgresql.Driver");
         c = DriverManager
            .getConnection("jdbc:postgresql://postgres:5432/test_db",
            "earthly", "password");
      } catch (Exception e) {
         e.printStackTrace();
         System.err.println(e.getClass().getName()+": "+e.getMessage());
         System.exit(0);
      }
      System.out.println("Opened database successfully");
   }
}
```

`./build.gradle`

```groovy
apply plugin: 'java'
apply plugin: 'application'

mainClassName = 'postgresclient.PostgreSQLJDBC'

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
    compile(group: 'org.postgresql', name: 'postgresql', version: '42.3.3')
    testCompile "junit:junit:4.12"
}

```
</details>

<details open>
<summary>Python</summary>

To copy the files for [this example ( Part 6 )](https://github.com/earthly/earthly/tree/main/examples/tutorial/python/part6) run

```bash
earthly --artifact github.com/earthly/earthly/examples/tutorial/python:main+part6/part6 ./part6
```
`./tests/test_db_connection.py`

```python
import unittest
import psycopg2

class MyIntegrationTests(unittest.TestCase):

    def test_db_connection_active(self):
        connection = psycopg2.connect(
            host="postgres",
            database="test_db",
            user="earthly",
            password="password")
        
        self.assertEqual(connection.closed, 0)

if __name__ == '__main__':
    unittest.main()
```

```yml
version: "3.9"
   
services:
  db:
    image: postgres
    container_name: db
    hostname: postgres
    environment:
      - POSTGRES_DB=test_db
      - POSTGRES_USER=earthly
      - POSTGRES_PASSWORD=password
    ports:
      - 5432:5432
    networks:
      - python/part6_default

networks:
  python/part6_default:
```

`./Earthfile`

```Dockerfile
VERSION 0.8
FROM python:3
WORKDIR /code

build:
    COPY ./requirements.txt .
    RUN pip install -r requirements.txt
    COPY . .

run-tests:
    FROM earthly/dind:alpine-3.19-docker-25.0.5-r0
    COPY ./docker-compose.yml .
    COPY ./tests ./tests
    RUN apk update
    RUN apk add postgresql-client
    WITH DOCKER --compose docker-compose.yml --load app:latest=+docker
        RUN while ! pg_isready --host=localhost --port=5432; do sleep 1; done ;\
          docker run --network=default_python/part6_default app python3 ./tests/test_db_connection.py
    END
```
</details>
