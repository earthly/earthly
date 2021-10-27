## Releasing instructions

### earthly
* Make sure you have access to the `earthly-technologies` organization secrets.
  ```bash
  ./earthly secrets ls /earthly-technologies
  ```
* Make sure you have uploaded your aws credentials to your user secrets.
  ```bash
  ./earthly secrets get /user/earthly-technologies/aws/credentials
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
* Update the CHANGELOG.md with the corresponding release notes and open a PR
  * Use a comparison such as https://github.com/earthly/earthly/compare/v0.3.0...main (replace the versions in the URL with the previously released version) to see which PRs will go into this release.
* Make sure that main build is green for all platforms (check build status for the latest commit on GitHub).
* Make sure nightly BuildKite builds are green for most recent build
  | Platform      | Status        |
  | ------------- | ------------- |
  | MacOS (x86)   | [![Build status](https://badge.buildkite.com/cc0627732806ab3b76cf13b02c498658b851056242ec28f62d.svg)](https://buildkite.com/earthly-technologies/earthly-mac-scheduled)
  | MacOS (M1)    | [![Build status](https://badge.buildkite.com/10a7331b2032fcc9f7f311c5218d12c1a18c317cd7fc9270ba.svg)](https://buildkite.com/earthly-technologies/earthly-m1-scheduled)
  | Windows (WSL) | [![Build status](https://badge.buildkite.com/19d9cf7fcfb3e0ee45adcb29abb5ef3cfcd994fba2d6dc148c.svg)](https://buildkite.com/earthly-technologies/earthly-windows-scheduled)
* Run
  ```bash
  ./release.sh
  ```
* Run
  ```bash
  ./earthly --push ./release+release-repo --RELEASE_TAG="$RELEASE_TAG"
  ```
  TODO: This step will be merged into the release.sh command once our staging environment is setup
* Merge branch `main` into `next`, then merge branch `next` into `main`.
* Update the version for the installation command in the following places:
<!-- vale HouseStyle.Spelling = NO -->
  * [ci-integration.md](../docs/ci-integration.md)
  * [circle-integration.md](../docs/ci-integration/guides/circle-integration.md)
  * [gh-actions-integration.md](../docs/ci-integration/guides/gh-actions-integration.md)
  * [codebuild-integration.md](../docs/ci-integration/guides/codebuild-integration.md)
  * [build-an-earthly-ci-image.md](../docs/ci-integration/build-an-earthly-ci-image.md)
<!-- vale HouseStyle.Spelling = YES -->
  * you can try doing that with:
    ```
    REGEX='\(earthly\/releases\/download\/\)v[0-9]\+\.[0-9]\+\.[0-9]\+\(\/\)'; grep -Ril './docs/' -e $REGEX | xargs -n1 sed -i 's/'$REGEX'/\1'$RELEASE_TAG'\2/g'
    ```
* Update the pinned image tags used in the following places:
<!-- vale HouseStyle.Spelling = NO -->
  * [all-in-one.md](../docs/docker-images/all-in-one.md)
  * [buildkit-standalone.md](../docs/docker-images/buildkit-standalone.md)
  * [build-an-earthly-ci-image.md](../docs/ci-integration/build-an-earthly-ci-image.md)
<!-- vale HouseStyle.Spelling = YES -->
  * you can try doing that with:
    ```shell
    REGEX='\(\searthly\/\(buildkitd\|earthly\):\)v[0-9]\+\.[0-9]\+\.[0-9]\+'; grep -Ril './docs/' -e $REGEX | xargs -n1 sed -i 's/'$REGEX'/\1'$RELEASE_TAG'/g'
    ```
* Update the Docker image documentation's tags with the new version, plus the prior two image versions under:
<!-- vale HouseStyle.Spelling = NO -->
  * [all-in-one.md](../docs/docker-images/all-in-one.md)
  * [buildkit-standalone.md](../docs/docker-images/buildkit-standalone.md)
<!-- vale HouseStyle.Spelling = YES -->
* After GitBook has processed the `main` branch, run a broken link checker over https://docs.earthly.dev. This one is fast and easy: https://www.deadlinkchecker.com/.
* Verify the [Homebrew release job](https://github.com/earthly/homebrew-earthly) has successfully run and has merged the new `release-v...` branch into `main`.
* Copy the release notes you have written before and paste them in the Earthly Community slack channel `#announcements`, together with a link to the release's GitHub page. If you have Slack markdown editing activated, you can copy the markdown version of the text.
* Ask Adam to tweet about the release.

### One-Time (clear this section when done during release)

* Add new one-time items here.

#### Performing a test release

To perform a test release to a personal repo, first:

1. fork a copy of both `earthly/earthly`, and `earthly/homebrew-earthly`
2. commit your changes you wish to release and push them to your personal repo.
3. save a copy of your GitHub token to `+secrets/user/github-token` (e.g. `earthly secrets set /user/github-token keep-it-secret`)

Then run:

  ```bash
  RELEASE_TAG=v0.5.10 GITHUB_USER=mygithubuser DOCKERHUB_USER=mydockerhubuser EARTHLY_REPO=earthly BREW_REPO=homebrew-earthly GITHUB_SECRET_PATH=+secrets/user/github-token ./release.sh
  ```

NOTE: apt and yum repos do not currently support test releases. (TODO: fix this)

#### Troubleshooting

If the release-homebrew fails with a rejected git push, you may have to delete the remote branch by running the following under the interactive debugger:

    git push "$GIT_USERNAME" --delete "release-$RELEASE_TAG"

#### Rollbacks

If you need to rollback/disable a version:

1. Go to [GitHub releases](https://github.com/earthly/earthly/releases), click on the `edit release` button, then check the `This is a prerelease` checkbox.
2. Check out the [earthly/homebrew-earthly](https://github.com/earthly/homebrew-earthly) repo, and run:

```bash
git checkout main
git revert --no-commit 123abc..HEAD # where `123abc` is the sha1 commit to roll back to
git commit # enter a message saying you are rolling back
git push
```

3. TODO need to create targets for apt and yum Earthfiles to perform rollbacks

### dind

Docker-in-Docker (dind) images change less frequently than earthly, but take a long time to build.
These images can be rebuilt by running:

  ```bash
  ./earthly --push ./release+release-dind --RELEASE_TAG="$RELEASE_TAG"
  ```

### VS Code syntax highlighting

* First set the version to publish:
  ```bash
  export VSCODE_RELEASE_TAG="v..."
  ```
  (You can [see what is already published](https://marketplace.visualstudio.com/items?itemName=earthly.earthfile-syntax-highlighting))
* Make sure that the version has release notes already in the [README](../contrib/earthfile-syntax-highlighting/README.md)
* Then publish it:
  ```bash
  ./earthly \
    --push \
    ./release+release-vscode-syntax-highlighting \
    --VSCODE_RELEASE_TAG="$VSCODE_RELEASE_TAG"
  ```
* Finally, tag git for future reference
  ```bash
  git tag "vscode-syntax-highlighting-$VSCODE_RELEASE_TAG"
  git push origin "vscode-syntax-highlighting-$VSCODE_RELEASE_TAG"
  ```

(If token has expired, Vlad can regenerate one following [this guide](https://code.visualstudio.com/api/working-with-extensions/publishing-extension#get-a-personal-access-token) and then setting it using `./earthly secrets set /earthly-technologies/vsce/token '...'`)
