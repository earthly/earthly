# C++ example

A complete C++ example is available on [the Basics page](../guides/basics.md).

```Dockerfile
# build.earth
FROM ubuntu:20.10

# for apt to be noninteractive
ENV DEBIAN_FRONTEND noninteractive
ENV DEBCONF_NONINTERACTIVE_SEEN true

RUN apt-get update && apt-get install -y build-essential cmake

WORKDIR /code

code:
  COPY src src
  SAVE IMAGE

build:
  FROM +code
  RUN cmake src
  # cache cmake temp files to prevent rebuilding .o files when the .cpp files don't change
  RUN --mount=type=cache,target=/code/CMakeFiles make
  SAVE ARTIFACT hello AS LOCAL "hello"

docker:
  COPY +build/hello /bin/hello
  ENTRYPOINT ["/bin/hello"]
  SAVE IMAGE cpp-example:latest

```

For the complete code see the [examples/cpp GitHub directory](https://github.com/earthly/earthly/tree/master/examples/cpp).
