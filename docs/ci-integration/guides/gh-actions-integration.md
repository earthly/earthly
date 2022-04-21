
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
      FORCE_COLOR: 1
    steps:
    - uses: actions/checkout@v2
      with:
          # By default, only the latest commit is checked out.
          # Earthly uses the branch name for tagging,
          # so we need to explicitly check out the branch.
          ref: ${{ github.head_ref || github.ref_name }}
    - name: Docker Login
      run: docker login --username "$DOCKERHUB_USERNAME" --password "$DOCKERHUB_TOKEN"
    - name: Download latest earthly
      run: "sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/download/v0.6.14/earthly-linux-amd64 -O /usr/local/bin/earthly && chmod +x /usr/local/bin/earthly'"
    - name: Earthly version
      run: earthly --version
    - name: Run build
      run: earthly --ci --push +build
```

For a complete guide on CI integration see the [CI integration guide](../overview.md).
