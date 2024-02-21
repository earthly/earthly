
# GitLab CI/CD integration

This example uses [GitLab CI/CD](https://docs.gitlab.com/ee/ci/) to build the Earthly target `+build`.


```yml
# .gitlab-ci.yml

services:
  - docker:dind

variables:
  DOCKER_HOST: tcp://docker:2375
  FORCE_COLOR: 1
  EARTHLY_EXEC_CMD: "/bin/sh"

image: earthly/earthly:v0.8.4

before_script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY

earthly:
  stage: build
  script:
    - earthly --ci --push -P +build
```

Note that in this particular configuration, the `earthly/earthly` image will first
start BuildKit under the same container via the image's entrypoint script; however
by setting `EARTHLY_EXEC_CMD=/bin/sh`, the `/usr/bin/earthly-entrypoint.sh` script
will present a shell rather than call the earthly binary. This bootstrapping occurs
before the `before_script` portion of the gitlab job executes.

In order to configure a registry mirror, users will need to configure a multi-line
string for `EARTHLY_ADDITIONAL_BUILDKIT_CONFIG` under the `variables` section. For example:

```yml
variables:
  EARTHLY_ADDITIONAL_BUILDKIT_CONFIG: |-
    [registry."docker.io"]
      mirrors = ["registry-mirror.example.com"]
```

A full example is available [on GitLab](https://gitlab.com/earthly-technologies/earthly-demo).

For a complete guide on CI integration see the [CI integration guide](../overview.md).
