
# Using Docker in Earthly

This guide walks through using Docker commands in Earthly.

## Basic usage

In order to use Docker commands (such as `docker run`), Earthly makes available isolated Docker daemons which are started and stopped on-demand. The reason for using isolated instances of Docker daemons is such that no preexisting Docker state (e.g. images, containers, networks, volumes) can influence the way the build executes. This allows Earthly to achieve high degrees of reproducibility.

Here is a quick example of running a `hello-world` docker container via `docker run` in Earthly:

```Dockerfile
hello:
    FROM earthly/dind:alpine-3.19-docker-25.0.5-r0
    WITH DOCKER --pull hello-world
        RUN docker run hello-world
    END
```

Let's break it down.

`FROM earthly/dind:alpine-3.19-docker-25.0.5-r0` inherits from an Earthly-supported docker-in-docker (dind) image. This is recommended, because `WITH DOCKER` requires all the Docker binaries (not just the client) to be present in the build environment.

`WITH DOCKER ... END` starts a Docker daemon for the purpose of running Docker commands against it. At the end of the execution, this also terminates the daemon and permanently deletes all of its data (e.g. daemon cached images).

`--pull hello-world` pulls the image `hello-world` from the Docker Hub. This option could have been replaced with the more traditional `docker pull hello-world`. However, the Earthly variant additionally stores the image in the Earthly cache, so that the actual pull is performed only if the image changes. Because the daemon cache is cleared after each run, `docker pull` would not achieve the same.

`RUN docker run hello-world` executes the `docker run` command in the context of the daemon created by `WITH DOCKER`.

## Loading images built by Earthly

A typical use of Docker in Earthly is running an image that has been built via Earthly itself. To achieve that, the option `WITH DOCKER --load ...=...` can be used. Here is an example:

```Dockerfile
build:
    ...
    ENTRYPOINT ...
    SAVE IMAGE my-image:latest

smoke-test:
    FROM earthly/dind:alpine-3.19-docker-25.0.5-r0
    WITH DOCKER --load test:latest=+build
        RUN docker run test:latest
    FROM earthly/dind:alpine
    WITH DOCKER --load +build
        RUN docker run my-image:latest
    END
```

`--load +build` takes the image produced by the target `+build` and loads it into the Docker daemon created by `WITH DOCKER`.

## Running docker-compose

It is possible to run `docker-compose` via `WITH DOCKER`, either explicitly, simply by running the `docker-compose` tool, or implicitly, via the `--compose` flag. The `--compose` flag allows you to specify a Docker compose stack that needs to be brought up before the execution of the `RUN` command. For example:

```Dockerfile
FROM earthly/dind:alpine-3.19-docker-25.0.5-r0
COPY docker-compose.yml ./
WITH DOCKER \
        --compose docker-compose.yml \
        --service db \
        --service api
    RUN docker run some-integration-test:latest
END
```

Using the `--compose` flag has the added benefit that any images needed by the compose stack will be automatically added to the pull list by Earthly, thus using cache efficiently.

## Performance

