# earth command reference

## earth

#### Synopsis

* ```
  earth [--build-arg <key>[=<value>]] [--secret|-s <secret-id>[=<value>]]
        [--push] [--no-output] [--no-cache] [--allow-privileged|-P]
        [--ssh-auth-sock <path-to-sock>]
        [--git-username <git-user>] [--git-password <git-pass>]
        [--git-url-instead-of <git-instead-of>]
        [--buildkit-host <bk-host>] [--buildkit-cache-size-mb <cache-size-mb>]
        [--buildkit-image <bk-image>]
        <target-ref>
  ```
  (target form)
* ```
  earth [--build-arg <key>[=<value>]] [--secret|-s <secret-id>[=<value>]]
        [--push] [--no-cache] [--allow-privileged|-P]
        [--ssh-auth-sock <path-to-sock>]
        [--git-username <git-user>] [--git-password <git-pass>]
        [--git-url-instead-of <git-instead-of>]
        [--buildkit-host <bk-host>] [--buildkit-cache-size-mb <cache-size-mb>]
        [--buildkit-image <bk-image>]
        --artifact|-a <artifact-ref>
  ```
  (artifact form)
* ```
  earth [--build-arg <key>[=<value>]] [--secret|-s <secret-id>[=<value>]]
        [--push] [--no-cache] [--allow-privileged|-P]
        [--ssh-auth-sock <path-to-sock>]
        [--git-username <git-user>] [--git-password <git-pass>]
        [--git-url-instead-of <git-instead-of>]
        [--buildkit-host <bk-host>] [--buildkit-cache-size-mb <cache-size-mb>]
        [--buildkit-image <bk-image>]
        --image|-i <image-ref>
  ```
  (image form)

#### Description

The command executes a build referenced by `<target-ref>` (*target form*), `<artifact-ref>` (*artifact form*) or `<image-ref>` (*image form*). In the *target form*, the referenced target and its dependencies are built. In the *artifact form*, the referenced artifact and its dependencies are built, but only the specified artifact is output. In the *image form*, the referenced image and its dependencies are built, but only the specified image is output.

If a buildkit daemon has not already been started, and the option `--buildkit-host` is not specified, this command also starts up a container named `earthly-buildkitd` to act as a build daemon.

The build has two stages:

* The build itself
* The output

During the build phase, the referenced target and all its direct or indirect dependencies are executed. During the output phase, all applicable artifacts with an `AS LOCAL` specification are written to the specified output location, and all applicable docker images are loaded onto the host's docker daemon. If the `--push` option is specified, the output phase additionally pushes any applicable docker images to remote registries and also all `RUN --push` commands are executed.

If the build phase does not succeed, not output is produced and no push instruction is executed. In this case, the command exits with a non-zero exit code.

The output of the two phases are separated by a `=== SUCCESS ===` marker.

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

Pushing only happens during the output stage, and only if the build has succeeded.

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

Also available as an env var setting: `SSH_AUTH_SOCK=<path-to-sock>`.

Sets the path to the SSH agent sock, which can be used for SSH authentication. SSH authentication is used by Earthly in order to perform git clone's underneath. On most systems, the env var `SSH_AUTH_SOCK` env var is already set if a SSH agent is running. For more information see the [Authentication page](../guides/auth.md).

##### `--git-username <git-user>`

Also available as an env var setting: `GIT_USERNAME=<git-user>`.

Sets the git username to use for non-SSH git authentication. For more information see the [Authentication page](../guides/auth.md).

To prevent the need to specify this option on every `earth` invocation, it is recommended to specify the env var form in a file like `.bashrc` or `.profile`.

##### `--git-password <git-pass>`

Also available as an env var setting: `GIT_PASSWORD=<git-pass>`.

Sets the git password to use for non-SSH git authentication. For more information see the [Authentication page](../guides/auth.md).

{% hint style='danger' %}
##### Important

For security reasons, it is strongly recommended to use the env var form of this setting and not the flag form.
{% endhint %}

##### `--git-url-instead-of <git-instead-of>`

Also available as an env var setting: `GIT_URL_INSTEAD_OF=<git-instead-of>`.

Rewrites git URLs of a certain pattern. Similar to [`git-config url.<base>.insteadOf`](https://git-scm.com/docs/git-config#Documentation/git-config.txt-urlltbasegtinsteadOf). Format: `<base>=<instead-of>`.

This setting allows rewriting all git URLs of the form `https://...` into `git@github.com:...`, or vice-versa.

For example:

* `--git-url-instead-of='git@github.com:=https://github.com/'` forces use of SSH-based URLs for GitHub (default)
* `--git-url-instead-of='https://github.com/=git@github.com:'` forces use of HTTPS-based URLs for GitHub

Currently, only one such pattern is allowed to be specified.

To prevent the need to specify this option on every `earth` invocation, it is recommended to specify the env var form in a file like `.bashrc` or `.profile`.

##### `--buildkit-host <bk-host>`

Also available as an env var setting: `EARTHLY_BUILDKIT_HOST=<bk-host>`.

Instructs `earth` to use an alternate buildkit host. When this option is specified, `earth` does not manage (starts/restarts as necessary) the buildkit daemon.

##### `--buildkit-cache-size-mb <cache-size-mb>`

Also available as an env var setting: `EARTHLY_BUILDKIT_CACHE_SIZE_MB=<cache-size-mb>`.

The total size of the buildkit cache, in MB. The buildkit daemon will allocate disk space for this size. Size less than `1000` (1GB) is not recommended. The default size if this option is not set is `10000` (10GB).

This setting is only used when the buildkit daemon is started (or restarted). In order to apply the setting immediately, issue the command

```bash
earth --buildkit-cache-size-mb <cache-size-mb> prune --reset
```

##### `--buildkit-image <bk-image>`

Also available as an env var setting: `EARTHLY_BUILDKITD_IMAGE=<bk-image>`.

Instructs earth to use an alternate image for buildkitd. The default image used is `earthly/buildkitd:<earth-version>`.

## earth prune

#### Synopsis

* ```
  earth [options] [--all|-a] prune (standard form)
  ```
* ```
  earth [options] --reset prune (reset form)
  ```

#### Description

The command `earth prune` eliminates Earthly cache. In the *standard form* it issues a prune command to the buildkit daemon. In the *reset form* it restarts the buildkit daemon, instructing it to completely delete the cache directory on startup, thus forcing it to start from scratch.

#### Options

##### `--all|-a`

Instructs earth to issue a "prune all" command to the buildkit daemon.

##### `--reset`

Restarts the buildkit daemon and completely resets the cache directory.

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
