# Releasing instructions

### earthly
* Make sure you have access to the `earthly-technologies` organization secrets.
  ```bash
  ./earthly secrets ls /earthly-technologies
  ```
* Make sure you have uploaded your aws credentials to your user secrets.
  ```bash
  ./earthly secrets ls /user/earthly-technologies/aws/credentials
  ```
* Choose the next [release tag](https://github.com/earthly/earthly/releases).
  ```bash
  export RELEASE_TAG="v..."
  ```
* Make sure you are on main
  ```bash
  git checkout main && git pull
  ```
* If you are releasing a new major or minor version (`X.0.0` or `X.X.0`), then you need to update the license text. Under `licenses/BSL`
  * Update the version under `Licensed Work` to point to the new version (without the `.0` suffix)
  * Update the year (`(c) 2021`) to point to the current year
  * Update the `Change Date` to the first April 1st after a three year period. So if today is Jan 2077, then update to 2080-04-01. If it's June 2077, then update to 2081-04-01.
  * Commit this to the `main` branch before continuing.
* Make sure that main build is green for all platforms (check build status for the latest commit on GitHub).
* Run
  ```bash
  ./earthly reset
  ```
* Run
  ```bash
  ./earthly --build-arg RELEASE_TAG --push -P ./release+release
  ```
* Go to the [releases page](https://github.com/earthly/earthly/releases) and edit the latest release to add release notes. Use a comparison such as https://github.com/earthly/earthly/compare/v0.3.0...v0.3.1 (replace the right versions in the URL) to see which PRs went into this release.
* Once everything looks good, uncheck the "pre-release" checkbox on the GitHub release page. This will make this release the "latest" when people install Earthly with the one-liner. Important: You **have** to do this before the next step.
* Run
  ```bash
  ./earthly --build-arg RELEASE_TAG --push ./release+release-homebrew
  ```
* Important: Subscribe to the PR that griswoldthecat created in homebrew-core, so that you can address any review comments that may come up.
* Run
  ```bash
  ./earthly --build-arg RELEASE_TAG --push ./release+release-repo
  ```
* Merge branch `main` into `next`, then merge branch `next` into `main`.
* Update the version for the installation command in the following places:
  * [ci-integration.md](../docs/ci-integration.md)
  * [circle-integration.md](../docs/ci-examples/circle-integration.md)
  * [gh-actions-integration.md](../docs/ci-examples/gh-actions-integration.md)
  * [codebuild-integration.md](../docs/ci-examples/codebuild-integration.md)
  * you can try doing that with:
    ```
    find . -name '*integration.md' | xargs -n1 sed -i 's/\(earthly\/releases\/download\/\)v[0-9]\+\.[0-9]\+\.[0-9]\+\(\/\)/\1'$RELEASE_TAG'\2/g'
    ```
* After gitbook has processed the `main` branch, run a broken link checker over https://docs.earthly.dev. This one is fast and easy: https://www.deadlinkchecker.com/.
* Copy the release notes you have written before and paste them in the Earthly Community slack channel `#announcements`, together with a link to the release's GitHub page. If you have Slack markdown editing activated, you can copy the markdown version of the text.
* Ask Adam to tweet about the release.

#### Troubleshooting

If the release-homebrew fails with a rejected git push, you may have to delete the remote branch by running the following under the interactive debugger:

    git push "$GIT_USERNAME" --delete "release-$RELEASE_TAG"

### VS Code syntax highlighting

* First set the version to publish:
  ```bash
  export VSCODE_RELEASE_TAG="v..."
  ```
  (You can see what is already published [here](https://marketplace.visualstudio.com/items?itemName=earthly.earthfile-syntax-highlighting))
* Make sure that the version has release notes already in the [README](../contrib/earthfile-syntax-highlighting/README.md)
* Then publish it:
  ```bash
  ./earthly \
    --build-arg VSCODE_RELEASE_TAG \
    --push \
    ./release+release-vscode-syntax-highlighting
  ```
* Finally, tag git for future reference
  ```bash
  git tag "vscode-syntax-highlighting-$VSCODE_RELEASE_TAG"
  git push origin "vscode-syntax-highlighting-$VSCODE_RELEASE_TAG"
  ```

(If token has expired, Vlad can regenerate one following [this guide](https://code.visualstudio.com/api/working-with-extensions/publishing-extension#get-a-personal-access-token) and then setting it using `./earthly secrets set /earthly-technologies/vsce/token '...'`)