It's recommended to use the `earthly/dind:alpine-3.19-docker-25.0.5-r0` image for running docker-in-docker. See the best-practices' section on using [with docker](../guides/best-practices.md#use-earthly-dind) for more details.

In cases when using `earthly/dind` is not possible, Earthly will attempt to install Docker in the image you have chosen. This has the drawback of not being able to use cache efficiently and is not recommended for performance reasons.

Another option is to use the Earthly UDC `INSTALL_DIND`. This will install Docker in the build environment, but you can control at what point in the build it happens, thus being able to optimize your cache use. For example:

```Dockerfile
FROM my-image:latest
DO github.com/earthly/lib+INSTALL_DIND
COPY ./docker-compose.yml ./
WITH DOCKER ...
    ...
END
```

In the above example, the `INSTALL_DIND` command is executed before the `COPY` command, thus ensuring that the Docker installation is cached and reused for subsequent builds when the file `docker-compose.yml` changes.

## Integration testing

For more information on integration testing and working with service dependencies see our [tutorial on integration testing in Earthly](./integration.md).

## Limitations of Docker in Earthly

The current implementation of Docker in Earthly has a number of limitations:

* Only one `RUN` command is allowed within the `WITH DOCKER` clause. The reason for this is that only one cache layer is used for the entire clause. You can, however, chain multiple shell commands together within a single `RUN` command. For example:
  ```Dockerfile
  WITH DOCKER
      RUN command1 && \
          command2 && \
          command3 && \
          ...
  END
  ```
* It is recommended that the target containing the `WITH DOCKER` clause inherits from a supported Docker-in-Docker (dind) image such as `earthly/dind:alpine-3.19-docker-25.0.5-r0` or `earthly/dind:ubuntu-23.04-docker-25.0.2-1`. If your build requires the use of an alternative environment as part of a test (e.g. to run commands like `sbt test` or `go test` together with a docker-compose stack), consider placing the test itself in a Docker image, then loading that image via `--load` and running the test as a Docker container.
* If you do not use an officially supported Docker-in-Docker image, Earthly will attempt to install Docker in whatever image you have chosen. This has the drawback of not being able to use cache efficiently and is not recommended for performance reasons.
* To maximize the use of cache, all external images used should be declared via the options `--pull` or `--compose`. Even though commands such as `docker run` automatically pull an image if it is not found locally, it will do so every single time the `WITH DOCKER` clause is executed, due to Docker caching not being preserved between runs. Pre-declaring the images ensures that they are properly cached by Earthly to minimize unnecessary redownloads.
* `docker build` cannot be used to build Dockerfiles. However, the Earthly command `FROM DOCKERFILE` can be used instead. See [alternative to docker build](#alternative-to-docker-build) below.
* The state of the Docker daemon within Earthly cannot be inspected on the host (e.g. for debugging purposes). For example, if a `docker-compose` stack fails, you cannot execute commands like `docker-compose logs` or `docker logs` on the host. However, you may use the interactive mode to drop into a shell within the build environment and execute such commands there. For more information, see the [debugging guide](./debugging.md).
* It is currently not possible to mount `/var/run/docker.sock` in order to use the host Docker daemon. This goes against Earthly's principles of keeping execution repeatable. Mounting the Docker socket may cause builds to depend on the host Daemon state (e.g. pre-cached images) in ways that may not be obvious or easy to reproduce if the build were executed in another environment.
* Creating external networks (the default docker behavior), will sometimes fail and produce 100% packet loss; however __internal__ networks do no have this issue. Use `docker network create --internal ...` or the `internal: true` compose option instead. See [#3495](https://github.com/earthly/earthly/issues/3495) for details.

## Alternatives to Docker in Earthly

It is not always necessary to execute docker commands within an Earthly build. Certain operations can be replicated with Earthly constructs.

### Alternative to docker run

In certain cases, simple `docker run` invocations can be replaced by a simple [`RUN --entrypoint`](../earthfile/earthfile.md#entrypoint). For example, the following:

```Dockerfile
FROM docker:19.03.13-dind
WITH DOCKER --pull hello-world
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

Running `docker build` within Earthly is discouraged, as it has a number of key limitations:

* Layer caching does not work. This is because `WITH DOCKER` does not preserve Docker cache between runs (other than `--pull`).
* Once an image is created, it cannot be exported as a build output in a form other than a TAR archive (e.g. it cannot be automatically loaded onto the host Docker daemon).

Instead of executing `docker build`, it is advisable to use the Earthly command `FROM DOCKERFILE`. For example, the command `docker build -t my-image:latest .` can be emulated by:

```Dockerfile
FROM DOCKERFILE .
SAVE IMAGE my-image:latest
```

## See also

* Reference for [`WITH DOCKER`](../earthfile/earthfile.md#with-docker)
* [Debugging techniques](./debugging.md)
* [Tutorial on integration testing in Earthly](./integration.md)
