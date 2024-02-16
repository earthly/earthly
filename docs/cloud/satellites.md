# Earthly Satellites

Earthly Satellites are [remote runners](../remote-runners.md) that work seamlessly with Earthly, using persistent cache to improve build times.
Satellites can be either [fully managed](https://earthly.dev/earthly-satellites) by Earthly Cloud or [self-hosted](./satellites/self-hosted.md) in your own environment.

## Get started with Earthly Cloud Satellites for free!

Fully managed Satellites are included with [Earthly Cloud](https://docs.earthly.dev/earthly-cloud/overview). Earthly Cloud is a SaaS build automation platform with consistent builds, ridiculous speed, and a next-gen developer experience that works seamlessly with any CI. *Get 6,000 build minutes/month as part of Earthly Cloud's no time limit free tier.* ***[Sign up today](https://cloud.earthly.dev/login).***

## Benefits

Typical use cases for Earthly Satellites include:

* **Speeding up CI builds** in sandboxed CI environments such as GitHub Actions, GitLab, CircleCI, and others. Most CI build times are improved by 2-20X with Satellites.
* **Sharing compute and cache with coworkers** or with the CI.
* **Executing cross-platform builds natively**. For example, executing builds on x86 architecture natively when you are working from an Apple Silicon machine (Apple M1/M2) and vice versa, arm64 builds from an x86 machine.
* **Benefiting from high-bandwidth internet access** from the satellite, allowing for fast downloads of dependencies and pushes for deployments. This is particularly useful if you are in a location with slow internet.
* **Using Earthly in restricted environments**, where privileged access or docker-in-docker are not supported.

## How Earthly Satellites work

### On your laptop

* You kick off the build from the command line, and Earthly uses a remote satellite for execution.
* The source files used are the ones you have locally in the current directory.
* The build logs from the satellite are streamed back to your terminal in real time, so you can see the progress of the build.
* The outputs of the build - images and artifacts - are downloaded back to your local machine upon success.
* Everything looks and feels as if it is executing on your computer in your terminal.
* In reality, the execution takes place in the cloud with high parallelism and a lot of caching.

### In your CI of choice

* The CI starts a build and invokes Earthly.
* Earthly starts the build on a remote satellite, executing each step in isolated containers.
* The same cache is used between runs on the same satellite, so parts that haven’t changed do not repeat.
* Logs are streamed back to the CI in real time.
* Any images, artifacts, or deployments that need to be pushed as part of the build are pushed directly from the satellite.
* Build pass/fail is returned as an exit code, so your CI can report the status accordingly.

## Getting started

### 1. Sign up for Earthly Cloud (free)

Earthly Satellites is part of Earthly Cloud. You can use it for free as part of our free tier. Get started with Earthly Cloud by visiting the [sign up](https://cloud.earthly.dev/login) page, and get 6,000 build minutes/month for free.

### 2. Launch a new satellite

Satellites are launched in one of the following two ways, depending on which kind of satellite you intend on creating.

#### Earthly Cloud

To launch a new managed Satellite on Earthly Cloud, run:

```bash
earthly sat launch <satellite-name>
```

The Satellite name can be any arbitrary string.

If you are part of multiple Earthly organizations, you may want to first select the org under which you would like to launch the satellite:

```bash
earthly org select <org-name>
earthly sat launch <satellite-name>
```

Once the satellite is created it will be automatically selected for use as part of your builds. The selection takes place by Earthly adding some information in your Earthly config file (usually located under `~/.earthly/config.yml`).

#### Self-Hosted

Self-Hosted Satellites are instead launched by running the satellite container directly. See the [self-hosted guide](./satellites/self-hosted.md) for instructions.

### 3. Run a build

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

For more information on using satellites, see the [Using satellites page](./satellites/using.md).

### 4. Invite your team

A final step is to invite your team to use the satellite. This can be done by running:

```bash
earthly org invite <email>
```

Once a user has been invited, you can forward them a link to the page [Using Satellites](./satellites/using.md) for them to get started.

## Managing Satellites

For more information on managing satellites, see the [Managing Satellites page](./satellites/managing.md).

## Satellite specs

When using Cloud Satellites, the size and architecture can be specified at launch time using the `--size` and `--platform` flags.
For the full list of supported options, please see the [Pricing Page](https://earthly.dev/pricing).

## Using Satellites in CI

A key benefit of using satellites in a CI environment is that the cache is shared between runs. This results in significant speedups in CIs that would otherwise have to start from scratch each time.

{% hint style='danger' %}
##### Note

If a satellite is shared between multiple CI pipelines, it is possible that it becomes overloaded by too many parallel builds. For best performance, you can create a dedicated satellite for each CI pipeline. See the [best practices guide](./satellites/best-practices.md) for more details.
{% endhint %}

To get started with using Earthly Satellites in CI, you can create a login token for access.

First, run

```bash
earthly account create-token <token-name>
```

to create your login token.

Copy and paste the value into an environment variable called `EARTHLY_TOKEN` in your CI environment.

Then as part of your CI script, simply select your satellite using one of these supported methods

* Selection command: `earthly sat select <satellite-name>`
* Setellite flag: `earthly --sat my-satellite +build`
* Environment variable: `EARTHLY_SATELLITE=my-satellite`

before running your Earthly targets.

Note that when using [Self-Hosted Satelites](./satellites/self-hosted.md), your CI runner must be able to access the satellite on the network where it is hosted.

{% hint style='danger' %}
##### Registry Login

It's best to avoid using an image registry like Dockerhub without authentication, since the IP address from the satellite easily become rate-limited.
A simple `docker login` command before running your build should be used to pass registry credentials to your satellite.
See our [Docker authentication](../guides/auth.md) guide for more details.

{% endhint %}

## Known limitations

* Pull-through cache is currently not supported

If you run into any issues please let us know either via [Slack](https://earthly.dev/slack), [GitHub issues](https://github.com/earthly/cloud-issues/issues) or by [emailing support](mailto:support+satellite@earthly.dev).
