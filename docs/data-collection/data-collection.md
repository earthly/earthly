# Earthly data collection

By default, Earthly collects anonymized data which we use for measuring performance of the earthly command.

## Installation ID

Earthly will create a universally unique installation ID (UUID v4) under `~/.earthly/install_id`, which
is used to track each installation. This ID is randomly created and does not contain any personal data.


## Anonymized data

In addition to the installation ID, earthly will also collect a one-way-hash of the
git repository name.

## CI platform

Earthly applies some heuristics to determine if it is running in a CI system, and will
report which CI system is detected (e.g. GitHub Actions, Circle CI, Travis CI, Jenkins, etc).

## Command and exit code

Earthly will report which command was run (e.g. build, prune, etc), the execution time, and corresponding exit code.
Command line arguments are *not* captured.

## Disabling analytics

To disable the collection of data, set the `disable_analytics` option to `true` under the global config file `~/.earthly/config.yml`.

For example:

```yaml
global:
    disable_analytics: true
```

This option is documented in the [Earthly configuration file page](../earthly-config/earthly-config.md).
