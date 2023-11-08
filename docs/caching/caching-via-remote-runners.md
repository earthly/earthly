# Caching via remote runners

Caching via remote runners in Earthly works by simply reusing the same runner for multiple builds. This allows the runner to retain the cache between execution, and thus reuse it. There is nothing special that needs to be configured for this to work. All of the features of caching in Earthly work as expected, including layer caching and cache mounts. The key to what makes remote runners great is the fact that the cache is local to the runner, and thus is available instantly, without the need for an upload/download step.

To learn more, see the [remote runners page](../remote-runners.md).

The [managing cache page](./managing-cache.md) contains information about how to reset the cache remotely, if needed.
