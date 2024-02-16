# Managing Satellites

This page describes how to manage [Earthly Satellites](../satellites.md).

## Launching a new satellite

Satellites are launched in one of the following two ways, depending on which kind of satellite you intend on creating.

### Earthly Cloud

To launch a new satellite on Earthly Cloud, run:

```bash
earthly sat launch <satellite-name>
```

The Satellite name can be any arbitrary string.

If you are part of multiple Earthly organizations, you may have to specify the org name under which you would like to launch the satellite:

```bash
earthly sat launch <satellite-name>
```

Once the satellite is created it will be automatically selected for use as part of your builds. The selection takes place by Earthly adding some information in your Earthly config file (usually located under `~/.earthly/config.yml`).

### Self-Hosted

Self-Hosted Satellites are created by running the satellite container directly. See the [self-hosted guide](self-hosted.md) for instructions.

## Removing a satellite

Satellites are removed in one of the following two ways, depending on where they are deployed.

### Earthly Cloud

To remove a satellite from Earthly Cloud, you can run:

```bash
earthly sat rm <satellite-name>
```

Note that it is best to remove a satellite while it is asleep, to prevent accidentally cancelling an ongoing build.

### Self-Hosted

Self-Hosted Satellites are typically removed by gracefully terminating the satellite container directly. Once the satellite has terminated, it will continue listing in an `offline` state in the output of `satellite ls`. This record is for historical or debugging purposes, however, it can be permanently removed by running `satellite rm <name>`.

See the [self-hosted guide](self-hosted.md) for more details.

## Listing satellites

To list the satellites available in your organization, run:

```bash
earthly sat ls
```

## Selecting a satellite

Selecting a satellite causes Earthly to use that satellite for any builds from that point onwards.

To select a satellite for use, run:

```bash
earthly sat select <satellite-name>
```

## Unselecting a satellite

Unselecting a satellite will cause Earthly to run builds locally from that point onwards.

To unselect a satellite, run:

```bash
earthly sat unselect
```

## Checking status of a satellite

Checking the status of a satellite allows you to view information about a satellite's current state, including whether it is being used right now, how much cache space has been used, version information and other information.

To check the status of a satellite, you can run:

```bash
earthly sat inspect <satellite-name>
```

Here is some example output of an inspect command:

```
Connecting to core-test...
Version github.com/earthly/buildkit v0.6.21 7a6f9e1ab2a3a3ddec5f9e612ef390af218a32bd
Platforms: linux/amd64 (native) linux/amd64/v2 linux/amd64/v3 linux/amd64/v4 linux/arm64 linux/riscv64 linux/ppc64le linux/s390x linux/386 linux/mips64le linux/mips64 linux/arm/v7 linux/arm/v6
Utilization: 0 other builds, 0/12 op load
GC stats: 9.0 GB cache, avg GC duration 275ms, all-time GC duration 2.754s, last GC duration 0s, last cleared 0 B
Instance state: Operational
Currently selected: No
```

## Clearing cache

There are two ways to clear the cache on a satellite. One may be faster than the other, depending on the size of the cache.

### Recreating the Underlying Satellite Instance (often faster)

Running the `update` command with `--drop-cache` will relaunch the instance with an empty cache volume.
Note that this operation can take a while, and the satellite may also receive any available updates during the process.

```bash
earthly satellite update --drop-cache my-satellite
```

Note: the `update` comamnd only works with Earthly Cloud satellites.

### Using the prune command

The `earthly prune` command also works on satellites.
It usually takes longer than running `satellite update`; however, it does not trigger a relaunch.
The prune command requires the satellite to be selected before running.

```bash
earthly prune -a
```

## Updating a satellite

The steps below apply only to Earthly Cloud satellites. For [Self-Hosted Satellites](self-hosted.md), you must upgrade the version being used in your deployment manually.

Earthly Cloud Satellites receive automatic version updates by default, unless they are pinned to a specific version (by using the `--version` launch flag).
Pinned versions will still receive minor patches, such as security updates.

### Auto-Update Maintenance Windows

Maintenance windows are set between 2AM and 4AM in the timezone where the satellite is launched by default.
The start time of the 2 hour window can be explicitly set by passing a 24-hr formatted time to the `--maintenance-window` flag of the launch command.

