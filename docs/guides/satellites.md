# Earthly Satellites

This feature is part of the Earthly Satellites paid plan.

{% hint style='danger' %}
##### Important

This feature is currently in **Beta** stage

* The feature may break or change significantly in future versions of Earthly.
* Give us feedback on [Slack](https://earthly.dev/slack) or via GitHub issues
{% endhint %}

Earthly Satellites are remote Buildkit instances managed by the Earthly team. They allow you to perform builds remotely, retaining cache between runs.

TODO: Diagram!

When using Earthly Satellites, even though the build executes remotely, the following pieces of functionality are still available:

* Build logs are streamed in real-time
* Outputs (images and artifacts) resulting from the build are transferred back to your local machine
* Commands under `LOCALLY` execute on your local machine
* Secrets available locally, including Docker/Podman credentials are passed to the satellite whenever needed by the build

## Benefits

Typical use cases for Earthly Satellites include:

* Speeding up CI builds in sandboxed CI environments such as GitHub Actions, GitLab, CircleCI, and others. Most CI build times are improved by 20-80% via Satellites.
* Executing builds on AMD64/Intel architecture natively when working from an Apple Silicon machine (Apple M1/M2).
* Sharing compute and cache with coworkers.
* Benefiting from high-bandwidth internet access from the satellite, thus allowing for fast downloads of dependencies and fast pushes for deployments. This is particularly useful if operating from a location with slow internet.
* Using Earthly in environments where privileged access or docker-in-docker are not supported.

## Security

As builds often handle sensitive pieces of data, Satellites are designed with security in mind. Here are some of Earthly's security considerations:

* Network communication and data at rest is secured using industry state of the art practices.
* Satellite instances run in sandboxed, isolated environments are are only accessible by users you invite onto the platform.
* The cache is not shared between satellites.
* Secrets used as part of the build are only kept in-memory temporarily, unless they are part of the [Earthly Cloud Secrets storage](../cloud-secrets.md), in which case they are encrypted at rest.
* In addition, Earthly is pursuing SOC 2 compliance.
* To read more about Earthly's security practices please see the [Security page](https://earthly.dev/security).

## Getting started

### 1. Request access

Satellites is currently a private beta feature. Please [contact us](mailto:support+satellite@earthly.dev) to join the beta.

### 2. Register an account

To get started with Satellites, you'll need to register an Earthly account, if you haven't already. You can do so by visiting [Earthly CI](https://ci.earthly.dev), or by using the CLI as described below.

```bash
earthly account register --email <email>
```

An email will be sent to you containing a verification token. Next run:

```bash
earthly account register --email <email> --token <token>
```

This command will prompt you to set a password, and to optionally register a public-key for password-less authentication.

### 3. Create or join an Earthly org

If you haven't already, create an Earthly org by running:

```bash
earthly org create <org-name>
```

To invite another user to join your org, run:

```bash
earthly org invite /<org-name>/ <email>
```

Note the slashes around the org name. Also, please note that the user must have an account on Earthly before they can be invited. (This is a temporary limitation which will be addressed in the future.)

### 4. Launch a new satellite

To launch a new satellite, run:

```bash
earthly sat launch <satellite-name>
```

The Satellite name can be any arbitrary string.

If you are part of multiple Earthly organizations, you may have to specify the org name under which you would like to launch the satellite:

```bash
earthly sat --org <org-name> launch <name>
```

Once the satellite is created it will be automatically selected for use as part of your builds. The selection takes place by Earthly adding some information in your Earthly config file (usually located under `~/.earthly/config.yml`).

### 5. Run a build

To execute a build using the newly created satellite, simply run Earthly like you always have. For example:

```bash
earthly +my-target
```

Because the satellite has been automatically selected in the step above, the build will be executed on it.

To go back to using your local machine for builds, you may "unselect" the satellite by running:

```bash
earthly sat unselect
```

You can always go back to using the satellite by running:

```bash
earthly sat select <satellite-name>
```

Or, you can use a satellite only for a specific build, even if it is not selected:

```bash
earthly --sat <satellite-name> +my-target
```

Conversely, if a satellite is currently selected, but you want to execute a build on your local machine, you can use the `--no-sat` flag:

```bash
earthly --no-sat +my-target
```

## Working with Satellites

To list the satellites available in your organization, run:

```bash
earthly sat ls
```

To check the status of a satellite, you can run:

```bash
earthly sat inspect <satellite-name>
```

To clear a satellite's cache, run the following command after selecting the satellite:

```bash
earthly purge -a
```

## Satellite specs

Satellites are currently only available in one size, and it has the following specs:

* 4 CPUs
* 16 GB of RAM
* 90 GB of cache storage
* 5 Gib internet bandwidth

## Using Satellites in CI

A key benefit of using satellites in a CI environment is that the cache is shared between runs. This results in significant speedups in CIs that would otherwise have to start from scratch each time.

{% hint style='danger' %}
##### Note

If a satellite is shared between multiple CI pipelines, it is possible that it becomes overloaded by too many parallel builds. For best performance, you can create a dedicated satellite for each CI pipeline.
{% endhint %}

To get started with using Earthly Satellites in CI, you can create a login token for access.

First, run

```bash
earthly account create-token <token-name>
```

to create your login token.

Copy and paste the value into an environment variable called `EARTHLY_TOKEN` in your CI environment.

Then as part of your CI script, just run

```bash
earthly satellite select <name>
```

before running your Earthly targets.

## Known limitations

* Satellites currently require a manual re-launch in order to get updated to the latest version available.
  ```bash
  earthly satellite rm <name>
  earthly satellite launch <name>
  ```
* The output phase (the phase in which a satellite outputs build results back to the local machine) is slower than it could be. To work around this issue, you can make use of the `--no-output` flag (assuming that local outputs are not needed). Note that the `--no-ouptut` flag works in conjunction with `--push`, as the pushing takes place from the satellite. We are working on ways in which local outputs can be synchronized more intelligently such that only a diff is transferred over the network.
* A user can only be invited into an Earthly org if they already have a user account. This is a temporary limitation which will be addressed in the future.
* Satellites in conjunction with `--save-inline-cache` or `--use-inline-cache` is currently unsupported. When using `--ci`, `--save-inline-cache` and `--use-inline-cache` will not be implicitly enabled when using Satellites.
* Running Earthly `v0.6.20` against an older satellite causes an Earthly crash. We have fixed this issue and will be available in a future release, but until then you can upgrade your satellite (remove and re-launch) to work around this problem.
