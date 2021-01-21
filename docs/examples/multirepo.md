# Multirepo example

In this example, we will walk through a simple multirepo setup that can be used with Earthly. The entire code of this exercise is available in the

* [examples/multirepo directory on GitHub](https://github.com/earthly/earthly/tree/main/examples/multirepo)
* This referenced [project of HTML static files](https://github.com/earthly/earthly-example-multirepo-static)
* This referenced [project of JS files](https://github.com/earthly/earthly-example-multirepo-js)

In this example, let's assume that we have a web application where HTML files are held in one repository and JS file are held in another. The complete application is a combination of both.

As such, each repository might have its self-contained `Earthfile`, specific to the programming language and setup of the repository. A third, main, repository might want to tie everything together in a complete web application.

Here is an example of what the main repository's `Earthfile` might look like

```Dockerfile
# Earthfile

FROM node:13.10.1-alpine3.11
WORKDIR /example-multirepo

docker:
    RUN npm install -g http-server
    COPY github.com/earthly/earthly-example-multirepo-static+html/* ./
    COPY github.com/earthly/earthly-example-multirepo-js+build/* ./
    EXPOSE 8080
    ENTRYPOINT ["http-server", "."]
    SAVE IMAGE example-multirepo:latest
```

Notice how the build is able to reference other repositories and copy artifacts from specific build targets. For example, the line

```Dockerfile
COPY github.com/earthly/earthly-example-multirepo-static+html/* ./
```

references the `html` target of the repository `github.com/earthly/earthly-example-multirepo-static` and copies all its artifacts in the current build environment. Earthly takes care of cloning that repository, executing its build for the `html` target and extracting the artifacts to be used here.

This command is also cache-aware. If the HEAD of the repository points to a different Git hash, Earthly knows to re-clone and build the repository again, using as much cache as relevant, depending on which files have changed.

You can also specify a specific tag or branch of the remote repository, to help keep builds consistent and avoid surprising changes. For that, you can use the syntax

```Dockerfile
COPY github.com/earthly/earthly-example-multirepo-static:v0.1.1+html/* ./
```

where `v0.1.1` is a tag or branch specifier.

To review this example with its complete code, check out the [examples/multirepo directory on GitHub](https://github.com/earthly/earthly/tree/main/examples/multirepo).
