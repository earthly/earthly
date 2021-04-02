
# Detailed explanation

Going back to the example earthfile definition, here is what each command does:

{% method %}
{% sample lang="Go" %}
```Dockerfile
# Earthfile

# The build starts from a docker image: golang:1.15-alpine3.13
FROM golang:1.15-alpine3.13
# We change the current working directory.
WORKDIR /go-example

# The above commands are inherited implicitly by all targets below
# (as if they started with FROM +base).

# Declare a target, build.
build:
    # Define the recipe of the target build as follows:

    # Copy main.go from the build context to the build environment.
    COPY main.go .
    # Run a go build command.
    # This uses the previously copied main.go file.
    RUN go build -o build/go-example main.go
    # Save the output of the build command as an artifact. Call this
    # artifact /go-example (it can be later referenced as +build/go-example).
    # In addition, store the artifact as a local file (on the host) named
    # build/go-example. This local file is only written if the entire build
    # succeeds.
    SAVE ARTIFACT build/go-example /go-example AS LOCAL build/go-example

# Declare a target, docker.
docker:
    # Define the recipe of the target docker as follows:

    # Copy the artifact /go-example produced by another target, +build, to the
    # current directory within the build container.
    COPY +build/go-example .
    # Set the entrypoint for the resulting docker image.
    ENTRYPOINT ["/go-example/go-example"]
    # Save the current state as a docker image, which will have the docker tag
    # go-example:latest. This image is only made available to the host's docker
    # if the entire build succeeds.
    SAVE IMAGE go-example:latest
```
{% sample lang="JavaScript" %}
```Dockerfile
# Earthfile

# The build starts from a docker image: node:13.10.1-alpine3.11
FROM node:13.10.1-alpine3.11
# We change the current working directory.
WORKDIR /js-example

# The above commands are inherited implicitly by all targets below
# (as if they started with FROM +base).

# Declare a target, build.
build:
    # Define the recipe of the target build as follows:

    # Copy index.js from the build context to the build environment.
    COPY index.js .
    # Save the index.js in an artifact dir called dist (it can be later
    # referenced as +build/dist). In addition, store the artifact as a
    # local file (on the host) named dist/index.js. This local file is only
    # written if the entire build succeeds.
    SAVE ARTIFACT index.js /dist/index.js AS LOCAL ./dist/index.js

# Declare a target, docker.
docker:
    # Define the recipe of the target docker as follows:

    # Copy the artifact /dist produced by another target, +build, to the
    # current directory within the build container.
    COPY +build/dist dist
    # Set the entrypoint for the resulting docker image.
    ENTRYPOINT ["node", "./dist/index.js"]
    # Save the current state as a docker image, which will have the docker tag
    # js-example:latest. This image is only made available to the host's docker
    # if the entire build succeeds.
    SAVE IMAGE js-example:latest
```
{% sample lang="Java" %}
```Dockerfile
# Earthfile

# The build starts from a docker image: openjdk:8-jdk-alpine
FROM openjdk:8-jdk-alpine
# We install gradle using alpine's apk command.
RUN apk add --update --no-cache gradle
# We change the current working directory.
WORKDIR /java-example

# The above commands are inherited implicitly by all targets below
# (as if they started with FROM +base).

# Declare a target, build.
build:
    # Define the recipe of the target build as follows:

    # Copy build.gradle and src from the build context to the build
    # environment.
    COPY build.gradle ./
    COPY src src
    # Run the gradle build and gradle install commands.
    # These use the previously copied src dir.
    RUN gradle build
    RUN gradle install
    # Save the output of the build command as artifacts. Call these
    # artifacts bin and lib (they can be later referenced as +build/bin and
    # +build/lib respectively).
    # In addition, store the artifacts as local dirs (on the host) named
    # build/bin and build/lib. These local dirs are only written if the entire
    # build succeeds.
    SAVE ARTIFACT build/install/java-example/bin /bin AS LOCAL build/bin
    SAVE ARTIFACT build/install/java-example/lib /lib AS LOCAL build/lib

# Declare a target, docker.
docker:
    # Define the recipe of the target docker as follows:

    # Copy the artifacts /bin and /lib produced by another target, +build, to
    # the current directory within the build container.
    COPY +build/bin bin
    COPY +build/lib lib
    # Set the entrypoint for the resulting docker image.
    ENTRYPOINT ["/java-example/bin/java-example"]
    # Save the current state as a docker image, which will have the docker tag
    # java-example:latest. This image is only made available to the host's
    # docker if the entire build succeeds.
    SAVE IMAGE java-example:latest
```
{% sample lang="Python" %}
```Dockerfile
# Earthfile

# The build starts from a python 3 docker image
FROM python:3
# We change the current working directory
WORKDIR /code

# The above commands are inherited implicitly by all targets below
# (as if they started with FROM +base).

#Declare a target, build
build:
     # Copy the source files from build context to the build environment 
    COPY src src
    # Save the python source in an artifact dir called src (it can be later
    SAVE ARTIFACT src /src

#Declare a target, docker
docker:
    #Define the recipe of the target docker as follows:

    # Copy the artifact /src from the target +build into the current directory with the build container
    COPY +build/src src

    #Set the entrypoint for the resulting docker image
    ENTRYPOINT ["python3", "./src/hello.py"]
    # Save the current state as a docker image, which will have the docker tag
    # python-example:latest. This image is only made available to the host's docker
    # if the entire build succeeds.
    SAVE IMAGE python-example:latest

```
{% endmethod %}

## Continue tutorial

👉 [Part 3: Adding dependencies in the mix](./part-3-adding-dependencies-in-the-mix.md)
