# Releasing instructions

### earth

* Make sure you have a GITHUB_TOKEN set. If you don't have a GITHUB_TOKEN, generate one [here](https://github.com/settings/tokens) with scope `repo`.
  ```bash
  export GITHUB_TOKEN="..."
  ```
* Choose the next [release tag](https://github.com/earthly/earthly/releases).
  ```bash
  export RELEASE_TAG="v..."
  ```
* Export github username.
  ```bash
  export GIT_USERNAME="..."
  ```
* Make sure you are on master
  ```bash
  git checkout master && git pull
  ```
* Run
  ```bash
  earth \
    --build-arg RELEASE_TAG \
    --secret GITHUB_TOKEN \
    --push -P ./release+release
  ```
* Run
  ```bash
  earth \
    --build-arg RELEASE_TAG \
    --build-arg GIT_USERNAME \
    --build-arg GIT_NAME="$(git config user.name)" \
    --build-arg GIT_EMAIL="$(git config user.email)" \
    --secret GITHUB_TOKEN \
    --push ./release+release-homebrew
  ```
* Merge branch `next` into `master`.

### VS Code syntax highlighting

TODO

### Using a fork of buildkit

* Build a buildkit image with

```bash
DOCKER_BUILDKIT=1 docker build -t earthly/buildkit:fix-ssh-auth-sock --target buildkit-buildkitd-linux .
```

* Push it to docker hub

```bash
docker push earthly/buildkit:fix-ssh-auth-sock
```

* Use it in our own build, but pin it to a specific sha256, just in case.
