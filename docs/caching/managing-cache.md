# Managing cache

This page describes how to manage the Earthly cache locally or on a remote runner, such as an Earthly Satellite.

## Local cache

### Local cache location

Earthly cache is persisted in a docker (or podman) volume called `earthly-cache` on your system. When Earthly starts for the first time, it brings up a BuildKit daemon in a Docker container, which initializes the `earthly-cache` volume. The volume is managed by Earthly's BuildKit daemon and there is a regular garbage-collection for old cache.

{% hint style='info' %}
#### Checking current cache size
You can check the current size of the cache by running:

```bash
sudo du -h /var/lib/docker/volumes/earthly-cache | tail -n 1
```
{% endhint %}

### Specifying local cache size limit

The default cache size is adaptable depending on available space on your system. It defaults to 10% or 10 GB, whichever is greater. If you would like to change the cache size, you can specify a different limit by modifying the `cache_size_mb` and/or `cache_size_pct` settings in the [configuration](../earthly-config/earthly-config.md). For example:

```yaml
global:
  cache_size_mb: 30000
  cache_size_pct: 70
```

### Resetting the local cache

The cache can be safely deleted manually, if the daemon is not running

```bash
docker stop earthly-buildkitd
docker rm earthly-buildkitd
docker volume rm earthly-cache
```

However, it is easier to simply use the command

```bash
earthly prune --reset
```

which restarts the daemon and resets the contents of the cache volume.

## Cache on a remote runner

### Configuring the cache size on a remote runner

If you are using [Earthly Satellites](../cloud/satellites.md), you can simply launch a bigger satellite via the `--size` flag: `earthly sat launch --size ...`.

If you are using a self-hosted remote runner, you can configure the cache policy by passing the appropriate [buildkit configuration](https://github.com/moby/buildkit/blob/master/docs/buildkitd.toml.md) to the buildkitd container.

### Resetting the cache on a remote runner

The command `earthly prune` will work on remote runners too, albeit without the `--reset` flag, which is not supported in a remote setting.

To cause a satellite to restart with a fresh cache, you can use the command `earthly sat update --drop-cache`.
