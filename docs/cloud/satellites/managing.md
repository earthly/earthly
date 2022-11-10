# Managing Satellites

This feature is part of the Earthly Satellites paid plan.

{% hint style='danger' %}
##### Important

This feature is currently in **Beta** stage

* The feature may break or change significantly in future versions of Earthly.
* Give us feedback on
  * [Slack](https://earthly.dev/slack)
  * [GitHub issues](https://github.com/earthly/earthly/issues)
  * [Emailing support](mailto:support+satellite@earthly.dev)
{% endhint %}

This page describes how to manage [Earthly Satellites](../satellites.md).

## Prerequisites

In order to manage Earthly Satellites, you must have an Earthly account and an Earthly organization, and you must request access to the Satellite private beta program. For more information, see the [Earthly Cloud overview](../overview.md) and the [Satellites page](../satellites.md).

## Managing Satellites

### Launching and removing satellites

To launch a new satellite, run:

```bash
earthly sat launch <satellite-name>
```

The Satellite name can be any arbitrary string.

If you are part of multiple Earthly organizations, you may have to specify the org name under which you would like to launch the satellite:

```bash
earthly sat --org <org-name> launch <satellite-name>
```

Once the satellite is created it will be automatically selected for use as part of your builds. The selection takes place by Earthly adding some information in your Earthly config file (usually located under `~/.earthly/config.yml`).

To remove a satellite, you can run:

```bash
earthly sat rm <satellite-name>
```

### Listing satellites

To list the satellites available in your organization, run:

```bash
earthly sat ls
```

### Selecting a satellite

Selecting a satellite causes Earthly to use that satellite for any builds from that point onwards.

To select a satellite for use, run:

```bash
earthly sat select <satellite-name>
```

### Unselecting a satellite

Unselecting a satellite will cause Earthly to run builds locally from that point onwards.

To unselect a satellite, run:

```bash
earthly sat unselect
```

### Checking status of a satellite

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

### Clearing cache

To clear the cache of a satellite, run the following while a satellite is selected:

```bash
earthly prune -a
```

### Upgrading a satellite

Currently, satellites do not have an auto-update mechanism built in. In order to get a newer version of a satellite, you need to manually remove and re-launch the satellite. Note that this operation resets the cache.

```bash
earthly sat rm <satellite-name>
earthly sat launch <satellite-name>
```

The newly launched satellite will always get the latest version available.

### Managing instance state

To save costs, satellites automatically enter a **sleep** state after 30 min of inactivity. While a satellite is asleep, you are not billed for any compute minutes.

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

### Inviting a user to use a satellite

Currently, all users who are part of an organization are allowed to use any satellite in the organization. To invite another user to join your org, run:

```bash
earthly org invite /<org-name>/ <email>
```

Note the slashes around the org name. Also, please note that **the user must have an account on Earthly before they can be invited**. (This is a temporary limitation which will be addressed in the future.)

Once a user has been invited, you can forward them a link to the page [Using Satellites](./using.md) for them to get started.

### Satellite IP address

The source IP address of the satellite for all internet traffic is `35.160.176.56`. This can be used for granting access to private resources or to production environments.
