The `dind` (docker-in-docker) image is designed for Earthfile targets that use the `WITH DOCKER` command.

See the ["use-earthly-dind" best-practice](https://docs.earthly.dev/best-practices#use-earthly-dind) for details.

## Tags

This image supports 3 linux distributions:
* alpine
* ubuntu:20.04
* ubuntu:23.04

For which the current latest tags (respectively) are:
* `alpine-3.19-docker-25.0.2-r0`
* `ubuntu-20.04-docker-24.0.5-1`
* `ubuntu-23.04-docker-24.0.5-1`

For other available tags, please check out https://hub.docker.com/r/earthly/dind/tags

## Outdated Tags

* `alpine`
* `ubuntu`

## Note

The outdated `ubuntu` image is incompatible with the earthly v0.7.14 (and fixed in v0.7.15).
Correspondingly the `alpine` image at one point was also incompatible with v0.7.14, but was updated with
a backwards-compatable fix.

Users, however, are encouraged to pin to specific version tags moving forward. The unversioned tags will be left as-is
to help backwards-breaking changes.

To ease this transition, one can make use of an `IF` command that depends on the `EARTHLY_VERSION` builtin argument:

```
VERSION 0.8

dind:
  FROM earthly/dind:alpine
  ARG EARTHLY_VERSION
  ARG SMALLEST_VERSION="$(echo -e "$EARTHLY_VERSION\nv0.7.14" | sort -V | head -n 1)"
  IF [ "$SMALLEST_VERSION" = "v0.7.14" ]
    # earthly is at v0.7.14 or newer, and must use the more recent dind:alpine-3.19-docker-25.0.2-r0 image
    FROM earthly/dind:alpine-3.19-docker-25.0.2-r0
  END

test:
  FROM +dind
  WITH DOCKER
    RUN docker --version # old versions of earthly will get 20.10.14, and newer will get 23.0.6
  END
```
