Earthly was designed to be used as the main CI build specification. It can be integrated with any CI system and used instead of the CI's own native syntax.

The main benefit of using Earthly in an existing CI is that it provides consistency between local and CI builds. This allows you to debug, test, and iterate on your CI builds locally, before pushing them to the CI. The fact that Earthly uses containers underneath provides a level of guarantee that the build behaves the same no matter where it is run from. This is what we call **repeatability**.

In this section, we will explore how to use Earthly in a well-known CI system, such as GitHub Actions. For more information on how to use Earthly in other CIs such as GitLab, Jenkins, or CircleCI, you can check out the [CI Integration page](../ci-integration/overview.md).

## Using Earthly in your current CI

To use Earthly in a CI, you typically encode the following steps in your CI's build configuration:

1. Download and install Earthly
2. Log in to image registries, such as DockerHub
3. Run Earthly

For example, here is how this would work in GitHub Actions:

```yaml
```yml
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
      FORCE_COLOR: 1
    steps:
    - uses: earthly/actions/setup-earthly@v1
      with:
        version: v0.7.0
    - uses: actions/checkout@v2
    - name: Put back the git branch into git (Earthly uses it for tagging)
      run: |
        branch=""
        if [ -n "$GITHUB_HEAD_REF" ]; then
          branch="$GITHUB_HEAD_REF"
        else
          branch="${GITHUB_REF##*/}"
        fi
        git checkout -b "$branch" || true
    - name: Docker Login
      run: docker login --username "$DOCKERHUB_USERNAME" --password "$DOCKERHUB_TOKEN"
    - name: Run build
      run: earthly --push +build
```

Here is an explanation of the steps above:

* The action `earthly/actions/setup-earthly@v1` downloads and installs Earthly. Running this action is similar to running the Earthly installation one-liner `sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/download/v0.6.30/earthly-linux-amd64 -O /usr/local/bin/earthly && chmod +x /usr/local/bin/earthly'`
* The command `docker login` performs a login to the DockerHub registry. This is required, to prevent rate-limiting issues when using popular base images.
* The command `earthly --ci --push +build` executes the build. The `--ci` flag is used here, in order to force the use of `--strict` mode -- in this mode, Earthly prevents the use of features that make the build less repeatable -- and also to disable local outputs. Artifacts and images resulting from the build are not needed within the CI environment. Any outputs should be pushed via `RUN --push` or `SAVE IMAGE --push` commands.

### Using Earthly Satellites in CI (optional)

Many CI systems are based on ephemeral sandboxes - meaning that when the build finishes, the environment is thrown away, and subsequent runs will have to start again from scratch, including downloading and reinstalling dependencies. This is the case for CIs such as GitHub Actions, GitLab CI, and CircleCI. This characteristic results in many steps being repeated on every run despite the fact that nothing related to those steps has changed.

Earthly Satellites are a great way to mitigate this issue. As satellites have persistent cache that is instantly available, pipelines can execute much faster, when running together with the warmed-up cache of the remote instance. This alone typically results in a 2-20X speedup, depending on setup.

To use a satellite in your CI, you need to:

1. Generate an Earthly token and set it as an environment variable in your CI
2. Specify the satellite name to use as part of the `earthly` command.

To generate an Earthly token, assuming that you are logged in, you can use the command

```bash
earthly account create-token my-ci-token
```

You would then set the produced token as a secret in your CI of choice, and then use that secret to set an environment variable. In GitHub Actions, this looks like this:

```yml
env:
  EARTHLY_TOKEN: ${{ secrets.EARTHLY_TOKEN }}
```

Then, to use the satellite, you would use the flag `--sat my-satellite` in your `earthly` command. For example:

```bash
earthly --sat my-satellite --ci --push +build
```

For more information on how to use Earthly Satellites, you can check out the [Satellites page](../cloud/satellites.md).

For more information about integrating Earthly with other CI systems, you can check out the [CI Integration page](../ci-integration/overview.md).
