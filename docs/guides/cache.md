# Managing cache

Earthly cache works similarly to Dockerfile layer-based caching. In fact, the same [technology](https://github.com/moby/buildkit) is used underneath.

## Cache layers

Many Earthfile commands create cache layers. A cache layer may be reused in a future build, if the conditions under which it is created are the same.

Examples of commands which create layers are `COPY` and `RUN`.

One of the main things influencing the conditions are "sources". Sources are created through commands like `COPY` and `GIT CLONE`. `RUN`, however, is not a source, even if the command itself involves downloading content from an external location. This means that a `RUN` command, on its own, would always be cached if it has been run under the same circumstances previously (except for the `RUN --push` variant).

For a primer into Dockerfile caching see [this article](https://pythonspeed.com/articles/docker-caching-model/). The same principles apply to Earthfiles.

## Cache location

Earthly cache is persisted in a docker volume called `earthly-cache` on your system. When Earthly starts for the first time, it brings up a BuildKit daemon in a Docker container, which initializes the `earthly-cache` volume. The volume is managed by Earthly's BuildKit daemon and there is a regular garbage-collection for old cache to keep this space at a maximum of approximately 10GB.

## Specifying cache size

Some builds may require more cache beyond the default 10GB allocated. In order to modify the size of the cache, you can modify the `cache_size_mb` setting in the [configuration](../earth-config/earth-config.md). For example:

```yaml
global:
  cache_size_mb: 50000
```

## Resetting cache

The cache can be safely deleted manually, if the daemon is not running

```bash
docker stop earthly-buildkitd
docker rm earthly-buildkitd
docker volume rm earthly-cache
```

However, it is easier to simply use the command

```bash
earth prune --reset
```

which restarts the daemon and resets the contents of the cache volume.

## See also

* [Advanced caching techniques](./advanced-caching.md)
