# Functions

{% hint style='hint' %}
#### UDCs have been renamed to Functions

Functions used to be called UDCs (User Defined Commands). Earthly 0.7 uses `COMMAND` instead of `FUNCTION`.
{% endhint %}

Earthly Functions are reusable sets of instructions that can be inserted in targets or other functions. In other words, it is a way to import common build steps which can be reused in multiple contexts.

Unlike targets, functions inherit the (1) build context and (2) the build environment from the caller. Meaning that

1. Any local `COPY` operation will use the directory where the calling Earthfile exists, as the source.
2. Any files, directories and dependencies created by a previous step of the caller are available to the function to operate on; and any file changes resulting from executing the function's commands are passed back to the caller as part of the build environment.

Thus, when importing and reusing functions across a complex build, it is very much like reusing libraries in a regular programming language.

## Usage

Functions are defined similarly to regular targets, with a couple of exceptions: the name is in ALL_UPPERCASE_SNAKE_CASE and the recipe must start with `FUNCTION`. For example:

```Dockerfile
MY_COPY:
    FUNCTION
    ARG src
    ARG dest=./
    ARG recursive=false
    RUN cp $(if $recursive =  "true"; then printf -- -r; fi) "$src" "$dest"
```

This function can be invoked from a target via `DO`

```Earthfile
build:
    FROM alpine:3.18
    WORKDIR /function-example
    RUN echo "hello" >./foo
    DO +MY_COPY --src=./foo --dest=./bar
    RUN cat ./bar # prints "hello"
```

A few things to note about this example:

* The definition of `MY_COPY` does not contain a `FROM` so the build environment it operates in is the build environment of the caller.
* This means that `+MY_COPY` has access to the file `./foo`.
* Although the copy file operation is performed within `+MY_COPY`, its effects are seen in the environment of the caller - so the resulting `./bar` is available to the caller.

## Scope

Functions create their own `ARG` scope, which is distinct from the caller. Any `ARG` that needs to be passed from the caller needs to be passed explicitly via `DO +MY_FUNCTION --<build-arg-key>=<build-arg-value>`, as in the following example.

```Dockerfile
build:
    ARG var=value-in-build
    # prints "something-else"
    DO +PRINT_VAR
    # prints "value-in-build"
    DO +PRINT_VAR --var=$var

PRINT_VAR:
    FUNCTION
    ARG var=something-else
    RUN echo "$var"
```

Global imports and global args are inherited from the `base` target of the same Earthfile where the command is defined in (this may be distinct from the `base` target of the caller).

```Dockerfile
VERSION 0.8

ARG --global a_global_var=value-in-global

build:
    # prints "value-in-global"
    DO +PRINT_VAR

PRINT_VAR:
    FUNCTION
    RUN echo "$a_global_var"
```

## Targets vs Functions

Targets and functions are Earthly's core primitives for organizing build recipes. They encapsulate build logic, and from afar they look pretty similar. However, the use-cases for each are vastly different.

In general, targets are used to produce specific build results, while functions are used as a way to reuse build logic, when certain commands are repeated in multiple places. As a real-world analogy, targets are more like factories, while functions are more like components that are used to put together factories.

Here is a comparison of the two primitives:

| | Targets | Functions |
| --- | --- | --- |
| Represents a collection of Earthly commands | ✅ | ✅ |
| Can reference other targets in its body | ✅ | ✅ |
| Can reference other functions in its body | ✅ | ✅ |
| Build context | The directory where the Earthfile resides | Inherited from the caller |
| Build environment, when no `FROM` is specified | Inherited from the base of its own Earthfile | Inherited from the caller |
| `IMPORT` statements | Inherited from the base of its own Earthfile | Inherited from the base of its own Earthfile |
| `ARG` context | Creates its own scope | Creates its own scope |
| Requires that `ARG`s be passed in explicitly | ✅ | ✅ |
| Global `ARG` context | Inherited from the base of its own Earthfile | Inherited from the base of its own Earthfile |
| Can output artifacts | ✅ | ❌ - can issue `SAVE ARTIFACT`, but it's the caller that emits the artifacts |
| Can output images | ✅ | ❌ - can issue `SAVE IMAGE`, but it's the caller that emits the images |
| Can be called via `earthly` CLI | ✅ | ❌ |
| Can be used in conjunction with an `IMPORT` | ✅ - `FROM some-import+my-target` | ✅ - `DO some-import+MY_FUNCTION` |
| Commands that can reference it | `FROM`, `BUILD`, `COPY`, `WITH DOCKER --load`, `FROM DOCKERFILE` | `DO` |
