
# Using Docker in Earthly

This guide walks through using Docker commands in Earthly as well as some possible use-cases, such as integration testing.

## Basic usage

In order to use Docker commands (such as `docker run`), Earthly makes available isolated Docker daemons which are started and stopped on-demand. The reason for using isolated instances of Docker daemons is such that no pre-existing Docker state (eg images, containers, networks, volumes) can influence the way the build executes, thus having no way to create inconsistencies when the same build is run in a different environment.

Here is a quick example of running a hello-world docker container via `docker run` in Earthly:

```Dockerfile
hello:
    FROM docker:19.03.13-dind
    WITH DOCKER
        DOCKER PULL hello-world
        RUN docker run hello-world
    END
```

Let's break it down.

`FROM docker:19.03.13-dind` inherits from an official docker-in-docker (dind) image. This is mandatory, because `WITH DOCKER` requires all the Docker binaries (not just the client) to be present in the build environment.

`WITH DOCKER ... END` starts a Docker daemon for the purpose of running Docker commands against it. At the end of the execution, this also terminates the daemon and permanently deletes all of its data (eg daemon cached images).

`DOCKER PULL hello-world` pulls the image `hello-world` from the Docker Hub. This command could have been replaced with the more traditional `docker pull hello-world`. However, the Earthly variant additionally stores the image in the Earthly cache, so that the actual pull is performed only if the images changes. Because the daemon cache is cleared after each run, `docker pull` would not achieve the same.

`RUN docker run hello-world` executes the `docker run` command in the context of the daemon created by `WITH DOCKER`.

## Loading images built by Earthly

A typical use of Docker in Earthly is running an image that has been built via Earthly itself. To achieve that, the command `DOCKER LOAD` can be used. Here is an example:

```Dockerfile
build:
    ...
    ENTRYPOINT ...
    SAVE IMAGE my-image:latest

smoke-test:
    FROM docker:19.03.13-dind
    WITH DOCKER
        DOCKER LOAD +build test:latest
        RUN docker run test:latest
    END
```

`DOCKER LOAD +build test:latest` takes the image produced by the target `+build` and loads it into the Docker daemon created by `WITH DOCKER` as the image with the tag `test:latest`. The tag can then be used to reference this image in other docker commands, such as `docker run`.

Notice that although the image name produced as output is `my-image:latest`. This image name is not available in the `WITH DOCKER` environment, however, as it is only used to tag for use outside of Earthly.

<!--
## Integration testing

TODO: Link to the integration testing example (docker-compose example in Earthly) when ready.
-->

## Limitations of Docker in Earthly

The current implementation of Docker in Earthly has a number of limitations:

* Only one `RUN` command is allowed within the `WITH DOCKER ... END` clause. The reason for this is that only one cache layer is used for the entire `WITH DOCKER ... END` clause. You can, however, chain multiple shell commands together within a single `RUN` command:
  ```Dockerfile
  WITH DOCKER
      RUN command1 && \
          command2 && \
          command3 && \
          ...
  END
  ```
* The target containing the `WITH DOCKER` clause has to inherit from an official Docker-in-Docker (dind) image such as `docker:19.03.13-dind`. If your build requires the use of an alternative environment as part of a test (eg to run commands like `sbt test` or `go test` together with a docker-compose stack), consider placing the test itself in a Docker image, then loading that image via `DOCKER LOAD` and running the test as a Docker container. (This limitation may be lifted in the future.)
* To maximize the use of cache, all external images used should be declared via the command `DOCKER PULL`. Even though commands such as `docker run` automatically pull an image if it is not found locally, it will do so every single time the `WITH DOCKER` clause is executed, due to Docker caching not being preserved between runs. Pre-declaring the images via `DOCKER PULL` ensures that they are properly cached by Earthly to minimize unnecessary redownloads.
* `docker-compose` needs to be installed separately. This can be achieved easily, however, via `apk add docker-compose` placed just after the `FROM` command:
  ```Dockerfile
  FROM docker:19.03.13-dind
  RUN apk --update --no-cache add docker-compose
  WITH DOCKER
      ...
  END
  ```
* `docker build` cannot be used to build Dockerfiles. However, the Earthly command `FROM DOCKERFILE` can be used instead. See [alternative to docker build](#alternative-to-docker-build) below.
* The state of the Docker daemon within Earthly cannot be inspected on the host (eg for debugging purposes). For example, if a `docker-compose` stack fails, you cannot execute commands like `docker-compose logs` or `docker logs` on the host. However, you may use the interactive mode to drop into a shell within the build environment and execute such commands there. For more information, see the [debugging guide](./debugging.md).

## Alternatives to Docker in Earthly

It is not always necessary to execute docker commands within an Earthly build. Certain operations can be replicated with Earthly constructs.

### Alternative to docker run

In certain cases, simple `docker run` invocations can be replaced by a simple [`RUN --entrypoint`](../earthfile/earthfile.md#entrypoint). For example, the following:

```Dockerfile
FROM docker:19.03.13-dind
WITH DOCKER
    DOCKER PULL hello-world
    RUN docker run hello-world
END
```

Can be rewritten as

```Dockerfile
FROM hello-world
RUN --entrypoint
```

This, of course, has limitations, such as not being able to mount volumes the same way `docker run -v ...`  could (instead, a `COPY` command could be used); or not being able to run multiple containers in parallel. However, when appropriate, it can simplify a build definition.

### Alternative to docker build

Running `docker build` within Earthly is discouraged, as it has a number of key limitaitons:

* Layer caching does not work. This is because `WITH DOCKER` does not preserve Docker cache between runs (other than `DOCKER PULL`).
* Once an image is created, it cannot be exported as a build output in a form other than a TAR archive (eg it cannot be loaded onto the host Docker daemon).

Instead of executing `docker build`, it is advisable to use the Earthly command `FROM DOCKERFILE`. For example, the command `docker build -t my-image:latest .` can be emulated by:

```Dockerfile
FROM DOCKERFILE .
SAVE IMAGE my-image:latest
```

## See also

* Reference for [`WITH DOCKER`](../earthfile/earthfile.md#with-docker-beta)
* Reference for [`DOCKER PULL`](../earthfile/earthfile.md#docker-pull-beta)
* Reference for [`DOCKER LOAD`](../earthfile/earthfile.md#docker-load-beta)
* [Debugging techniques](./debugging.md)
