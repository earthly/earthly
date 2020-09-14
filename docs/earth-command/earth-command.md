# Earth command reference

## earth

#### Synopsis

* Target form
  ```
  earth [--build-arg <key>[=<value>]] [--secret|-s <secret-id>[=<value>]]
        [--push] [--no-output] [--no-cache] [--allow-privileged|-P]
        [--ssh-auth-sock <path-to-sock>]
        [--buildkit-host <bk-host>]
        [--interactive|-i]
        <target-ref>
  ```
* Artifact form
  ```
  earth [--build-arg <key>[=<value>]] [--secret|-s <secret-id>[=<value>]]
        [--push] [--no-cache] [--allow-privileged|-P]
        [--ssh-auth-sock <path-to-sock>]
        [--buildkit-host <bk-host>]
        [--interactive|-i]
        --artifact|-a <artifact-ref> [<dest-path>]
  ```
* Image form
  ```
  earth [--build-arg <key>[=<value>]] [--secret|-s <secret-id>[=<value>]]
        [--push] [--no-cache] [--allow-privileged|-P]
        [--ssh-auth-sock <path-to-sock>]
        [--buildkit-host <bk-host>]
        [--interactive|-i]
        --image <target-ref>
  ```

#### Description

The command executes a build referenced by `<target-ref>` (*target form* and *image form*) or `<artifact-ref>` (*artifact form*). In the *target form*, the referenced target and its dependencies are built. In the *artifact form*, the referenced artifact and its dependencies are built, but only the specified artifact is output. The output path of the artifact can be optionally overriden by `<dest-path>`. In the *image form*, the image produced by the referenced target and its dependencies are built, but only the specified image is output.

If a buildkit daemon has not already been started, and the option `--buildkit-host` is not specified, this command also starts up a container named `earthly-buildkitd` to act as a build daemon.

The execution has two phases:

* The build
* The output

During the build phase, the referenced target and all its direct or indirect dependencies are executed. During the output phase, all applicable artifacts with an `AS LOCAL` specification are written to the specified output location, and all applicable docker images are loaded onto the host's docker daemon. If the `--push` option is specified, the output phase additionally pushes any applicable docker images to remote registries and also all `RUN --push` commands are executed.

Remote targets only output images and no artifacts, by default.

If the build phase does not succeed, not output is produced and no push instruction is executed. In this case, the command exits with a non-zero exit code.

The output of the two phases are separated by a `=== SUCCESS ===` marker.

#### Target and Artifact Reference

The `<target-ref>` can reference both local and remote targets.

##### Local Reference

`+<target-name>` will reference a target in the local earthfile in the current directory.

`<local-path>+<target-name>` will reference a local earthfile in a different directory as
specified by `<local-path>`, which must start with `./`, `../`, or `/`.

##### Remote Reference

`<gitvendor>/<namespace>/<project>/path/in/project[:some-tag]+<target-name>` will access a remote git repository.

##### Artifact Reference

The `<artifact-ref>` can reference artifacts built by targets. `<target-ref>/<artifact-path>` will reference a build target's artifact.

##### Examples

See the [Target, artifact, and image referencing guide](../guides/target-ref) for more details and examples.

#### Options

##### `--build-arg <key>[=<value>]`

Also available as an env var setting: `EARTHLY_BUILD_ARGS"<key>=<value>,<key>=<value>,..."`.

