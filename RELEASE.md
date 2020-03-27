# Release checklist

* Make sure you have a GITHUB_TOKEN set.
  ```bash
  export GITHUB_TOKEN="..."
  ```
* Choose the next [release tag](https://github.com/vladaionescu/earthly/releases).
  ```bash
  export RELEASE_TAG="v..."
  ```
* Run
  ```bash
  earth --build-arg RELEASE_TAG --secret GITHUB_TOKEN --push -P +release
  ```
* Update [README.md](./README.md) installation instructions to use the newly released version.
