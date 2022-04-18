
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
      run: echo "${{ secrets.DOCKERHUB_TOKEN }}" | docker login -u ${{ secrets.DOCKERHUB_USERNAME }} --password-stdin
    - name: Setup Earthly
      uses: earthly/actions-setup@v1
      with:
        version: "v0.6.10"
    - name: Run build
      run: earthly --ci --push +build
```

For a complete guide on CI integration see the [CI integration guide](../overview.md).
