
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

image: earthly/earthly:v0.7.22

before_script:
    - earthly bootstrap
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY

earthly:
  stage: build
  script:
    - earthly --ci --push -P +build
```

A full example is available [on GitLab](https://gitlab.com/earthly-technologies/earthly-demo).

For a complete guide on CI integration see the [CI integration guide](../overview.md).
