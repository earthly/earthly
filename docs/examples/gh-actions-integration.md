
# GitHub Actions integration

Here is an example a GitHub Actions build, where we build the Earthly target `+build`.

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
      GIT_URL_INSTEAD_OF: "https://github.com/=git@github.com:"
      FORCE_COLOR: 1
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
      run: "sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/download/v0.3.13/earth-linux-amd64 -O /usr/local/bin/earth && chmod +x /usr/local/bin/earth'"
    - name: Earth version
      run: earth --version
    - name: Run build
      run: earth --push +build
```

For a complete guide on CI integration see the [CI integration guide](../guides/ci-integration.md).
