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
* Make sure you are on main
  ```bash
  git checkout main && git pull
  ```
* Make sure that main build is green for all platforms (check build status for the latest commit on GitHub).
* Run
  ```bash
  ./earth reset
  ```
* Run
  ```bash
  ./earth \
    --build-arg RELEASE_TAG \
    --secret GITHUB_TOKEN \
    --push -P ./release+release
  ```
* Run
  ```bash
  ./earth \
    --build-arg RELEASE_TAG \
    --build-arg GIT_USERNAME \
    --build-arg GIT_NAME="$(git config user.name)" \
    --build-arg GIT_EMAIL="$(git config user.email)" \
    --secret GITHUB_TOKEN \
    --push ./release+release-homebrew
  ```
* Merge branch `next` into `main`.
* Update the version for the installation command in the following places:
  * [ci-integration.md](../docs/guides/ci-integration.md)
  * [circle-integration.md](../docs/examples/circle-integration.md)
  * [gh-actions-integration.md](../docs/examples/gh-actions-integration.md)
* Go to the [releases page](https://github.com/earthly/earthly/releases) and edit the latest release to add release notes. Use a comparison such as https://github.com/earthly/earthly/compare/v0.3.0...v0.3.1 (replace the right versions in the URL) to see which PRs went into this release.
* Post link to release & homebrew PR in the `#release` channel on internal Slack.
* Ask Adam to tweet about the release.

### VS Code syntax highlighting

* First set the version to publish:
  ```bash
  export VSCODE_RELEASE_TAG=...
  ```
  (You can see what is already published [here](https://marketplace.visualstudio.com/items?itemName=earthly.earthfile-syntax-highlighting))
* Make sure that the version has release notes already in the [README](../contrib/vscode-syntax-highlighting/README.md)
* Ask Vlad for a token
  ```bash
  VSCE_TOKEN=.....
  ```
  (Vlad can generate one following [this guide](https://code.visualstudio.com/api/working-with-extensions/publishing-extension#get-a-personal-access-token))

* Then publish it:
  ```bash
  ./earth \
    --build-arg VSCODE_RELEASE_TAG \
    --secret VSCE_TOKEN \
    --push \
    ./release+release-vscode-syntax-highlighting
  ```
* Finally, tag git for future reference
  ```bash
  git tag "vscode-syntax-highlighting-$VSCODE_RELEASE_TAG"
  git push origin "vscode-syntax-highlighting-$VSCODE_RELEASE_TAG"
  ```
