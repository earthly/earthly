# Experimental Features

Earthly makes use of feature flags to release new and experimental features.
These features must be explicitly enabled to use them.

Earthly uses [semantic versioning](http://semver.org/); once a new feature
has reached stability, a new major or minor version of Earthly will be released with
the feature enabled by default.

## Specifying Version and features

Each earthfile should list the current earthly version it depends on using the [`VERSION`](../earthfile/earthfile.md#version) command.
The `VERSION` command was first introduced under `0.5` and is currently optional; however, it will become manditory in a future version.

```Dockerfile
VERSION [--use-copy-include-patterns] <version-number>
```

## Feature flags

| Feature flag | status | description |
| --- | --- | --- |
| `--use-copy-include-patterns` | experimental | speeds up COPY transfers |

##### `--use-copy-include-patterns`

*Speeds up COPY transfers.*

When enabled, Earthly will only send the files listed for the specific [`COPY`](../earthfile/earthfile.md#copy) command.
Without this feature, Earthly sends the entire directory of files excluding files listed in the [`.earthignore` file](../earthfile/earthignore.md).
