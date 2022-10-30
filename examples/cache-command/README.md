# Using the CACHE Command

This guide shows how the `CACHE` command can be used in Earthly to optimize Targets that would normally perform better with incremental changes. Such as when downloading 3rd-party dependencies.

In such cases, we can optimize the cache so that results from a previous build can be reused when the files a cache depends on change.

Importantly, the results of the cache are persisted, and can be used in subsequent Targets, or available in a resulting Docker image.

## Use Cases

Examples of when the `CACHE` command becomes useful include the following.

Java/Maven projects using `SNAPSHOT` builds:
* In such projects, dependencies may change frequently, but we can avoid downloading *all* of them each time in our Earthly Target

Node.JS/NPM projects with vendored `node_modules` using version ranges:
  * We can avoid downloading all dependencies to our `node_modules` directory and instead only download those which have changed

Elixir's Dialyzer: 
  * Dialyzer is compiled quickly with incremental code changes and we can take advantage of that in our Earthly Target

In these examples, adding the `CACHE` command allows an Earthly Target to perform more closely to how you would expect when running the commands natively in your local dev environment. 

## Usage

The `CACHE` command takes a singe argument, which is a directory containing the files you want to cache for future resuse.

```Earthfile
CACHE <directory>
```

The `<directory>` is also persisted into the resulting image at the end of the Target's execution. Meaning that you can access the directory in a subsequent Target by referencing it with `FROM`, or if you run `SAVE IMAGE`, the directory will be accessible in the resulting Docker image.

## Simple Example
A classic example for uing `CACHE` is within a Target which downloads new dependenices frequently. 

Here's an example Target that uses `CACHE` in an NPM build:
```Earthfile
VERSION --use-cache-command 0.6

deps:
  CACHE ./node_modules
  COPY package*.json .
  RUN npm install
```

In the above example, our `node_modules` directory is cached within the Earthly Target. When our `package.json` changes, Earthly resuses the contents of `./node_modules` from the previous execution so that only the dependencies that have changed since last time will be downloaded to `node_modules`.

Without the `CACHE` command in the above example, if there's any change to the `package.json` file, *all* of the dependencies would need to be downloaded each time.

## Tradeoffs

Although this technique can be used to speed up execution, it may have the following drawbacks that you should consider.

### Cache Bloat
Old files no longer used in the cache directory are not automatically cleaned up. 

For example, if a dependency was cached previously but no longer used (perhaps updated to a new version), the old file would permentantly exist in the cache. Doing this enough times may result in a cache directory that becomes quite large. 

A simple solution is to delete the cache, either by running the Target with a [`--no-cache`](https://docs.earthly.dev/docs/guides/advanced-local-caching#option-2-mount-based-caching-advanced) flag, or by running [`earthly prune`](https://docs.earthly.dev/docs/earthly-command#earthly-prune).


### Reduced Repeatability

Since the target uses a local cache, it can perform differently on different machines. For example, it may work differently on your local machine than it might on your CI or a colleague's machine.

Consider a situation where a dependency has been removed from a registry such as NPM. Although unlikely, this could result in a build that works on a machine that already has it cached, but on another machine the build may consistently fail.

A worse situation may occur if a build tool doesn't ignore a cached dependency that was removed from a manifest file (e.g. `pom.xml`, `package.json`, `build.gradle`). In such a case, a build  could behave differently (or break) on another machine if it unintentionally uses a cached dependency that should be removed.

## More Detailed Examples

Some fully functional example projects are contained within subdirectories to demonstrate the CACHE command in more detail. 

Take a look at the Earthfiles in these examples to see how they work. 

Feel free to reach out via Github Issue or on our [Community Slack channel](https://earthly.dev/slack) if you have any questions or ideas. Thanks for using Earthly!
