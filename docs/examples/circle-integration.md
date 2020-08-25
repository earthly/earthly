
# Circle CI integration

Here is an example of a Circle CI build, where we build the Earthly target `+build`.

```yml
# .circleci/config.yml

version: 2.1
jobs:
  build:
    machine:
      image: ubuntu-1604:201903-01
    environment:
      - GIT_URL_INSTEAD_OF: "https://github.com/=git@github.com:"
    steps:
      - checkout
      - run: docker login --username "$DOCKERHUB_USERNAME" --password "$DOCKERHUB_TOKEN"
      - run: "sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/download/v0.3.2/earth-linux-amd64 -O /usr/local/bin/earth && chmod +x /usr/local/bin/earth'"
      - run: earth --version
      - run: earth --push +build
```

For a complete guide on CI integration see the [CI integration guide](../guides/ci-integration.md).
