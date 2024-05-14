# Target, artifact and function referencing

This page describes the different types of references used in Earthly:

* Target references: `<project-ref>+my-target`
* Artifact references: `<project-ref>+my-target/my-artifact.bin`
* Image references (same as target references)
* Function references: `<project-ref>+MY_FUNCTION`
* Project references (the prefix of the above references): `github.com/foo/bar`, `./my/local/path`

## Target reference

Target references point to an Earthly target. They have the general form

`<project-ref>+<target>`

Target references distinguish themselves from function references (see below) by having a name in all-lower-case, kebab-case (e.g. `+my-target`).

Here are some examples:

* `+build`
* `./js+deps`
* `github.com/earthly/earthly:v0.8.10+earthly`

## Artifact reference

Artifact references are similar to target references, except that they have an artifact path at the end. It has the following form

`<target-ref>/<artifact-path>`

Here are some examples:

* `+build/my-artifact`
* `+build/some/artifact/deep/in/a/dir`
* `./js+build/dist`
* `github.com/earthly/earthly:v0.8.10+earthly/earthly`

## Image reference

Because there can only be one image per target, image references have the exact same format as target references.

The only difference is the context where they are used. For example, a `FROM` command takes an image reference. While a `BUILD` command takes a target reference.

## Function reference

Function references point to a function in an Earthfile. They have the general form

`<project-ref>+<function>`

Function references distinguish themselves from target references by having a name in all-caps, snake-case (e.g. `+MY_FUNCTION`).

Here are some examples:

* `+COMPILE`
* `./js+NPM_INSTALL`
* `github.com/earthly/earthly:v0.8.10+DOWNLOAD_DIND`

For more information on functions, see the [functions guide](./functions.md).

## Project references

<img src="img/ref-infographic-v2.png" alt="Target and artifact reference syntax" title="Reference targets using +" width="800px" />

Project references appear in target, artifact and function references. They point to the Earthfile containing the respective target, artifact or function. Below are the different types of project references available in Earthly.

### Local, internal

The simplest form, is where a target, function or artifact is referenced from the same Earthfile. In this case, the project reference is simply **the empty string**. Here are some examples of this type of project reference being used in various other references:

| Project ref | Target ref | Artifact ref | Function ref |
|----|----|----|----|
| (**empty string**) | `+<target-name>` | `+<target-name>/<artifact-path>` | `+<function-name>` |
| (**empty string**) | `+build` | `+build/out.bin` | `+COMPILE` |

In this form, Earthly will look for the target within the same Earthfile. We call this type of referencing local, internal. Local, because it comes from the same system, and internal, because it is within the same Earthfile.

### Local, external

Another form, is where a target, function or artifact is referenced from a different directory. In this form, the path to that directory is specified before `+`. It must always start with either `./`, `../` or `/`, on any operating system (including Windows). Example:

| Project ref | Target ref | Artifact ref | Function ref |
|----|----|----|----|
| `./path/to/another/dir` | `./path/to/another/dir+<target-name>` | `./path/to/another/dir+<target-name>/<artifact-path>` | `./path/to/another/dir+<function-name>` |
| `./js` | `./js+build` | `./js+build/out.bin` | `./js+COMPILE` |

It is recommended that relative paths are used, for portability reasons: the working directory checked out by different users will be different, making absolute paths infeasible in most cases.

### Remote

Another form of a project reference is the remote form. In this form, the recipe and the build context are imported from a remote location. It has the following form:

| Project ref | Target ref | Artifact ref | Function ref |
|----|----|----|----|
| `<vendor>/<namespace>/<project>/path/in/project[:some-tag]` | `<vendor>/<namespace>/<project>/path/in/project[:some-tag]+<target-name>` | `<vendor>/<namespace>/<project>/path/in/project[:some-tag]+<target-name>/<artifact-path>` | `<vendor>/<namespace>/<project>/path/in/project[:some-tag]+<function-name>` |
| `github.com/earthly/earthly/buildkitd` | `github.com/earthly/earthly/buildkitd+build` | `github.com/earthly/earthly/buildkitd+build/out.bin` | `github.com/earthly/earthly/buildkitd+COMPILE` |
| `github.com/earthly/earthly:v0.8.10` | `github.com/earthly/earthly:v0.8.10+build` | `github.com/earthly/earthly:v0.8.10+build/out.bin` | `github.com/earthly/earthly:v0.8.10+COMPILE` |

### Import reference

Finally, the last form of project referencing is an import reference. Import references may only exist after an `IMPORT` command, which helps resolve the reference to a full project reference of the types above.

| Import command | Project ref | Target ref | Artifact ref | Function ref |
|----|----|----|----|----|
| `IMPORT <full-project-ref> AS <import-alias>` | `<import-alias>` | `<import-alias>+<target-name>` | `<import-alias>+<target-name>/<artifact-path>` | `<import-alias>+<function-name>` |
| `IMPORT github.com/earthly/earthly/buildkitd` | `buildkitd` | `buildkitd+build` | `buildkitd+build/out.bin` | `buildkitd+COMPILE` |
| `IMPORT github.com/earthly/earthly:v0.8.10` | `earthly` | `earthly+build` | `earthly+build/out.bin` | `earthly+COMPILE` |

Here is an example in an Earthfile:

```Dockerfile
IMPORT github.com/earthly/earthly/buildkitd

...

BUILD buildkitd+buildkitd
```

## Implicit Base Target Reference

All Earthfiles start with a base recipe. This is the only recipe which does not have an explicit target name - the name is always implied to be `base`. All other target implicitly inherit from `base`. You can imagine that all recipes start with an implicit `FROM +base`

```
# base recipe
FROM golang:1.15-alpine3.13
WORKDIR /go-example

build:
    # implicit FROM +base
    RUN echo "Hello World"
```

## Canonical form

Most references have a canonical form. It is essentially the remote form of the same target, with repository and tag inferred. The canonical form can be useful as a universal identifier for a target.

For example, depending on where the files are stored, the `+build` target could have the canonical form `github.com/some-user/some-project/some/deep/dir:master+build`, where `github.com/some-user/some-project` was inferred as the Git location, based on the Git remote called `origin`, and `/some/deep/dir` was inferred as the sub-directory where `+build` exists within that repository. The Earthly tag is inferred using the following algorithm:

* If the current HEAD has at least one Git tag, then use the first Git tag listed by Git, otherwise
* If the repository is not in detached HEAD mode, use the current branch, otherwise
* Use the current Git hash.

If no Git context is detected by Earthly, then the target does not have a canonical form.
