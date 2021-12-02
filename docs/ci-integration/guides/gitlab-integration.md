
# GitLab CI/CD integration

Here is an example a GitLab CI/CD build, where we build the Earthly target `+build`.

```yml
# .gitlab-ci.yml

image: docker
services:
  - docker:dind

before_script:
    - apk update && apk add git
    - wget https://github.com/earthly/earthly/releases/download/v0.6.1/earthly-linux-amd64 -O /usr/local/bin/earthly
    - chmod +x /usr/local/bin/earthly
    - /usr/local/bin/earthly bootstrap
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY

earthly:
  stage: build
  script:
    - earthly --ci --push -P +build
```

A full example is available [on GitLab](https://gitlab.com/brandon176/earthly-demo/-/tree/main).

For a complete guide on CI integration see the [CI integration guide](../overview.md).
