
# GitHub Actions integration

Here is an example of a GitHub Actions build that uses the [earthly/actions-setup](https://github.com/earthly/actions-setup).

This example assumes an [Earthfile](../../earthfile/earthfile.md) exists with a `+build` target:

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
    - uses: earthly/actions-setup@v1
      with:
        version: v0.7.8
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
    - name: Earthly version
      run: earthly --version
    - name: Run build
      run: earthly --push +build
```

Alternatively, you can skip using the `earthly/actions-setup` job and include
a step to download earthly instead:

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
    - uses: actions/checkout@v3
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
    - name: Download latest earthly
      run: "sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/download/v0.7.22/earthly-linux-amd64 -O /usr/local/bin/earthly && chmod +x /usr/local/bin/earthly'"
    - name: Earthly version
      run: earthly --version
    - name: Run build
      run: earthly --ci --push +build
```

For a complete guide on CI integration see the [CI integration guide](../overview.md).

{% hint style='danger' %}
## actions/checkout ref argument

The example deliberately does not use the [`ref`](https://github.com/actions/checkout#checkout-a-different-branch) `actions/checkout@v3` option,
as it can lead to inconsistent builds where a user chooses to re-run an older commit which is no longer at the head of the branch.
{% endhint %}
