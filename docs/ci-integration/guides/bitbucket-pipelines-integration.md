# Bitbucket Pipelines integration

Bitbucket Pipelines run in a shared Docker environment and do not support running Earthly builds directly due to [restrictions](https://jira.atlassian.com/browse/BCLOUD-21419) that Bitbucket has put in place.

You can however, run Earthly builds on Bitbucket pipelines via [remote runners](../../remote-runners.md) such as [Earthly Satellites](../../cloud/satellites.md). Because Bitbucket Pipelines run as containers you can also use the official Earthly Docker image. Here is an example of a Bitbucket Pipeline build. This example assumes your Earthfile has a `+build` target defined.

```yml
# ./bitbucket-pipelines.yml

image: earthly/earthly:v0.8.6

pipelines:
  default:
    - step:
        name: "Set Earthly token"
        script:
          - export EARTHLY_TOKEN=$EARTHLY_TOKEN
    - step:
        name: "Docker login"
        script:
          - docker login --username "$DOCKERHUB_USERNAME" --password "$DOCKERHUB_TOKEN"
    - step:
        name: "Build"
        script:
          - earthly --ci --push --sat $EARTHLY_SAT --org $EARTHLY_ORG +build
```

For a complete guide on CI integration see the [CI integration guide](../overview.md).