Here's an example showing a satellite launched with a custom maintenance window of 4AM to 6AM:

```bash
earthly satellite launch --maintenance-window 04:00 my-satellite
```

Note that updates will only happen during the maintenance window while the satellite is asleep. If the satellite remains in use 
for the entire duration of the maintenance window, then the update will be re-attempted the next day.

### Version Pinning

If you want to prevent your satellite from automatically upgrading to a new earthly version, you can pin your version using the `--version` flag.
Note that satellites on pinned versions may still receive auto-updates during a maintenance window; however, these updates will be limited to stability or security patches, rather than version updates.

```bash
earthly satellite launch --version v0.6.29 my-satellite
```

### Weekends-Only Mode

In cases where nightly maintenance schedules are not suitable, you can configure your auto-updates so that they will only run on Saturday using the following launch flag:

```bash
earthly satellite launch --maintenance-weekends-only my-satellite
```

{% hint style='info' %}
##### Note
The day of the week is determined based on the UTC timezone.
Maintenance windows will be expanded to better incorporate user timezones and custom schedules in the future.
{% endhint %}

### Manually Updating a Satellite

Satellites can also be manually updated using the `update` command. The update command can be used to not only trigger a version upgrade, but also to change other parameters of the satellite, such as feature-flags or its cache. 
Satellites must be in a sleep state before an update can be started. You can use the `earthly satellite sleep` command to do this manually.
Below are some examples of how you can use the `update` command.

The following example updates a satellite to the latest revision, respecting any pinned versions:

```bash
earthly satellite update my-satellite
```

#### Dropping Cache

The following command updates the satellite and clears its cache during the process. If no updates are available, this command will still clear the cache.

```bash
earthly satellite update --drop-cache my-satellite
```

#### Changing Earthly Version

Updates can also specify a new pinned earthly version using the `--version` flag:

```bash
earthly satellite update --version v0.6.29 my-satellite
```

#### Changing Satellite Size

The size of an existing satellite can be altered using the `--size` flag.
Note that changing the size of a satellite will also drop its existing cache.

```
earthly satellite update --size xlarge my-satellite
```

#### Changing Feature Flags

Feature flags can be set during an update as well.
When any feature flags are passed in, the entire set of existing feature flags are replaced with the new set.
Passing no feature flags will retain the existing flags.

Note feature-flags are typically used to preview unreleased features. They are considered highly experimental.

```bash
earthly satellite update --feature-flag cache-pct=30 my-satellite
```

{% hint style='info' %}
##### Note
It's not currently possible to completely clear out the flags using the update command; you will have to destroy and recreate the satellite.
{% endhint %}

#### Satellite Revision System

In more detail, satellite versions are controlled via Earthly's internal revisioning system. 
A satellite revision includes an earthly version plus a revision increment, where each earthly version may contain a number or ordered revisions.
Revision increments are released to patch stability or performance within a specific earthly version.

You can view your current satellite version and revision number using the `earthly satellite inspect` command.


## Managing instance state

To save costs, Earthly Cloud satellites automatically enter a **sleep** state after 30 min of inactivity. While a satellite is asleep, you are not billed for any compute minutes.

The satellite will automatically **wake up** when a new build is started while it's in a sleep state. This is visible during the `Init` phase of the Earthly log.

If you want more fine-grain control over your Satellite's state, you can also manually put it to sleep using the command:

```bash
earthly sat sleep <satellite-name>
```

Similarly, a Satellite can be manually woken up using:

```bash
earthly sat wake <satellite-name>
```

Note that the [`inspect`](#checking-status-of-a-satellite) command will show you if a Satellite is currently awake or asleep.

## Inviting a user to use a satellite

Currently, all users who are part of an organization are allowed to use any satellite in the organization. To invite another user to join your org, run:

```bash
earthly org invite <email>
```

Once a user has been invited, you can forward them a link to the page [Using Satellites](./using.md) for them to get started.

## Satellite IP address

The source IP address of an Earthly Cloud satellite for all internet traffic is `35.160.176.56`. This can be used for granting access to private resources or to production environments.
