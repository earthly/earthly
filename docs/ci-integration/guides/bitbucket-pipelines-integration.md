# Bitbucket Pipelines integration

Bitbucket Pipelines run in a shared Docker environment and do not support running Earthly builds directly due to [restrictions](https://jira.atlassian.com/browse/BCLOUD-21419) that Bitbucket has put in place.

You can however, run Earthly builds on Bitbucket pipelines via [Earthly Satellites](../../cloud/satellites.md). Because Bitbucket Pipelines run as containers you can also use the official Earthly Docker image. Here is an example of a Bitbucket Pipeline build. This example assumes `+unit-test` and `+build` targets in your Earthfile.

```yml
# ./bitbucket-pipelines.yml

image: earthly/earthly:latest

pipelines:
  default:
    - step:
        name: "Set Earthly token"
        script:
          - export EARTHLY_TOKEN=$EARTHLY_TOKEN
          # See https://docs.earthly.dev/docs/earthly-command#earthly-account-create-token to obtain a token.
          - earthly --version
    - step:
        name: "Docker login"
        script:
          - docker login --username "$DOCKERHUB_USERNAME" --password "$DOCKERHUB_TOKEN"
    - step:
        name: "Unit tests"
        script:
          - earthly --sat $EARTHLY_SAT --org $EARTHLY_ORG +unit-test
    - step:
        name: "Build"
        script:
          - earthly --push --sat $EARTHLY_SAT --org $EARTHLY_ORG +build
```

For a complete guide on CI integration see the [CI integration guide](../overview.md).

{% hint style='danger' %}

## earthly/earthly:latest Docker image

The example uses the `:latest` tag for the Earthly image. In production it is recommended to specify a specific Earthly version instead of using latest to mitigate the risks of potentially breaking changes as new Earthly images are pushed.
