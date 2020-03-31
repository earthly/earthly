# CI integration

Integrating Earthly into your CI is simply a matter of automating the same steps you would use for your local installation. In this guide, we will walk through this process.

## Step 1: Ensure Docker and Git are available

The first step is to ensure that Earthly's pre-requisites are available. On many CI systems both of these are available in the default base image or environment. Refer to your provider's documentation.

Vendors known to include these dependencies:

* CircleCI image `ubuntu-1604:201903-01`
* GitHub actions `ubuntu-latest`
* Travis dist `xenial`
* GitLab image `docker:git` with service `docker:dind` added.
* Azure DevOps vmImage `Ubuntu-16.04`

## Step 2: Install earth command

The next step is to install the `earth` command. For this, you need to run the command:

```bash
sudo /bin/sh -c 'curl -s https://api.github.com/repos/vladaionescu/earthly-releases/releases/latest | grep browser_download_url | grep linux-amd64 | cut -d : -f 2- | tr -d \x22 | wget -P /usr/local/bin/ -i - && mv /usr/local/bin/earth-linux-amd64 /usr/local/bin/earth && chmod +x /usr/local/bin/earth'
```

Note that if you are using YAML to specify your CI, you should put this command in quotes. Note also that `\x22` in this command then becomes `\\x22` in a quoted YAML string. For example:

```yml
run: "sudo /bin/sh -c 'curl -s https://api.github.com/repos/vladaionescu/earthly-releases/releases/latest | grep browser_download_url | grep linux-amd64 | cut -d : -f 2- | tr -d \\x22 | wget -P /usr/local/bin/ -i - && mv /usr/local/bin/earth-linux-amd64 /usr/local/bin/earth && chmod +x /usr/local/bin/earth'"
```

## Step 3: Configure earth

Depending on your needs, you may need to ensure that Git has authenticated access and / or that Docker is logged in so that it has access to private repositories.

To authenticate Git, you may either use SSH-based authentication, or username-password-based authentication. See the [Authentication page for more information](./auth.md). In case you don't need any Git authentication, you might want to force all GitHub URLs to be transformed to `https://github.com/...` instead of `git@github.com:...`. For this, you can add an environment variable to configure this behavior:

```bash
export GIT_URL_INSTEAD_OF="https://github.com/=git@github.com:"
```

The way you configure environment variables in your CI will vary.

To log in Docker, simply run

```bash
docker login --username '<username>' --password '<password>'
```

{% hint style='info' %}
##### Note

Make sure that secrets (like `<password>` above) are not exposed in plain text. You may need to configure an environment variable with your CI vendor.
{% endhint %}

## Step 4: Run the build

This is often as simple as

```bash
earth +target-name
```

If you would like to enable pushing Docker images to registries and also running `RUN --push` commands, you might use

```bash
earth --push +target-name
```

If you need to pass secrets to the Earthly build, you might also use the `--secret` flag, mentioning the env var where the secret is kept.

```bash
earth --secret SOME_SECRET_ENV_VAR +target-name
```

For more information see the [earth command reference](../earth-command/earth-command.md).

## Complete examples

#### CircleCI

```yml
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
      - run: "sudo /bin/sh -c 'curl -s https://api.github.com/repos/vladaionescu/earthly-releases/releases/latest | grep browser_download_url | grep linux-amd64 | cut -d : -f 2- | tr -d \\x22 | wget -P /usr/local/bin/ -i - && mv /usr/local/bin/earth-linux-amd64 /usr/local/bin/earth && chmod +x /usr/local/bin/earth'"
      - run: earth --version
      - run: earth --push +build
```

#### GitHub Actions

```yml
name: CI

on:
  push:
    branches: [ master ]
  pull_request:
    branches: [ master ]

jobs:
  build:
    runs-on: ubuntu-latest
    env:
      DOCKERHUB_USERNAME: ${{ secrets.DOCKERHUB_USERNAME }}
      DOCKERHUB_TOKEN: ${{ secrets.DOCKERHUB_TOKEN }}
      GIT_URL_INSTEAD_OF: "https://github.com/=git@github.com:"
    steps:
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
    - name: Download latest earth
      run: "sudo /bin/sh -c 'curl -s https://api.github.com/repos/vladaionescu/earthly-releases/releases/latest | grep browser_download_url | grep linux-amd64 | cut -d : -f 2- | tr -d \\x22 | wget -P /usr/local/bin/ -i - && mv /usr/local/bin/earth-linux-amd64 /usr/local/bin/earth && chmod +x /usr/local/bin/earth'"
    - name: Earth version
      run: earth --version
    - name: Run build
      run: earth --push +build
```
