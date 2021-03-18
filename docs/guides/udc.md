# User-defined commands

{% hint style='danger' %}
##### Important

This feature is currently in **Experimental** stage

* The feature may break, be changed drastically with no warning, or be removed altogether in future versions of Earthly.
* Check the [GitHub tracking issue](https://github.com/earthly/earthly/issues/581) for any known problems.
* Give us feedback on [Slack](https://earthly.dev/slack) in the `#udc` channel.
{% endhint %}

User-defined commands (UDCs) are templates (much like functions in regular programming languages), which can be used to define a series of steps to be executed in sequence. In other words, it is a way to import common build steps which can be reused in multiple contexts.

Unlike targets, UDCs inherit the (1) build context and (2) the build environment from the caller. Meaning that

1. Any local `COPY` opperation will use the directory where the calling Earthfile exists, as the source.
2. Any files, directories and dependencies created by a previous step of the caller are available to the UDC to operate on; and any file changes resulting from executing the UDC commands are passed back to the caller as part of the build environment.

Thus, when importing and reusing UDCs across a complex build, it is much like reusing libraries in a regular programming language.

## Usage

UDCs are defined very much like regular targets, with a couple of exceptions: the name is in all-uppercase, snake-case and the recipe must start with `COMMAND`. For example:

```Dockerfile
MY_COPY:
    COMMAND
    ARG src
    ARG dest=./
    ARG recursive=false
    RUN cp $(if $recursive =  "true"; then printf -- -r; fi) "$src" "$dest"
```

This UDC can be invoked from a target via `DO`

```Earthfile
build:
    FROM alpine:3.13
    WORKDIR /udc-example
    RUN echo "hello" >./foo
    DO +MY_COPY --src=./foo --dest=./bar
    RUN cat ./bar # prints "hello"
```

A few things to note about this example:

* The definition of `MY_COPY` does not contain a `FROM` so the build environment it operates in is the build environment of the caller.
* This means that `+MY_COPY` has access to the file `./foo`.
* Although the copy file operation is performed within `+MY_COPY`, its effects are seen in the environment of the caller - so the resulting `./bar` is available to the caller.

## Scope

UDCs create their own `ARG` scope, which is distinct from the caller. Any `ARG` that needs to be passed from the caller needs to be passed explicitly via `DO +COMMAND --<build-arg-key>=<build-arg-value>`, as in the following example.

```Dockerfile
build:
    ARG var=value-in-build
    # prints "something-else"
    DO +PRINT_VAR
    # prints "value-in-build"
    DO +PRINT_VAR --var=$var

PRINT_VAR:
    COMMAND
    ARG var=something-else
    RUN echo "$var"
```

Global imports and global args are inherited from the `base` target of the same Earthfile where the command is defined in (this may be distinct from the `base` target of the caller).

```Dockerfile
ARG a_global_var=value-in-global

build:
    # prints "value-in-global"
    DO +PRINT_VAR

PRINT_VAR:
    COMMAND
    RUN echo "$a_global_var"
```
