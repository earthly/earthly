# 🌎 Earthly - a build system for inhabitants of planet earth

*Parallel, reproducible, consistent and portable builds for the same era as your code*

## Why

TODO: Simplify this by moving details into launch blog post. Keep only the essence here.

Why does the world need a new build system?

We live in an era of containers, CI/CD, automation, rich set of programming languages, varying code structures (mono/poly-repos) and open-source collaboration. None of the build systems out there serve well these new trends:

* They don't take advantage of container-based isolation to make the builds portable
* Builds are not easily reproducible - they often depend on already installed dependencies
* They don't provide a way to import open-source recipes
* Many are programming language-specific, making them unattractive to be used as a one-stop-shop
* Parallelization ranges from difficult to almost impossible
* Importing and reusability are primitive or difficult to use across repos
* Caching is difficult to master, making it impossible to scale the build system in a mono-repo
* They are difficult to use, leading to a build system guru situation (only one person knows how the build works)
* Or a combination of the above.

Of the popular choices out there, the options that come close are [Bazel](https://bazel.build/) and [Dockerfiles](https://docs.docker.com/engine/reference/builder/). Both are excellent, yet there are challenges: Bazel is difficult to adopt because it requires an all-in approach (are you ready to completely rewrite the build.gradle's in your org into [Bazel BUILD files](https://docs.bazel.build/versions/master/tutorial/java.html)?). Dockerfiles are great, but they only output images. You can subsequently use docker run commands and mounted volumes to output other kinds of artifacts - but that requires that you now wrap your docker builds into Makefiles or some other build system.

Earthly accepts the reality that for some languages, the best build system is provided by the community of that language (like gradle, webpack, sbt etc), yet it adds the Dockerfile-like caching on top, plus more flexibility in defining hierarchies of dependencies for caching and reusability.

### Benefits

Here's what you get with Earthly:

* Consistent environment between developers and with CI
* Programming language agnostic
* No need to install project-specific dependencies - the build is self-contained
* First-party support for Docker
* Syntax easy to understand even with no previous experience
* Efficient use of caching, that is easy to understand by anyone
* Simple parallelism, without the gnarly race conditions
* Like a docker build but it can also yield classical artifacts (packages, binaries etc)
  owned by the user, not by root
* Mono-repo friendly (reference targets within subdirs, with no strings attached)
* Multi-repo friendly (reference targets in other repos just as easily)
* In a complex setup, ability to trigger rebuild of only affected targets and tests
* An import system that is easy to use and has no hidden implications

## Dive in

... TODO
