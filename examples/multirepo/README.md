# Multirepo example

In this example, we will walk through a simple multirepo setup that can be used with Earthly. The entire code of this exercise is available in the

* This directory
* This referenced [project of HTML static files](https://github.com/earthly/earthly-example-multirepo-static)
* This referenced [project of JS files](https://github.com/earthly/earthly-example-multirepo-js)

In this example, let's assume that we have a web application where HTML files are held in one repository and JS file are held in another. The complete application is a combination of both.

As such, each repository might have its self-contained `Earthfile`, specific to the programming language and setup of the repository. A third, main, repository might want to tie everything together in a complete web application.

## Run

To run this build execute

```bash
earthly +docker
```

in this directory, or, without cloning the Earthly repo, run this anywhere

```
earthly github.com/earthly/earthly/examples/multirepo:main+docker
```

Then, run the resulting container:

```
docker run --rm -p 127.0.0.1:8080:8080 earthly/examples:multirepo
```

and load `http://127.0.0.1:8080` in your browser.

## Explanation

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

## Compare

Compare this example with the example presented in [monorepo](../monorepo).
