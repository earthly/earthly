# Releasing instructions

### earth
* Make sure you have access to the `earthly-technologies` organization secrets.
  ```bash
  ./earth secrets ls /earthly-technologies
  ```
* Choose the next [release tag](https://github.com/earthly/earthly/releases).
  ```bash
  export RELEASE_TAG="v..."
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
    --push -P ./release+release
  ```
* Run
  ```bash
  ./earth \
    --build-arg RELEASE_TAG \
    --push ./release+release-homebrew
  ```
* Merge branch `main` into `next`, then merge branch `next` into `main`.
* Update the version for the installation command in the following places:
  * [ci-integration.md](../docs/guides/ci-integration.md)
  * [circle-integration.md](../docs/examples/circle-integration.md)
  * [gh-actions-integration.md](../docs/examples/gh-actions-integration.md)
* Go to the [releases page](https://github.com/earthly/earthly/releases) and edit the latest release to add release notes. Use a comparison such as https://github.com/earthly/earthly/compare/v0.3.0...v0.3.1 (replace the right versions in the URL) to see which PRs went into this release.
* Once everything looks good, uncheck the "pre-release" checkbox on the GitHub release page. This will make this release the "latest" when people install Earthly with the one-liner.
* Post link to release & homebrew PR in the `#release` channel on internal Slack.
* Copy the release notes you have written before and paste them in the Earthly Community slack channel `#announcements`, together with a link to the release's GitHub page. If you have Slack markdown editing activated, you can copy the markdown version of the text.
* Ask Adam to tweet about the release.

### VS Code syntax highlighting

* First set the version to publish:
  ```bash
  export VSCODE_RELEASE_TAG=...
  ```
  (You can see what is already published [here](https://marketplace.visualstudio.com/items?itemName=earthly.earthfile-syntax-highlighting))
* Make sure that the version has release notes already in the [README](../contrib/earthfile-syntax-highlighting/README.md)
* Then publish it:
  ```bash
  ./earth \
    --build-arg VSCODE_RELEASE_TAG \
    --push \
    ./release+release-vscode-syntax-highlighting
  ```
* Finally, tag git for future reference
  ```bash
  git tag "vscode-syntax-highlighting-$VSCODE_RELEASE_TAG"
  git push origin "vscode-syntax-highlighting-$VSCODE_RELEASE_TAG"
  ```

* (If token has expired, Vlad can regenerate one following [this guide](https://code.visualstudio.com/api/working-with-extensions/publishing-extension#get-a-personal-access-token)). The set it using `./earth secrets set /earthly-technologies/vsce/token '...'`
