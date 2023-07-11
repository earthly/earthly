Earthly was designed to be used as the main CI build specification. It can be integrated with any CI system and used instead of the CI's own native syntax.

The main benefit of using Earthly in CI is that it provides consistency between local and CI builds. This allows you to debug, test, and iterate on your builds locally, before pushing them to the CI. The fact that Earthly uses containers underneath provides a level of guarantee that the build behaves the same no matter where it is run from. This is what we call **repeatability**.

In this section, we will explore how to use Earthly in a traditional CI system, such as GitHub Actions. **If you would like an introduction to using Earthly CI, see [Part 8b](./part-8b-using-earthly-ci.md) of this tutorial.**

For more information on how to use Earthly in other CIs such as GitLab, Jenkins, or CircleCI, you can check out the [CI Integration page](../ci-integration/overview.md).

## Using Earthly in Your Current CI

To use Earthly in a CI, you typically encode the following steps in your CI's build configuration:

1. Download and install Earthly
2. Set up any credentials needed for the build
3. Log in to image registries, such as DockerHub
4. Run Earthly

As part of this, you may need to set up credentials for Earthly Cloud, if you are using Earthly Satellites or Earthly Secrets. For this, you can use the following command:

```bash
earthly account create-token my-ci-token
```

Finally, here is a complete example of how to run Earthly in GitHub Actions:

```yaml
# .github/workflows/ci.yml

name: CI

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    env:
      DOCKERHUB_USERNAME: ${{ secrets.DOCKERHUB_USERNAME }}
      DOCKERHUB_TOKEN: ${{ secrets.DOCKERHUB_TOKEN }}
      EARTHLY_TOKEN: ${{ secrets.EARTHLY_TOKEN }}
      FORCE_COLOR: 1
    steps:
    - uses: earthly/actions/setup-earthly@v1
      with:
        version: v0.7.10
    - uses: actions/checkout@v2
    - name: Docker Login
      run: docker login --username "$DOCKERHUB_USERNAME" --password "$DOCKERHUB_TOKEN"
    - name: Run build
      run: earthly --org <org-name> --sat <satellite-name> --ci --push +build
```

Here is an explanation of the steps above:

* The action `earthly/actions/setup-earthly@v1` downloads and installs Earthly. Running this action is similar to running the Earthly installation one-liner `sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/download/v0.7.11/earthly-linux-amd64 -O /usr/local/bin/earthly && chmod +x /usr/local/bin/earthly'`
* The command `docker login` performs a login to the DockerHub registry. This is required, to prevent rate-limiting issues when using popular base images.
* The command `earthly --org ... --sat ... --ci --push +build` executes the build. The `--ci` flag is used here, in order to force the use of `--strict` mode. In `--strict` mode, Earthly prevents the use of features that make the build less repeatable and also disables local outputs -- because artifacts and images resulting from the build are not needed within the CI environment. Any outputs should be pushed via `RUN --push` or `SAVE IMAGE --push` commands. The flags `--org` and `--sat` allow you to select the organization and satellite to use for the build. If no satellite is specified, the build will be executed in the CI environment itself, with limited caching.

For more information about integrating Earthly with other traditional CI systems, you can check out the [CI Integration page](../ci-integration/overview.md).
