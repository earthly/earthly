# Integration Testing With Earthly

Running unit tests in a build pipeline is relatively simple. By definition, unit tests have no external dependencies. Things get more interesting when we want to test how our service integrates with other services and external systems. A service may have dependencies on external file systems, on databases, on external message queues, or other services. An ergonomic and effective development environment should have simple ways to construct and run integration tests. It should be easy to run these tests locally on the developer machine and in the build pipeline. 

## A Scenario

Imagine this scenario: You work on a team that owns a service, that, in production runs containerized as part of a microservice architecture. The totality of the services that make-up production are too numerous reasonably run on a developer machine but thankfully your service only depends on a handful of other services. You have unit tests but need ways to test your points of integration with your dependencies. You also want ways to test your service end to end, that is you want to exercise it in ways a service consumer might. 

For simplicity's sake, this guide will use a Scala application with a very simple set of dependencies. In principle, these steps will work with any service whose dependencies can be specified in a docker-compose file.

## The Example App

The app has one purpose, it returns the first 5 countries alphabetically via standard out. To do this it connects to a database through a data access layer and returns the results from a Postgres table. The specifics of the app are not relevant to the testing method but are outlined here for clarity.

The application has unit tests that don't require any dependencies. Additionally it has integration tests that exercise the data access layer and the connection to the database. 

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

  "An table" should "have country data" in {
    val dal = new DataAccessLayer()
    assert(dal.countries(5).transact(xa).unsafeRunSync.size == 5)
  }
}
```
Ouput:
``` bash
>sbt it:test
[info] DatabaseIntegrationTest:
[info] An table
[info] - should have country data
[info] Run completed in 2 seconds, 954 milliseconds.
[info] Total number of tests run: 1
```
{% endmethod %}

## The Earthfile

We start with a simple earthly file that can build and create a docker file for our app. See [Basic](./basics) guide for more details on that.

{% hint style='info' %}
##### Note
This guide assumes you are using a docker image based on [docker:dind](https://hub.docker.com/_/docker) or have installed docker in docker into your base container. It is possible in the future that earthly will not require the dind for `WITH DOCKER` commands and at the time any base image will work successfully. 
{% endhint %}


{% method %}
{% sample lang="Base Earth Target" %}

We start from an alpine docker in docker image and the dependencies we need to build and test our app. These include the jdk and docker-compose. 
``` Dockerfile
FROM docker:19.03.7-dind
    
WORKDIR /scala-example
RUN apk add openjdk11 docker-compose bash wget
```

[Full file](https://github.com/earthly/earthly-example-scala/blob/master/integration/Earthfile)

{% sample lang="SBT" %}

We then install SBT, Scala Build Tool, for building our application.

``` Dockerfile
sbt: 
    #Scala
    # Defaults if not specified in --build-arg
    ARG sbt_version=1.3.2
    ARG sbt_home=/usr/local/sbt

    # Download and extract from archive
    RUN mkdir -pv "$sbt_home"
    RUN wget -qO - "https://github.com/sbt/sbt/releases/download/v$sbt_version/sbt-$sbt_version.tgz" >/tmp/sbt.tgz
    RUN tar xzf /tmp/sbt.tgz -C "$sbt_home" --strip-components=1
    RUN ln -sv "$sbt_home"/bin/sbt /usr/bin/

    # This triggers a bunch of useful downloads.
    RUN sbt sbtVersion
    SAVE IMAGE 
```

[Full file](https://github.com/earthly/earthly/blob/master/examples/integration/Earthfile)

{% sample lang="Build Target" %}

We also copy in our project files and setup our build target.  
``` Dockerfile
project-files:
    FROM +sbt
    COPY build.sbt ./
    COPY project project
    # Run sbt for caching purposes.
    RUN touch a.scala && sbt compile && rm a.scala
    SAVE IMAGE

build:
    FROM +project-files
    COPY src src
    RUN sbt compile
    SAVE IMAGE 
