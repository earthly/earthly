
# Basics

You're here to learn how to use Earthly. Welcome! Earthly is a build system based upon modern development technologies available today, like Docker.

First, a few concepts:

* Earthly - the build system as a whole
* `earth` - the CLI tool used to interact with Earthly
* earthfile - a file named `build.earth`, which contains a series of targets and their respective recipes
* buildkitd - a [daemon built by Docker team](https://github.com/moby/buildkit) and used by Earthly to execute builds. It executes LLB, the same low-level primitives used when building Dockerfiles. The buildkitd daemon is started automatically in a docker container, by `earth`, when executing builds.
* recipe - a specific series of build steps
* target - the label used to identify a recipe
* build context - the main directory made available to the build for copying files from
* artifact - a file resulting from executing a target (not all targets have artifacts)
* image - a docker image resulting from executing a target (not all targets have images)

Now that we got that out of the way, we can now be more specific:

Earthly is a build system where all recipes are executed in docker containers, with layer caching principles similar to Dockerfiles, and where the only thing shared between recipes are declared artifacts and images.

A key difference from a Dockerfile build is that Earthly can be used to build not just images, but also artifacts - files that can be written back onto the host filesystem.

## Target referencing

Targets have a particular referencing convention which helps Earthly to identify which recipe to execute.

### Local, internal

The simplest form, is where a target is referenced from the same directory:

`+<target-name>`

For example,

`+build`

In this form, Earthly will look for the target within the same dir (or within the same earthfile). We call this type of referencing local, internal. Local, because it comes from the same system, and internal, because it is within the same earthfile.

### Local, external

Another form, is where a target is referenced from a different directory. In this form, the path to that directory is specified before `+`. It must always start with either `./`, `../` or `/`, on any operating system (including Windows). Example:

`./path/to/another/dir+<target-name>`

For example:

`./js+build`

It is recommended that relative paths are used, for portability reasons: the working directory checked out by different users will be different, making absolute paths infeasible in most cases.


### Remote

Finally, the last form of target referencing is the remote form. In this form, the recipe and the build context are imported from a remote location. It has the following form:

`<vendor>/<namespace>/<project>/path/in/project:some-tag+<target-name>`

For example:

`github.com/vladaionescu/earthly/buildkitd+buildkitd`

or

`github.com/vladaionescu/earthly:v0.1+all`

## Earthfile basics

Earthfiles are always named `build.earth`, regardless of their location in the codebase. Earthfiles have the following rough structure:

```
...
base recipe
...

target:
    ...
    recipe
    ...

target:
    ...
    recipe
    ...
```

For example, here is an earthfile of an example go app:

```earth
FROM golang:1.13-alpine3.11
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

You will notice that the recipes look very much like Dockerfiles. This is an intentional design decision. Existing Dockerfiles can be ported to earthfiles by copy-pasting them over and then tweaking them slightly. Compared to Dockerfile syntax, some commands are new (like `SAVE ARTIFACT`), others have additional semantics (like `COPY +target/some-artifact`) and other semantics are removed (like `FROM ... AS ...` and `COPY --from`).

All earthfiles start with a base recipe. This is the only recipe which does not have an explicit target name - the name is always implied to be `base`. All other target implicitly inherit from `base`. You can imagine that all recipes start with an implicit `FROM +base`.

In this particular example, we can see two explicit targets: `build` and `docker`. In order to execute the build, we can run, for example:

```bash
earth +docker
```

The output might look like this:

**TODO: Redo output**

```
...
```

Notice how to the left of `|`, within the output, we can see some targets like `+base`, `+build` and `+docker` . Notice how the output is interleaved between `+docker` and `+build`. This is because the system executes independent build steps in parallel. The reason this is possible effortlessly is because only very few things are shared between the builds of the recipes and those things are declared and obvious. The rest is completely isolated.

In addition, notice how even though the base is used as part of both `build` and `docker`, it is only executed once. This is because the system deduplicates execution, where possible.

Furthermore, the fact that the `docker` target depends on the `build` target is visible within the command `COPY +build/go-example .`. Through this command, the system knows that it also needs to build the target `+build`, in order to satisfy the dependency on the artifact.

Finally, notice how the output of the build: the docker image `go-example:latest` and the file `build/go-example` is only written after the build is declared a success. This is due to another isolation principle of Earthly: a build either succeeds completely or it fails altogether.

Going back to the example earthfile definition, here is what each command does:

```earth
# The build starts from a docker image: golang:1.13-alpine3.11
FROM golang:1.13-alpine3.11
# We change the current working directory.
WORKDIR /go-example

# The above commands are inherited implicitly by all targets below
# (as if they started with FROM +base).

# Declare a target, build.
build:
    # Define the recipe of the target build as follows:

    # Copy main.go from the build context to the current dir within the build container, as a layer.
    COPY main.go .
    # Run a go build command as a layer. This uses the previously copied main.go file.
    RUN go build -o build/go-example main.go
    # Save the output of the build command as an artifact. Call this artifact /go-example (it
    # can be later referenced as +build/go-example. In addition, store the artifact as a
    # local file (on the host) named build/go-example. This local file is only written if the
    # entire build succeeds.
    SAVE ARTIFACT build/go-example /go-example AS LOCAL build/go-example

# Declare a target, docker.
docker:
    # Define the recipe of the target docker as follows:

    # Copy the artifact /go-example produced by another target, +build, to the current directory
    # within the build container.
    COPY +build/go-example .
    # Set the entrypoint for the resulting docker image.
    ENTRYPOINT ["/go-example/go-example"]
    # Save the current state as a docker image, which will have the docker tag go-example:latest.
    # This image is only made available to the host's docker if the entire build succeeds.
    SAVE IMAGE go-example:latest
```
