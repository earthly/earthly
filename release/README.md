# Releasing instructions

### earth
* Ask Vlad to test the latest version on his Mac
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

First set the version to publish:

```bash
export VSCODE_RELEASE_TAG=v0.0.3
```

Ask Vlad for a token

```bash
VSCE_TOKEN=.....
```
(vlad can generate one following this guide: https://code.visualstudio.com/api/working-with-extensions/publishing-extension#get-a-personal-access-token )

Then publish it:
```bash
earth \
  --build-arg VSCODE_RELEASE_TAG \
  --secret VSCE_TOKEN \
  --push \
  +release-vscode-syntax-highlighting
```

### Using a fork of buildkit

We maintain a fork of buildkit under https://github.com/earthly/buildkit which has an earthly-master branch which contains our patches.

Github actions performs a build of the buildkit docker image; under the build output expand the section titled `push buildkit docker image`
and look for the last line, which is similar to:

```
earthly-master: digest: sha256:40303c69b24c23c63c417efac1b6641e53ecc526598edf101a965c3dc54dddc3 size: 1158
```

Then update the [../buildkitd/Earthfile](buildkit earthfile) FROM entry with the updated sha256 string:

```
FROM earthly/buildkit:earthly-master@sha256:40303c69b24c23c63c417efac1b6641e53ecc526598edf101a965c3dc54dddc3
```
