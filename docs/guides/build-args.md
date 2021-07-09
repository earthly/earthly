# Build Arguments and Secrets

## Introduction

One of the core features of earthly is support for build arguments. Build arguments
can be used to dynamically set environment variables inside the context of [RUN commands](../earthfile/earthfile.md#run).

Build arguments can be passed between targets or from the command line. They encourage
writing generic Earthfiles and ultimately promote greater code-reuse.

Additionally, earthly defines secrets which are similar to build arguments, but are exposed as environment
variables when explicitly allowed.

## A Quick Example

Arguments are declared with the [ARG](../earthfile/earthfile.md#arg) keyword.

Let's consider a "hello world" example that allows us to change who is being greeted (e.g. hello banana, hello eggplant etc).
We will create a hello target that accepts the `name` argument:

```Dockerfile
FROM alpine:latest

hello:
    ARG name
    RUN echo "hello $name"
```

Then we will specify a value for the `name` argument on the command line when we invoke `earthly`:

```bash
earthly --build-arg name=world +hello
```

This will output

```
    buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
alpine:latest | --> Load metadata linux/amd64
         +foo | --> FROM alpine:latest
         +foo | [██████████] resolve docker.io/library/alpine:latest@sha256:69e70a79f2d41ab5d637de98c1e0b055206ba40a8145e7bddb55ccc04e13cf8f ... 100%
         +foo | name=world
         +foo | --> RUN echo "hello $name"
         +foo | hello world
       output | --> exporting outputs
```

If we re-run `earthly --build-arg name=world +hello`, we will see that the echo command is cached (and won't re-display the hello world text):

```
+foo | *cached* --> RUN echo "hello $name"
```

## Setting Argument Values

Argument values can be set multiple ways:

1. On the command line

   The value can be directly specified on the command line (as shown in the previous example):
   
   ```
   earthly --build-arg name=world +hello
   ```

2. From an environment variable

   If no value is given for name, then earthly will look for the value in the corresponding
   environment variable on the localhost:
   
   ```bash
   export name="banana"
   earthly --build-arg name +hello
   ```

3. From a `.env` file

   It is also possible to create an `.env` file to contain the build arguments to pass
   to earthly. First create an `.env` file with:
   
   ```
   name eggplant
   ```
   
   Then simply run earthly:
   
   ```bash
   earthly +hello
   ```

## Passing Argument values to targets

Build arguments can also be set when calling build targets. If multiple build arguments values are defined for the same argument name,
earthly will build the target for each value; this makes it easy to configure a "build matrix" within Earthly.

For example, we can create a new `greetings` target which calls `+hello` multiple times:

```dockerfile
greetings:
    BUILD \
        --build-arg name=world \
        --build-arg name=banana \
        --build-arg name=eggplant \
        +hello
```

Then when we call `earthly +greetings`, earthly will call `+hello` three times:

```
     buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
 alpine:latest | --> Load metadata linux/amd64
         +base | --> FROM alpine:latest
         +base | [██████████] resolve docker.io/library/alpine:latest@sha256:69e70a79f2d41ab5d637de98c1e0b055206ba40a8145e7bddb55ccc04e13cf8f ... 100%
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

In addition to the `BUILD` command, the `--build-arg` flag can also be used with `FROM`, `COPY` and a number of other commands.

## Passing secrets to RUN commands

Secrets are similar to build arguments; however, they are *not* defined in targets, but instead are explicitly defined for each `RUN` command that is permitted to access them.

Here's an example Earthfile that accesses a secret stored under `+secrets/passwd` and exposes it under the environment variable `mypassword`:

```dockerfile
FROM alpine:latest
hush:
    RUN --secret mypassword=+secrets/passwd echo "my password is $mypassword"
```

{% hint style='info' %}
It's also possible to temporarily mount a secret as a file:

```dockerfile
RUN --mount type=secret,target=/root/mypassword,id=+secrets/passwd echo "my password is $(cat /root/mypassword)"
```

The file will not be saved to the image snapshot.
{% endhint %}

The value for `+secrets/passwd` must then be supplied when earthly is invoked. This can be either done directly via:

```bash
earthly --secret passwd=itsasecret +hush
```

or if the value is omitted, then earthly will attempt to lookup the value from an environment variable on the localhost:

```bash
passwd=itsasecret \
earthly --secret passwd +hush
```

Alternatively, earthly offers [cloud-based secrets](cloud-secrets.md) if you need to share secrets between colleagues.

Once earthly is invoked, it will output:

```
+hush | --> RUN echo "my password is $mypassword"
+hush | my password is itsasecret
```

{% hint style='info' %}
### How Arguments affect caching

Commands in earthly must be re-evaluated when the command itself changes (e.g. `echo "hello $name"` is changed to `echo "greetings $name"`), or when
one of it's inputs has changed (e.g. `--build-arg name=world` is changed to `--build-arg name=banana`). Earthly creates a hash based on both the contents
of the command and the contents of all defined arguments of the target build context.

However, in the case of secrets, the contents of the secret *is not* included in the hash; therefore, if the contents of a secret changes, earthly is unable to
detect such a change, and thus the command will not be re-evaluated.
{% endhint %}

## Storage of local secrets

Earthly stores the contents of command-line-supplied secrets in memory on the localhost. When a `RUN` command that requires a secret is evaluated by BuildKit, the BuildKit
daemon will request the secret from the earthly command-line process and will temporarily mount the secret inside the runc container that is evaluating the `RUN` command.
Once the command finishes the secret is unmounted. It will not persist as an environment variable within the saved container snapshot. Secrets will persist in-memory
until the earthly command exits.

Earthly also supports cloud-based shared secrets which can be stored in the cloud. Secrets are never stored in the cloud unless a user creates an earthly account and
explicitly calls the `earthly secrets set ...` command to transmit the secret to the earthly cloud-based secrets server.
For more information about cloud-based secrets, check out our [cloud-based secrets management guide](cloud-secrets.md).
