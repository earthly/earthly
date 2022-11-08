# Build Arguments and Secrets

## Introduction

One of the core features of Earthly is support for build arguments. Build arguments
can be used to dynamically set environment variables inside the context of [RUN commands](../earthfile/earthfile.md#run).

Build arguments can be passed between targets or from the command line. They encourage
writing generic Earthfiles and ultimately promote greater code-reuse.

Additionally, Earthly defines secrets which are similar to build arguments, but are exposed as environment
variables when explicitly allowed.

## A Quick Example

Arguments are declared with the [ARG](../earthfile/earthfile.md#arg) keyword.

Let's consider a "hello world" example that allows us to change who is being greeted (e.g. hello banana, hello eggplant etc).
We will create a hello target that accepts the `name` argument:

```Dockerfile
VERSION 0.6
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
         +foo | [██████████] 100% resolve docker.io/library/alpine:latest@sha256:21a3deaa0d32a8057914f36584b5288d2e5ecc984380bc0118285c70fa8c9300
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
        +base | [██████████] 100% resolve docker.io/library/alpine:latest@sha256:21a3deaa0d32a8057914f36584b5288d2e5ecc984380bc0118285c70fa8c9300
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

4. From a `.env` file

   It is also possible to create an `.env` file to contain the build arguments to pass
   to earthly. First create an `.env` file with:
   
   ```
   name=eggplant
   ```
   
   Then simply run earthly:
   
   ```bash
   earthly +hello
   ```

## Passing Argument values to targets

Build arguments can also be set when calling build targets. If multiple build arguments values are defined for the same argument name,
Earthly will build the target for each value; this makes it easy to configure a "build matrix" within Earthly.

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

## Passing secrets to RUN commands

Secrets are similar to build arguments; however, they are *not* defined in targets, but instead are explicitly defined for each `RUN` command that is permitted to access them.

Here's an example Earthfile that accesses a secret stored under `+secrets/passwd` and exposes it under the environment variable `mypassword`:

```dockerfile
FROM alpine:latest
hush:
    RUN --secret mypassword=+secrets/passwd echo "my password is $mypassword"
```

If the environment variable name is identical to the secret ID. For example to accesses a secret stored under `+secrets/passwd` and exposes it under the environment variable `passwd`  you can use the shorthand :

```dockerfile
FROM alpine:latest
hush:
    RUN --secret passwd echo "my password is $passwd"
```

{% hint style='info' %}
It's also possible to temporarily mount a secret as a file:

```dockerfile
RUN --mount type=secret,target=/root/mypassword,id=+secrets/passwd echo "my password is $(cat /root/mypassword)"
```

The file will not be saved to the image snapshot.
{% endhint %}

## Setting secret values

The value for `+secrets/passwd` in examples above must then be supplied when earthly is invoked.

This is possible in a few ways:


1. Directly, on the command line:

   ```bash
   earthly --secret passwd=itsasecret +hush
   ```

2. Via an environment variable:

   ```bash
   export passwd=itsasecret
   earthly --secret passwd +hush
   ```

   If the value of the secret is omitted on the command line Earthly will lookup the environment variable with that name.

3. Via the environment variable `EARTHLY_SECRETS`

   ```bash
   export EARTHLY_SECRETS="passwd=itsasecret"
   earthly +hush
   ```

   Multiple secrets can be specified by separating them with a comma.

4. Via the `.env` file.

   Create a `.env` file in the same directory where you plan to run `earthly` from. Its contents should be:
   
   ```
   passwd=itsasecret
   ```
   
   Then simply run earthly:
   
   ```bash
   earthly +hello
   ```

5. Via cloud-based secrets. This option helps share secrets within a wider team. To read more about this see the [cloud-based secrets guide](../cloud/cloud-secrets.md).

Regardless of the approach chosen from above, once earthly is invoked, in our example, it will output:

```
+hush | --> RUN echo "my password is $mypassword"
+hush | my password is itsasecret
```

{% hint style='info' %}
### How Arguments and Secrets affect caching

Commands in earthly must be re-evaluated when the command itself changes (e.g. `echo "hello $name"` is changed to `echo "greetings $name"`), or when
one of its inputs has changed (e.g. `--name=world` is changed to `--name=banana`). Earthly creates a hash based on both the contents
of the command and the contents of all defined arguments of the target build context.

However, in the case of secrets, the contents of the secret *is not* included in the hash; therefore, if the contents of a secret changes, Earthly is unable to
detect such a change, and thus the command will not be re-evaluated.
{% endhint %}

## Storage of local secrets

Earthly stores the contents of command-line-supplied secrets in memory on the localhost. When a `RUN` command that requires a secret is evaluated by BuildKit, the BuildKit
daemon will request the secret from the earthly command-line process and will temporarily mount the secret inside the runc container that is evaluating the `RUN` command.
Once the command finishes the secret is unmounted. It will not persist as an environment variable within the saved container snapshot. Secrets will be kept in-memory
until the earthly command exits.

Earthly also supports cloud-based shared secrets which can be stored in the cloud. Secrets are never stored in the cloud unless a user creates an earthly account and
explicitly calls the `earthly secrets set ...` command to transmit the secret to the earthly cloud-based secrets server.
For more information about cloud-based secrets, check out our [cloud-based secrets management guide](../cloud/cloud-secrets.md).
