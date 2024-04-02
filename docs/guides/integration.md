# Integration Testing With Earthly

Running unit tests in a build pipeline is relatively simple. By definition, unit tests have no external dependencies. Things get more interesting when we want to test how our service integrates with other services and external systems. A service may have dependencies on external file systems, on databases, on external message queues, or other services. An ergonomic and effective development environment should have simple ways to construct and run integration tests. It should be easy to run these tests locally on the developer machine and in the build pipeline. 

** This guide will take an existing application with integration tests and show how they can be easily run inside earthly, both in the local development environment as well as in the build pipeline. **
## Prerequisites 

*This integration approach can work with most applications and development stacks. See [examples](https://github.com/earthly/earthly/tree/main/examples) for guidance on using earthly in other languages.*

### Our Application

The application we start with is simple. It returns the first 5 countries alphabetically via standard out. It has unit tests and integration tests. The integration tests require a datastore with the correct data in place.  

{% method %}
{% sample lang="App" %}
Application code:
```scala
Object Main extends App {
  val dal = new DataAccessLayer()
  val dv = new DataVersion()

  if(dv.version() > 1)
  {
    implicit val cs = IO.contextShift(ExecutionContext.global)
    val xa = Transactor.fromDriverManager[IO](
      "org.postgresql.Driver", 
      "jdbc:postgresql://localhost:5432/iso3166", 
      "postgres",
      "postgres"
    )

    val countries = dal.countries(5)
                       .transact(xa).unsafeRunSync
                       .toList.map(_.name).mkString(", ")

    println(s"The first 5 countries alphabetically are: $countries")
  }
}
```
The output of running the application:

``` bash
> sbt run
The first 5 countries are Afghanistan, Albania, Algeria, American Samoa, Andorra
```

{% sample lang="Unit Test" %}
``` scala
class DataVersionSpec extends FlatSpec {

  val dv = new DataVersion()
  "Data Version " should " be positive" in {
    assert(dv.version > 0)
  }
}
```
Output of running unit tests:
``` bash
> sbt test
[info] DataVersionSpec:
[info] Data Version
[info] - should be positive
[info] Run completed in 810 milliseconds.
[info] Total number of tests run: 1

```

{% sample lang="Integration Test" %}
Integration test:
``` scala
class DatabaseIntegrationTest extends FlatSpec {
  implicit val cs = IO.contextShift(ExecutionContext.global)

  val xa = Transactor.fromDriverManager[IO](
    "org.postgresql.Driver", 
    "jdbc:postgresql://localhost:5432/iso3166", 
    "postgres",
    "postgres"
  )

  "A table" should "have country data" in {
    val dal = new DataAccessLayer()
    assert(dal.countries(5).transact(xa).unsafeRunSync.size == 5)
  }
}
```
Output:
``` bash
>sbt it:test
[info] DatabaseIntegrationTest:
[info] A table
[info] - should have country data
[info] Run completed in 2 seconds, 954 milliseconds.
[info] Total number of tests run: 1
```
{% sample lang="Service Dependencies" %}

The Docker compose configuration specifies the application's dependencies. It is useful for local development and can be started and stopped using `docker-compose up -d` and `docker-compose down`.
This will also be essential for our Earthly integration tests.

Docker Compose:
``` yaml
version: "3"
services:
  postgres:
    container_name: local-postgres
    image: aa8y/postgres-dataset:iso3166
    ports:
      - 5432:5432
    hostname: postgres
    environment:
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
  postgres-ui:
    container_name: local-postgres-ui
    image: adminer:latest
    depends_on:
      - postgres
    ports:
      - 8080:8080
    hostname: postgres-ui
```

{% endmethod %}

### The Basic Earthfile

We start with a simple Earthfile that can build and create a docker image for our app. See the [Basics guide](../basics/basics.md) for more details, as well as examples in many programming languages.

{% method %}
{% sample lang="Base Earthly Target" %}

We start from an appropriate docker image and set up a working directory. 
``` Dockerfile
VERSION 0.8
FROM earthly/dind:alpine-3.19-docker-25.0.3-r1
WORKDIR /scala-example
RUN apk add openjdk11 bash wget postgresql-client
```

[Full file](https://github.com/earthly/earthly-example-scala/blob/main/integration/Earthfile)

{% sample lang="Project Files" %}
We then install SBT

``` Dockerfile
sbt: 
    #Scala
    # Defaults if not specified
    ARG sbt_version=1.3.2
    ARG sbt_home=/usr/local/sbt

    # Download and extract from archive
    RUN mkdir -pv "$sbt_home"
    RUN wget -qO - "https://github.com/sbt/sbt/releases/download/v$sbt_version/sbt-$sbt_version.tgz" >/tmp/sbt.tgz
    RUN tar xzf /tmp/sbt.tgz -C "$sbt_home" --strip-components=1
    RUN ln -sv "$sbt_home"/bin/sbt /usr/bin/

    # This triggers a bunch of useful downloads.
    RUN sbt sbtVersion
```

We then copy in our build files and run Scala Build Tool, so that we can cache our dependencies

``` Dockerfile
project-files:
    FROM +sbt
    COPY build.sbt ./
    COPY project project
    # Run sbt for caching purposes.
    RUN touch a.scala && sbt compile && rm a.scala
```

<!-- due to gitbook bug, https://github.com/earthly/earthly/blob/main/examples/integration-test/Earthfile changed to https://tinyurl.com/4m6hbd6a -->

[Full file](https://tinyurl.com/4m6hbd6a)

{% sample lang="Compile" %}

We also set up our build target.
``` Dockerfile
build:
    FROM +project-files
    COPY src src
    RUN sbt compile
```
[Full file](https://tinyurl.com/4m6hbd6a)

{% sample lang="Unit Test" %}

For unit tests, we copy in the source and run the tests.

``` Dockerfile

unit-test:
    FROM +project-file
    COPY src src
    RUN sbt test

```
[Full file](https://tinyurl.com/4m6hbd6a)

{% sample lang="Docker" %}

We then build a Dockerfile.

``` Dockerfile
docker:
    FROM +project-file
    COPY src src
    RUN sbt assembly
    ENTRYPOINT ["java","-cp","target/scala-2.12/scala-example-assembly-1.0.jar","Main"]
    SAVE IMAGE scala-example:latest 
```
[Full file](https://github.com/earthly/earthly-example-scala/blob/main/integration/Earthfile)

{% endmethod %}

See the [Basics Guide](../basics/basics.md) for more details on these steps, including how they might differ in Go, JavaScript, Java, and Python.

## In-App Integration Testing 

Since our service has a docker-compose file of dependencies, running integration tests is easy.

Our integration target needs to copy in our source code and our Dockerfile and then inside a `WITH DOCKER` start the tests:
``` Dockerfile
integration-test:
    FROM +project-files
    COPY src src
    COPY docker-compose.yml ./ 
    WITH DOCKER --compose docker-compose.yml
        RUN while ! pg_isready --host=localhost --port=5432 --dbname=iso3166 --username=postgres; do sleep 1; done ;\
            sbt it:test
    END
```
The `WITH DOCKER` has a `--compose` flag that we use to start up our docker-compose and run our integration tests in that context.

We can now run our tests both locally and in the CI pipeline, in a reproducible way:

``` bash
> earthly -P +integration-test
+integration-test | Creating local-postgres ... done
+integration-test | Creating local-postgres-ui ... done
+integration-test | +integration-test | [info] Loading settings for project scala-example-build from plugins.sbt ...
+integration-test | [info] DatabaseIntegrationTest:
+integration-test | [info] A table
+integration-test | [info] - should have country data
+integration-test | [info] Run completed in 7 seconds, 923 milliseconds.
+integration-test | [info] Tests: succeeded 1, failed 0, canceled 0, ignored 0, pending 0
+integration-test | Stopping local-postgres-ui ... done
+integration-test | Stopping local-postgres    ... done
+integration-test | Removing local-postgres-ui ... done
+integration-test | Removing local-postgres    ... done
+integration-test | Removing network scala-example_default
+integration-test | Target github.com/earthly/earthly-example-scala/integration:main+integration-test built successfully
...
```
This means that if an integration test fails in the build pipeline, you can easily reproduce it locally.  

## End to End Integration Tests

Our first integration test used was part of the service we were testing. This is one way to exercise integration code paths. Another useful form of integration testing is end-to-end testing. In this form of integration testing, we start up the application and test it from the outside. 

In our simplified case example, with a single code path, a test that verifies the application starts and produces the desired output is sufficient. 

{% method %}
{% sample lang="Test Script" %}
``` bash
source "./assert.sh"
set -v
results=$(docker run --network=host earthly/examples:integration)
expected="The first 5 countries alphabetically are: Afghanistan, Albania, Algeria, American Samoa, Andorra"

assert_eq "$expected" "$results"n
```

{% sample lang="Earth File" %}
``` dockerfile
smoke-test:
    FROM +project-files
    COPY docker-compose.yml ./ 
    COPY src/smoketest ./ 
    WITH DOCKER --compose docker-compose.yml --load=+docker
        RUN while ! pg_isready --host=localhost --port=5432 --dbname=iso3166 --username=postgres; do sleep 1; done ;\
            ./smoketest.sh
    END
```
{% endmethod %}

Output:
We can then run this and check that our application with its dependencies, produces the correct output.

``` Dockerfile
> earthly -P +smoke-test
+smoke-test | --> WITH DOCKER RUN for i in {1..30}; do nc -z localhost 5432 && break; sleep 1; done; docker run --network=host earthly/examples:integration
+smoke-test | Loading images...
+smoke-test | Loaded image: aa8y/postgres-dataset:iso3166
+smoke-test | Loaded image: adminer:latest
+smoke-test | Loaded image: earthly/examples:integration
+smoke-test | ...done
+smoke-test | Creating network "scala-example_default" with the default driver
+smoke-test | Creating local-postgres ... done
+smoke-test | Creating local-postgres-ui ... done
+smoke-test | +smoke-test | The first 5 countries alphabetically are: Afghanistan, Albania, Algeria, American Samoa, Andorra
+smoke-test | Stopping local-postgres-ui ... done
+smoke-test | Stopping local-postgres    ... done
+smoke-test | Removing local-postgres-ui ... done
+smoke-test | Removing local-postgres    ... done
+smoke-test | Removing network scala-example_default
+smoke-test | Target github.com/earthly/earthly-example-scala/integration:main+smoke-test built successfully
=========================== SUCCESS ===========================
...
```

## Bringing It All Together

Adding these testing targets to an all target, we now can unit test, integration test, and dockerize and push our software in a single command. Using this approach, integration tests that fail sporadically for environmental reasons and can't be reproduced consistently should be a thing of the past.

``` Dockerfile
all:
  BUILD +build
  BUILD +unit-test
  BUILD +integration-test
  BUILD +smoke-test
```

``` bash
> earthly -P +all
...
+all | Target github.com/earthly/earthly-example-scala/integration:main+all built successfully
=========================== SUCCESS ===========================
```

There we have it, a reproducible integration process. If you have questions about the example, [ask](https://gitter.im/earthly-room/community)

## See also
* [Docker In Earthly](./docker-in-earthly.md)
* [Source code for example](https://github.com/earthly/earthly/tree/main/examples/integration-test)
* [Integration Testing vs Unit Testing](https://blog.earthly.dev/unit-vs-integration/)
