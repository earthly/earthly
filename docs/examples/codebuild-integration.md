
# AWS CodeBuild integration

Here is an example of an AWS CodeBuild build, where we build the Earthly target `+build`.

{% hint style='info' %}
##### Note

Ensure when you're creating your CodeBuild Project that you enable the `Privileged` flag
in order to allow Earthly build Docker images.

{% endhint %}

```yml
# ./buildspec.yml
version: 0.2

env:
  variables:
    GIT_URL_INSTEAD_OF: "https://github.com/=git@github.com:"

phases:
  install:
    commands:
      - wget https://github.com/earthly/earthly/releases/latest/download/earth-linux-amd64 -O /usr/local/bin/earth && chmod +x /usr/local/bin/earth
  pre_build:
    commands:
      - echo Logging in to Docker
      - docker login --username "$DOCKERHUB_USERNAME" --password "$DOCKERHUB_TOKEN"
  build:
    commands:
      - earth --version
      - earth --push +build
```

For a complete guide on CI integration see the [CI integration guide](../guides/ci-integration.md).