Overides the value of the build arg `<key>`. If `<value>` is not specified, then the value becomes the value of the environment variable with the same name as `<key>`. For more information see the [`ARG` Earthfile command](../earthfile/earthfile.md#arg).

##### `--secret|-s <secret-id>[=<value>]`

Also available as an env var setting: `EARTHLY_SECRETS="<secret-id>=<value>,<secret-id>=<value>,..."`.

Passes a secret with ID `<secret-id>` to the build environments. If `<value>` is not specified, then the value becomes the value of the environment variable with the same name as `<secret-id>`.

The secret can be referenced within Earthfile recipes as `RUN --secret <arbitrary-env-var-name>=+secrets/<secret-id>`. For more information see the [`RUN --secret` Earthfile command](../earthfile/earthfile.md#run).

##### `--push`

Also available as an env var setting: `EARTHLY_PUSH=true`.

Instructs Earthly to push any docker images declared with the `--push` flag to remote docker registries and to run any `RUN --push` commands. For more information see the [`SAVE IMAGE` Earthfile command](../earthfile/earthfile.md#save-image) and the [`RUN --push` Earthfile command](../earthfile/earthfile.md#run).

Pushing only happens during the output phase, and only if the build has succeeded.

##### `--no-output`

Also available as an env var setting: `EARTHLY_NO_OUTPUT=true`.

Instructs Earthly not to output any images or artifacts. This option cannot be used with the *artifact form* or the *image form*.

##### `--no-cache`

Also available as an env var setting: `EARTHLY_NO_CACHE=true`.

Instructs Earthly to ignore any cache when building. It does, however, continue to store new cache formed as part of the build (to be possibly used on future invocations).

##### `--allow-privileged|-P`

Also available as an env var setting: `EARTHLY_ALLOW_PRIVILEGED=true`.

Permits the build to use the --privileged flag in RUN commands. For more information see the [`RUN --privileged` command](../earthfile/earthfile.md#run).

##### `--ssh-auth-sock <path-to-sock>`

Also available as an env var setting: `EARTHLY_SSH_AUTH_SOCK=<path-to-sock>`.

Sets the path to the SSH agent sock, which can be used for SSH authentication. SSH authentication is used by Earthly in order to perform git clone's underneath.

On Linux systems, this setting defaults to the value of the env var $SSH_AUTH_SOCK. On most systems, the env var `SSH_AUTH_SOCK` env var is already set if an SSH agent is running.

On Mac systems, this setting defaults to `/run/host-services/ssh-auth.sock` to match recommendation in [the official Docker documentation](https://docs.docker.com/docker-for-mac/osxfs/#ssh-agent-forwarding).

For more information see the [Authentication page](../guides/auth.md).

##### `--git-username <git-user>` (deprecated)

Also available as an env var setting: `GIT_USERNAME=<git-user>`.

This option is now deprecated. Please use the [configuration file](../earth-config/earth-config.md) instead.

##### `--git-password <git-pass>` (deprecated)

Also available as an env var setting: `GIT_PASSWORD=<git-pass>`.

This option is now deprecated. Please use the [configuration file](../earth-config/earth-config.md) instead.

##### `--git-url-instead-of <git-instead-of>` (deprecated)

Also available as an env var setting: `GIT_URL_INSTEAD_OF=<git-instead-of>`.

This option is now deprecated. Please use the [configuration file](../earth-config/earth-config.md) instead.

##### `--interactive|-i` (**experimental**)

Also available as an env var setting: `EARTHLY_INTERACTIVE=true`.

Enable interactive debugging mode. By default when a RUN command fails, earth will display the error and exit. If the interactive mode is enabled and an error occurs, an
interactive shell is presented which can be used for investigating the error interactively. Due to technical limitations, only a single interactive shell can be used on the system at any given time. This feature is experimental and may change over time.


## earth prune

#### Synopsis

* Standard form
  ```
  earth [options] prune [--all|-a]
  ```
* Reset form
  ```
  earth [options] prune --reset
  ```

#### Description

The command `earth prune` eliminates Earthly cache. In the *standard form* it issues a prune command to the buildkit daemon. In the *reset form* it restarts the buildkit daemon, instructing it to completely delete the cache directory on startup, thus forcing it to start from scratch.

#### Options

##### `--all|-a`

Instructs earth to issue a "prune all" command to the buildkit daemon.

##### `--reset`

Restarts the buildkit daemon and completely resets the cache directory.

## earth bootstrap

#### Synopsis

* ```
  earth bootstrap
  ```

#### Description

Installs bash and zsh shell completion for earth.


## earth --help

#### Synopsis

* ```
  earth --help
  ```
* ```
  earth <command> --help
  ```

#### Description

Prints help information about earth.

## earth --version

#### Synopsis

* ```
  earth --version
  ```

#### Description

Prints version information about earth.
