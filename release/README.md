# Releasing instructions

### earth
* Test on Mac using `earth --artifact github.com/earthly/earthly+for-darwin/earth ./earth && docker pull earthly/buildkitd:master && ./earth -P github.com/earthly/earthly+test-all`
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
* Once pushed, download and install the released earth locally for both linux and mac and test it using `earth -P github.com/earthly/earthly+test-all`
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
* Update the version for the installation command in the following places:
  * [ci-integration.md](../docs/guides/ci-integration.md)
  * [circle-integration.md](../docs/examples/circle-integration.md)
  * [gh-actions-integration.md](../docs/examples/gh-actions-integration.md)
* Go to the [releases page](https://github.com/earthly/earthly/releases) and edit the latest release to add release notes. Use a comparison such as https://github.com/earthly/earthly/compare/v0.3.0...v0.3.1 (replace the right versions in the URL) to see which PRs went into this release.
* Keep an eye on the Homebrew PR and once merged, test it out on a mac using
  ```bash
  brew upgrade earthly
  earth -P github.com/earthly/earthly+test-all
  ```

### VS Code syntax highlighting

First set the version to publish:

```bash
export VSCODE_RELEASE_TAG=v0.0.3
```

(you can see what is already published at https://marketplace.visualstudio.com/items?itemName=earthly.earthfile-syntax-highlighting )

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
