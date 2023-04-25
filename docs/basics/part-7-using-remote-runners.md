Earthly has the ability to run builds both locally and remotely. If you followed the standard [installation instructions](https://earthly.dev/get-earthly), then you most likely have only run local builds so far. In this section, we will explore how to use remote runners to perform builds on remote machines.

## Remote Runners

Earthly is able to use remote runners for performing builds on remote machines. When Earthly uses a remote runner, the inputs of the build are picked up from the local environment, then the execution takes place remotely, including any pushes (`RUN --push` commands, and `SAVE IMAGE --push` commands), but any local outputs are sent back to the local environment. All this takes place while your local Earthly process still provides the logs of the build in real time locally.

Remote runners are especially useful in a few specific circumstances:

* You want to reuse cache between CI runs to dramatically speed up builds (more on this in part 8).
* You want to share compute and cache with colleagues and/or with the CI.
* You have a build that requires a lot of resources, and you want to run it on a machine with more resources than your local machine.
* You have a build that requires to run on a specific CPU architecture natively.
* You have a slow internet connection.

There are two types of remote runners:

* Remote Buildkit (self-hosted)
* Earthly Satellites (managed by Earthly)

### Using a Remote Buildkit

A common way to use remote runners is to deploy your own instance of Buildkit and have Earthly connect to it. The [remote Buildkit page](../ci-integration/remote-buildkit.md) has more information on how to set this up.

Once the remote Buildkit is up and running, you may use it from Earthly by setting the configuration option `buildkit_host` to the address of the remote Buildkit. For example, if the remote Buildkit is running on `earthly-remote-buildkit.example.com`, you may set the configuration option

```bash
earthly config set global.buildkit_host earthly-remote-buildkit.example.com
```

And then run Earthly builds as usual.

```bash
earthly +my-target
```

Another option is to use the `--buildkit-host` flag on the command line, instead of setting the configuration option.

```bash
earthly --buildkit-host earthly-remote-buildkit.example.com +my-target
```

### Using an Earthly Satellite

Another way to use remote runners is to use Earthly Satellites. Earthly Satellites are remote runners managed by the Earthly team. They are a paid feature as part of the [Earthly Satellites or Earthly CI plans](https://earthly.dev/pricing).

To get started, first you need to create an Earthly Cloud account.

```bash
earthly account register --email <your-email>
```

Follow instructions in the email received to complete the registration. You will additionally need to create an organization.

```bash
earthly org create my-org
```

You must subscribe to a paid plan to use Earthly Satellites. The subscription has a 14-day trial -- your credit card is not charged if you cancel before then.

[**Click here to start your subscription**](https://buy.stripe.com/8wM9Es4BT4Vvb4YbIJ)

Then, you can create a satellite.

```bash
earthly sat launch my-satellite
```

Once a satellite has been launched it is automatically selected for use. If you ever need to switch the satellite yourself, you can can use the command...

```bash
earthly sat select my-satellite
```

Additionally, you can go back to performing local builds with the command...

```bash
earthly sat unselect
```

And then run Earthly builds as usual.

```bash
earthly +my-target
```

Or, you can use a satellite as part of the build without selecting first

```bash
earthly --sat my-satellite +my-target
```

For more information, check out the [Earthly Satellites](../cloud/satellites.md) page.

### Secrets and remote builds

When running remote builds, some operations might require access to secrets. For example, if you are pushing images to a private registry, or if you are logged in to DockerHub to prevent rate limiting. Earthly will automatically pass the credentials from your local machine to the remote runner.

Any secret that is available locally, including Docker/Podman credentials, will be passed to the remote runner whenever needed by the build.

For more information about secrets, see the [Args and secrets page](../guides/build-args.md) and the [authenticating Git and image registries page](../guides/auth.md).
