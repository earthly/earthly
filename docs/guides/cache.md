# Managing cache

Earthly cache works similarly to Dockerfile layer-based caching. In fact, the same [technology](https://github.com/moby/buildkit) is used underneath.

## Cache layers

Many Earthfile commands create cache layers. A cache layer may be reused in a future build, if the conditions under which it is created are the same.

Examples of commands which create layers are `COPY` and `RUN`.

One of the main things influencing the conditions are "sources". Sources are created through commands like `COPY` and `GIT CLONE`. `RUN`, however, is not a source, even if the command itself involves downloading content from an external location. This means that a `RUN` command, on its own, would always be cached if it has been run under the same circumstances previously (except for the `RUN --push` variant).

For a primer into Dockerfile caching see [this article](https://pythonspeed.com/articles/docker-caching-model/). The same principles apply to Earthfiles.

## Cache location

Earthly cache is persisted in a directory located at `/tmp/earthly` on your system. When Earthly starts for the first time, it brings up a BuildKit daemon in a Docker container, which reserves some disk space for the cache (by default, 10GB) as a loop device.

## Resetting cache

The cache can be safely deleted manually, if the daemon is not running

```bash
docker stop earthly-buildkitd
sudo rm -rf /tmp/earthly
```

Or you can also issue the earth command

```bash
earth prune --reset
```

which restarts the daemon and resets the cache directory.
