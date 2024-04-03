## Releasing instructions

### earthly
* Make sure you have access to the `earthly-technologies` organization secrets.
  ```bash
  ./earthly secrets --org earthly-technologies --project core ls
  ```
* Choose the next [release tag](https://github.com/earthly/earthly/releases).
  ```bash
  export RELEASE_TAG="v..."
  ```
* Is it a pre-release?
  ```bash
  export PRERELEASE="true-or-false"
  ```
* Make sure you are on main
  ```bash
  git checkout main && git pull
  ```
* Update the CHANGELOG.md with the corresponding release notes and open a PR
  * Use a comparison such as https://github.com/earthly/earthly/compare/v0.8.7...main (replace the versions in the URL with the previously released version) or a tool such as `gitk` (aka `git-gui`) to see which PRs will go into this release.
* Make sure that main build is green for all platforms (check build status for the latest commit on GitHub).
* Make sure the following build status are green:
  | Platform      | Status        |
  | ------------- | ------------- |
  | MacOS (x86)   | [![Build status](https://badge.buildkite.com/cc0627732806ab3b76cf13b02c498658b851056242ec28f62d.svg)](https://buildkite.com/earthly-technologies/earthly-mac-scheduled)
  | MacOS (M1)    | [![Build status](https://badge.buildkite.com/10a7331b2032fcc9f7f311c5218d12c1a18c317cd7fc9270ba.svg)](https://buildkite.com/earthly-technologies/earthly-m1-scheduled)
* Run
  ```bash
  cd release
  env -i HOME="$HOME" PATH="$PATH" SSH_AUTH_SOCK="$SSH_AUTH_SOCK" RELEASE_TAG="$RELEASE_TAG" USER="$USER" PRERELEASE="$PRERELEASE" ./release.sh
  ```
* Wait for the [Merge main to docs-0.8 on New Earthly Release](../.github/workflows/release-merge-docs.yml) workflow to complete; this workflow automatically merges `main` into `docs-0.8`. You can watch for it here: [![Merge main to docs-0.8 on New Earthly Release](https://github.com/earthly/earthly/actions/workflows/release-merge-docs.yml/badge.svg)](https://github.com/earthly/earthly/actions/workflows/release-merge-docs.yml)
In case the workflow fails the manual process is:
  ```shell
    git checkout docs-0.8 && git pull && git merge main && git push
    ```
* Updating the Earthly version in our docs:
  [Renovate](https://www.mend.io/renovate/) will open a PR targeting `docs-0.8` branch to update all docs as soon as a new release is available in this repo which you should then review & merge (An example PR can be found [here](https://github.com/earthly/earthly/pull/3285/files)).
* Merge `docs-0.8` into `main`.
  ```shell
    git checkout main && git merge docs-0.8 && git push
    ```
  * Note: If you don't have permissions to push directly to `main` branch, do the following:
    * `git checkout -b soon-to-be-main && git push origin soon-to-be-main`
    * Open a PR against the new branch and get it approved; IMPORTANT: don't squash-merge via github
    * Once all (required) checks pass, try pushing the branch again:
    `git checkout main && git push`

<!-- vale HouseStyle.Spelling = YES -->
* Wait for the [Check Docs for Broken Links](../.github/workflows/docs-checks-links.yml) workflow to complete; this workflow validates https://docs.earthly.dev does not contain any broken links. You can watch for it here: [![Check Docs for Broken Links](https://github.com/earthly/earthly/actions/workflows/docs-checks-links.yml/badge.svg?event=push)](https://github.com/earthly/earthly/actions/workflows/docs-checks-links.yml)
* Verify the [Homebrew release job](https://github.com/earthly/homebrew-earthly) has successfully run and has merged the new `release-v...` branch into `main`.
* Copy the release notes you have written before and paste them in the Earthly Community slack channel `#announcements`, together with a link to the release's GitHub page. If you have Slack markdown editing activated, you can copy the markdown version of the text.

### One-Time (clear this section when done during release)

* Add new one-time items here.

#### Performing a test release

To perform a test release to a personal repo, first:

1. fork a copy of both `earthly/earthly`, and `earthly/homebrew-earthly`
2. commit your changes you wish to release and push them to your personal repo.
3. save a copy of your GitHub token to `user/github-token` (e.g. `earthly secrets set /user/github-token keep-it-secret`)

Then run:

  ```bash
  RELEASE_TAG=v0.5.10 GITHUB_USER=mygithubuser DOCKERHUB_USER=mydockerhubuser EARTHLY_REPO=earthly BREW_REPO=homebrew-earthly GITHUB_SECRET_PATH=user/github-token ./release.sh
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
3. Mark the release title in [CHANGELOG.log](../CHANGELOG.md) as `(aborted release/not recommended)`, e.g.:
`## v0.7.18 - 2023-09-18 (aborted release/not recommended)`
4. TODO need to create targets for apt and yum Earthfiles to perform rollbacks

### dind

Docker-in-Docker (dind) images change less frequently than earthly, but take a long time to build.
earthly/dind images and their releases are maintained in [project repo](https://github.com/earthly/dind).

### Syntax Highlighting Releases

We currently have syntax highlighting for the following:
1. [vscode + github](https://github.com/earthly/earthfile-grammar)
1. [intellij](https://github.com/earthly/earthly-intellij-plugin) (py, go, java)
1. [vim](https://github.com/earthly/earthly.vim)
1. [sublime](https://github.com/earthly/sublimetext-earthly-syntax)
1. [emacs](https://github.com/earthly/earthly-emacs)


#### VSCode + Github

Release instructions can be found in the [project repo](https://github.com/earthly/earthfile-grammar#how-to-release).

#### Intellij

Intellij pulls its syntax highlighting from the [same repo used by VSCODE + Github](https://github.com/earthly/earthfile-grammar) and so should be released after to keep up to date.

1. Go to the [repo](https://github.com/earthly/earthfile-grammar)
1. Make relevant changes to the branches + test in this order:
    1. py
    1. go
    1. main
1. Sign + release the changes from each branch in this order:
    1. py
    1. go
    1. main

    Follow the instructions on how to sign and release as written in the [README](https://github.com/earthly/earthly-intellij-plugin#signing-requires-earthly-technologies-org-membership)

#### Vim

1. Go to the [repo](https://github.com/earthly/earthly.vim)
1. Make relevant updates and test
1. Once merged to main it will be released

#### Sublime Text

1. Go to the [repo](https://github.com/earthly/sublimetext-earthly-syntax)
1. Make relevant updates and test
1. Once merged to main it will be released

#### Emacs

1. Go to the [repo](https://github.com/earthly/earthly-emacs)
1. Make relevant updates and test
1. Once merged to main it will be released
