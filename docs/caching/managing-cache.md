# Managing cache

This page describes how to manage the Earthly cache locally or on a remote runner, such as an Earthly Satellite.

## Local cache

### Local cache location

Earthly cache is persisted in a docker (or podman) volume called `earthly-cache` on your system. When Earthly starts for the first time, it brings up a BuildKit daemon in a Docker container, which initializes the `earthly-cache` volume. The volume is managed by Earthly's BuildKit daemon and there is a regular garbage-collection for old cache.

### Specifying the local cache size limit

The default cache size is adaptable depending on available space on your system. It defaults to `min(55%, max(10%, 20GB))`. If you would like to change the cache size, you can specify a different limit by modifying the `cache_size_mb` and/or `cache_size_pct` settings in the [configuration](../earthly-config/earthly-config.md). For example:

```yaml
global:
  cache_size_mb: 30000
  cache_size_pct: 70
```

{% hint style='info' %}
#### Checking current size of the cache volume
You can check the current size of the cache volume by running:

```bash
sudo du -h /var/lib/docker/volumes/earthly-cache | tail -n 1
```
{% endhint %}

### Resetting the local cache

To reset the cache, you can issue the command

```bash
earthly prune
```

You can also safely delete the cache manually, if the daemon is not running

```bash
docker stop earthly-buildkitd
docker rm earthly-buildkitd
docker volume rm earthly-cache
```

Earthly also has a command that automates the above:

```bash
earthly prune --reset
```

## Cache on a remote runner / Earthly Satellite

### Configuring the cache size on a remote runner

If you are using [Earthly Satellites](../cloud/satellites.md), you can simply launch a bigger satellite via the `--size` flag: `earthly sat launch --size ...`.

If you are using a self-hosted remote runner, you can configure the cache policy by passing the appropriate [buildkit configuration](https://github.com/moby/buildkit/blob/master/docs/buildkitd.toml.md) to the [buildkit container](../ci-integration/remote-buildkit.md).

### Resetting the cache on a remote runner

The command `earthly prune` will work on remote runners too, albeit without the `--reset` flag, which is not supported in a remote setting.

To cause a satellite to restart with a fresh cache, you can use the command `earthly sat update --drop-cache`.

## Auto-skip cache

The auto-skip cache is a cache that is used to skip large parts of a build in certain situations. It is used by the `earthly --auto-skip` and `BUILD --auto-skip` commands.

Unlike the layer cache and the cache mounts, the auto-skip cache is global and is stored in a cloud database.

To clear the entire auto-skip cache for your Earthly org, you can use the command `earthly prune-auto-skip`.

To clear the auto-skip cache for an entire repository, you can use the command `earthly prune-auto-skip --path github.com/foo/bar --deep`.

To clear the auto-skip cache for a specific target, you can use the command `earthly prune-auto-skip --path github.com/foo/bar --target +my-target`.
