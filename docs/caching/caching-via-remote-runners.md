# Caching via remote runners

Caching via remote runners in Earthly works by simply reusing the same runner for multiple builds. The runner retains the cache between executions, and thus is able to perform significantly better than any caching mechanism that relies on upload and download. There is nothing special that needs to be configured for this to work. All of the features of caching in Earthly work as expected, including layer caching and cache mounts.

To learn more, see the [remote runners page](../remote-runners.md).

The [managing cache page](./managing-cache.md) contains information about how to reset the cache remotely, if needed.
