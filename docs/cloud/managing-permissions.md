# Managing permissions

This page covers managing permissions for Earthly Cloud products, such as Earthly Cloud Secrets, and Earthly Satellites.

## Overview

Earthly Cloud's permissions model has two security boundaries:

* Earthly orgs
* Earthly projects

Users may be invited to a specific organization, and optionally to specific projects within the organization.

Earthly orgs may contain the following shared resources:

* Satellites
* Earthly projects

Earthly projects, in turn, may contain the following resources:

* Secrets
* Pipelines
* Build history, including build logs

## Earthly org access levels

Within an Earthly org, users may be granted one of the following access levels:

* `read`: Can view the org, projects, and user membership. Can view, inspect, wake and build on satellites. Can also stream and share logs.
* `read+secrets`: Same as read, but can also view and use secrets.
* `write`: Everything in `read+secrets`, plus the ability to create and modify satellites, projects, and secrets.
* `admin`: Can manage the org, including adding and removing users, and managing projects, secrets and satellites.

Having a certain level of access for a given org automatically grants the same level of access to all projects within that org.

### Managing access to an Earthly org

To grant access to an Earthly org, you must invite the user to the org. This can be done by running:

```bash
earthly org invite --permission <access-level> <email>
```

If the user is already part of the org, you can change their access level by running:

```bash
earthly org member update --permission <permission> <email>
```

If you want to revoke access to an Earthly org, you can do so by running:

```bash
earthly org member rm <email>
```

## Earthly project access levels

Within an Earthly project, users may be granted one of the following access levels:

* `read`: Can view the project, including the build history and build logs.
* `read+secrets`: Same as read, but can also view and use secrets.
* `write`: Everything in `read+secrets`, plus the ability to create and modify secrets.
* `admin`: Everything in `write`, plus the ability to manage the project's users.

### Managing access to an Earthly project

To grant access to an Earthly project, you must invite the user to the project. This can be done by running:

```bash
earthly project --project <project-name> member add --permission <access-level> <email>
```

{% hint style='info' %}
##### Note
You can only invite a user to a project if they are already part of the organization.
{% endhint %}

If the user is already part of the project, you can change their access level by running:

```bash
earthly project --project <project-name> member update --permission <permission> <email>
```

If you want to revoke access to an Earthly project, you can do so by running:

```bash
earthly project --project <project-name> member rm <email>
```
