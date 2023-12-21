# Build Arguments and Variables

## Introduction

One of the core features of Earthly is support for build arguments. Build arguments are declared with `ARG` and
can be used to dynamically set environment variables inside the context of [RUN commands](../earthfile/earthfile.md#run).

Build arguments can be passed between targets or from the command line. They encourage
writing generic Earthfiles and ultimately promote greater code-reuse.

Another closely related primitive that Earthly offers is the variable (declared with `LET`). Variables are similar to build arguments, except that they cannot be used as parameters.

## A Quick Example

Arguments are declared either with the [ARG](../earthfile/earthfile.md#arg) keyword.

Let's consider a "hello world" example that allows us to change who is being greeted (e.g. hello banana, hello eggplant etc).
We will create a hello target that accepts the `name` argument:

```Dockerfile
VERSION 0.8
FROM alpine:latest

hello:
    ARG name
    RUN echo "hello $name"
```

Then we will specify a value for the `name` argument on the command line when we invoke `earthly`:

```bash
earthly +hello --name=world
```

This will output

```
    buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
alpine:latest | --> Load metadata linux/arm64
         +foo | --> FROM alpine:latest
         +foo | 100% resolve docker.io/library/alpine:latest@sha256:21a3deaa0d32a8057914f36584b5288d2e5ecc984380bc0118285c70fa8c9300
         +foo | name=world
         +foo | --> RUN echo "hello $name"
         +foo | hello world
       output | --> exporting outputs
```

If we re-run `earthly +hello --name=world`, we will see that the echo command is cached (and won't re-display the hello world text):

```
+foo | *cached* --> RUN echo "hello $name"
```

## Default values

Arguments may also have default values, which may be either constant or dynamic. For example, the following target will greet the name identified by the arg `name` (which has a default value of John), with the current time:

```Dockerfile
hello:
   ARG time=$(date +%H:%M)
   ARG name=John
   RUN echo "hello $name, it is $time"
```

```
alpine:latest | --> Load metadata linux/arm64
        +base | --> FROM alpine:latest
        +base | 100% resolve docker.io/library/alpine:latest@sha256:21a3deaa0d32a8057914f36584b5288d2e5ecc984380bc0118285c70fa8c9300
       +hello | --> ARG time = RUN $(date +%H:%M)
       +hello | --> RUN echo "hello $name, it is $time"
       +hello | hello John, it is 23:21
       output | --> exporting outputs
```

If an arg has no default value, then the default value is the empty string.

## Overriding Argument Values

Argument values can be set multiple ways:

1. On the command line

   The value can be directly specified on the command line (as shown in the previous example):
   
   ```
   earthly +hello --HELLO=world --FOO=bar
   ```

2. From environment variables

   Similar to above, except that the value is an environment variable:
   
   ```bash
   export HELLO="world"
   export FOO="bar"
   earthly +hello --HELLO="$HELLO" --FOO="$FOO"
   ```

3. Via the `EARTHLY_BUILD_ARGS` environment variable

    The value can also be set via the `EARTHLY_BUILD_ARGS` environment variable.
    
    ```bash
    export EARTHLY_BUILD_ARGS="HELLO=world,FOO=bar"
    earthly +hello
    ```

    This may be useful if you have a set of build args that you'd like to always use and would prefer not to have to specify them on the command line every time. The `EARTHLY_BUILD_ARGS` environment variable may also be stored in your `~/.bashrc` file, or some other shell-specific startup script.

4. From an `.arg` file

   It is also possible to create an `.arg` file to contain the build arguments to pass
   to earthly. First create an `.arg` file with:
   
   ```
   name=eggplant
   ```
   
   Then simply run earthly:
   
   ```bash
   earthly +hello
   ```

## Passing Argument values to targets

Build arguments can also be set when calling build targets.

```Dockerfile
greeting:
   BUILD +hello --name=world

hello:
    ARG name
    RUN echo "hello $name"
```

Arg overrides within the same Earthfile are passed automatically to each other. In the example below, if you are calling `earthly +greeting --name=world`, the `--name=world` override will be passed to `+hello` as well.

```Dockerfile
greeting:
   BUILD +hello

hello:
   ARG name
   RUN echo "hello $name"
```

This behavior does not apply to references to other Earthfiles. In order to pass arguments to other Earthfiles, you must either explicitly pass the argument. For example:

```Dockerfile
ARG name
BUILD +hello --name=$name
```

Or you can use the `--pass-args` flag to pass all arguments to the target:

```Dockerfile
BUILD --pass-args +hello
```

### Matrix builds

If multiple build arguments values are defined for the same argument name, Earthly will build the target for each value; this makes it easy to configure a "build matrix" within Earthly.

For example, we can create a new `greetings` target which calls `+hello` multiple times:

```dockerfile
greetings:
    BUILD +hello \
        --name=world \
        --name=banana \
        --name=eggplant
```

Then when we call `earthly +greetings`, earthly will call `+hello` three times:

```
     buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
 alpine:latest | --> Load metadata linux/amd64
         +base | --> FROM alpine:latest
         +base | resolve docker.io/library/alpine:latest@sha256:69e70a79f2d41ab5d637de98c1e0b055206ba40a8145e7bddb55ccc04e13cf8f ... 100%
        +hello | name=banana
        +hello | --> RUN echo "hello $name"
        +hello | name=eggplant
        +hello | --> RUN echo "hello $name"
        +hello | name=world
        +hello | --> RUN echo "hello $name"
        +hello | hello banana
        +hello | hello eggplant
        +hello | hello world
        output | --> exporting outputs
```

In addition to the `BUILD` command, build args can also be used with `FROM`, `COPY`, `WITH DOCKER --load` and a number of other commands:

```Dockerfile
BUILD +hello --name=world
COPY (+hello/file.txt --name=world) ./
FROM +hello --name=world
WITH DOCKER --load=(+hello --name=world)
  ...
END
```

Another way to pass build args is by specifying a dynamic value, delimited by `$(...)`. For example, in the following, the value of the arg `name` will be set as the output of the shell command `echo world` (which, of course is simply `world`):

```Dockerfile
BUILD +hello --name=$(echo world)
```

## Variables

Variables are similar to build arguments, except that they cannot be used as parameters. You can think of variables as "private" build arguments (or local variables). To declare a variable, you can use the `LET` command.

Variables can also be mutated via the `SET` command. For example:

```Dockerfile
hello:
   LET name = "world"
   RUN echo "hello $name"
   SET name = "banana"
   RUN echo "hello $name"
```

This can be useful when you would like to decide on the value of a variable based on an `IF` condition, or if you would like to construct the value of the variable via a `FOR` loop.

For more information on `LET` see the [`LET` Earthfile reference](../earthfile/earthfile.md#let).
