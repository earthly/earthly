# CACHE Example

This example demonstrates how you can use the `CACHE` command in Earthly to optimize builds that benefit from incremental changes to files such as external dependencies, or source code with programming languages that support incremental compilation. 

In such builds, we can optimize the cache so that results from a previous build can be reused when the files a cache depends on change.

## Use Cases

Some example use cases are:

* Java/Maven project using `SNAPSHOT` builds that change frequently
  * In such projects, we can reuse other dependencies which haven't changed such that our Earthly build only downloads those which are necessary
* Node.JS/NPM projects with vendored `node_modules` with version ranges
  * Similar to the Maven example, we can reduce the amount of time spent building these projects by reusing the dependencies in `node_modules` which have not changed.
* Elixir's Dialyzer 
  * Dialyzer is compiled quickly with incremental code changes and we can take advantage of that in our Earthly builds

## The CACHE Command
All of these examples are made more efficient using the `CACHE` command available in your Earthfile.

The `CACHE` command takes a path to a directory containing the 

## Examples

Some examples are contained here in subdirectories to demonstrate how the CACHE command works. Explore the Earthfiles in those subdirectories to learn more, and feel free to reach out to use via Github Issue or our Slack channel if you have any questions or ideas.