```
[Full file](https://github.com/earthly/earthly/blob/master/examples/integration/Earthfile)

{% sample lang="Unit Test Target" %}

For unit tests, we copy in the source and run the tests.

``` Dockerfile

unit-test:
    FROM +build
    COPY src src
    RUN sbt test

```
[Full file](https://github.com/earthly/earthly/blob/master/examples/integration/Earthfile)

{% sample lang="Docker Target" %}

A Dockerfile is built using the output of `sbt assembly`. 

{% hint style='info' %}
If minimal docker image size is important, other approaches, such as `sbt package` should be considered.
{% endhint %}
``` Dockerfile
docker:
    FROM +build
    COPY src src
    RUN sbt assembly
    ENTRYPOINT ["java","-cp","target/scala-2.12/scala-example-assembly-1.0.jar","Main"]
    SAVE IMAGE scala-example:lates 
```
[Full file](https://github.com/earthly/earthly-example-scala/blob/master/integration/Earthfile)

{% endmethod %}

 See [Basics Guide](./basics.md) for more details on these steps, including how they might differ in Go, Javascript, Java and Python.


## Integration Testing Step 1 - Define Your Dependencies

The first step of integration testing in earthly is to define all the service dependencies in a docker-compose file. 

#### Docker Compose:
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

This docker-compose file is essential for our integration testing, but it is also useful for local development. It can be started and stopped using `docker-compose up -d` and `docker-compose down`.


{% hint style='info' %}
#### SQL Data

In a real-world example, we would likely use a tool like [Flyway](https://flywaydb.org/) to manage our database structure and base data. For simplicity here, however, we are using a Postgres image with country data included `aa8y/postgres-dataset:iso3166`.
{% endhint %}

{% hint style='info' %}
#### Adminer
Our example application has one direct dependency, Postgres, and a helper app `adminer` which is a web UI for postgres. It is strictly not nessary for running our integration tests but useful in local development and serves to show how this solution scales to multiple dependencies.
{% endhint %}

## Integration Testing Step 2 - Run your tests

With our docker-compose ready, we can now add an integration test step to our Earthfile. 

Our integration target needs to copy in our source code and our docker file before it starts the tests:
``` Dockerfile
integration-test:
    FROM +build
    COPY src src
    COPY docker-compose.yml ./ 
```
Next, we use the `WITH DOCKER` statement to start up the docker daemon in our build context.

```
   WITH DOCKER 
   ...
   END
```
Following that, we pull each of our images into the build context using `DOCKER PULL`. While strictly not necessary, using the `DOCKER PULL` command ensures our image pulls are cached by earthly and ensures faster builds.

```
 DOCKER PULL aa8y/postgres-dataset:iso3166
 DOCKER PULL adminer:latest 
```
To run our integration tests, we now start up our docker-compose, wait for it to start up, run our test, and then stop it. We do this in a single run command. 

```
   RUN docker-compose up -d && \
            for i in {1..30}; do nc -z localhost 5432 && break; sleep 1; done; \
            sbt it:test && \
            docker-compose down 
```
{% hint style='info' %}
#### About netcat (nc)

This statement is a simple loop, that will block for up to 30 seconds or until we can read from port 5432 on localhost. 
```
for i in {1..30}; do nc -z localhost 5432 && break; sleep 1; done; \
```

Our application will connect to Postgres via localhost:5432. This step, therefore, ensures that don't run our tests until the database is up. There are many other ways to accomplish this, including READY checks, application-specific code, and scripts like [wait for it](https://github.com/vishnubob/wait-for-it). 

 Coordinating in among services is a complicated area out of the scope of this guide.


{% endhint %}


### Combined

Putting this all together we get:

``` Dockerfile
integration-test:
    FROM +build
    COPY src src
    COPY docker-compose.yml ./ 
    WITH DOCKER 
        DOCKER PULL aa8y/postgres-dataset:iso3166
        DOCKER PULL adminer:latest
        RUN docker-compose up -d && \
            for i in {1..30}; do nc -z localhost 5432 && break; sleep 1; done; \
            sbt it:test && \
            docker-compose down 
    END

