The `dind` (docker-in-docker) image is designed for Earthfile targets that use the `WITH DOCKER` command.

See the ["use-earthly-dind" best-practice](https://docs.earthly.dev/best-practices#use-earthly-dind) for details.

## Tags

* `alpine-3.18-docker-23.0.6-r4`
* `ubuntu-20.04-docker-24.0.5-1`
* `ubuntu-23.04-docker-24.0.5-1`

## Outdated Tags

* `alpine`
* `ubuntu`

## Note

The `alpine` and `ubuntu` images are incompatible with the earthly v0.7.14 onwards; newer versions of earthly should use the tags that contain an OS specific version in them.

To ease this transition, one can make use of an `IF` command that depends on the `EARTHLY_VERSION` builtin argument:

```
VERSION 0.7

dind:
  FROM earthly/dind:alpine
  ARG EARTHLY_VERSION
  ARG SMALLEST_VERSION="$(echo -e "$EARTHLY_VERSION\nv0.7.14" | sort -V | head -n 1)"
  IF [ "$SMALLEST_VERSION" = "v0.7.14" ]
    # earthly is at v0.7.14 or newer, and must use the more recent dind:alpine-3.18-docker-23.0.6-r4 image
    FROM earthly/dind:alpine-3.18-docker-23.0.6-r4
  END

test:
  FROM +dind
  WITH DOCKER
    RUN docker --version # old versions of earthly will get 20.10.14, and newer will get 23.0.6
  END
```