```
We can now run our it tests both locally and in the CI pipeline, in a reproducible way:

``` bash
> earth -P +integration-test
+integration-test | Creating local-postgres ... done
+integration-test | Creating local-postgres-ui ... done
+integration-test | +integration-test | [info] Loading settings for project scala-example-build from plugins.sbt ...
+integration-test | [info] DatabaseIntegrationTest:
+integration-test | [info] An table
+integration-test | [info] - should have country data
+integration-test | [info] Run completed in 7 seconds, 923 milliseconds.
+integration-test | [info] Tests: succeeded 1, failed 0, canceled 0, ignored 0, pending 0
+integration-test | Stopping local-postgres-ui ... done
+integration-test | Stopping local-postgres    ... done
+integration-test | Removing local-postgres-ui ... done
+integration-test | Removing local-postgres    ... done
+integration-test | Removing network scala-example_default
+integration-test | Target github.com/earthly/earthly-example-scala/integration:master+integration-test built successfully
...
```
This means that if an integration test fails in the build pipeline, you can easily reproduce it locally.  


## Running End to End Integration Tests

Our first integration test run used a testing harness inside the service under test. This is one way to exercise integration code paths and could be called whitebox integration testing. Another useful form of integration testing is end to end integration testing. In this form of integration testing, we start up the application and test it from the outside. 

In our simplified case example, with a single code path, a smoke test is sufficient. We start up the application, with its dependencies, and verify it runs successfully.


```
        RUN docker-compose up -d && \ 
            for i in {1..30}; do nc -z localhost 5432 && break; sleep 1; done; \
            docker run --network=host scala-example:latest && \
            docker-compose down
```
{% hint style='info' %}
#### Docker Networking
Note the `-network=host` flag passed to `docker run`. 
```
docker run --network=host scala-example:latest 
```
This tells docker to share the host network with this container, allowing it to access docker-compose ports using localhost.

{% endhint %}

Full Example:
``` Dockerfile
smoke-test:
    FROM +base
    COPY docker-compose.yml ./ 
    WITH DOCKER
        DOCKER PULL aa8y/postgres-dataset:iso3166
        DOCKER PULL adminer:latest
        DOCKER LOAD +docker scala-example:latest
        RUN docker-compose up -d && \ 
            for i in {1..30}; do nc -z localhost 5432 && break; sleep 1; done; \
            docker run --network=host scala-example:latest && \
            docker-compose down 
    END
```

Output:
``` Dockerfile
> earth -P +smoke-test
+smoke-test | --> WITH DOCKER RUN docker-compose up -d && for i in {1..30}; do nc -z localhost 5432 && break; sleep 1; done; docker run --network=host scala-example:latest && docker-compose down
+smoke-test | Loading images...
+smoke-test | Loaded image: aa8y/postgres-dataset:iso3166
+smoke-test | Loaded image: adminer:latest
+smoke-test | Loaded image: scala-example:latest
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
+smoke-test | Target github.com/earthly/earthly-example-scala/integration:master+smoke-test built successfully
=========================== SUCCESS ===========================
...
```

In more complex scenarios, this example could be extended to run tests against the service under test. Making http calls and verifying outputs using your preferred testing framework.

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
> earth -P +all
...
+all | Target github.com/earthly/earthly-example-scala/integration:master+all built successfully
=========================== SUCCESS ===========================
```

There we have it, a reproducible integration process. The full example can be found [here](). If you have questions about the example, [ask them here](https://gitter.im/earthly-room/community)

## See also
* [Docker In Earthly](./docker-in-earthly.md)
* [Source code for example](https://github.com/earthly/earthly/tree/master/examples/integration)
