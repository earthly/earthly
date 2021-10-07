# Monorepo example

In this example, we will walk through a simple monorepo setup that can be used with Earthly.

In this example, let's assume that we have a web application where HTML files are held in one directory and JS files are held in another. The complete application is a combination of both.

As such, each directory might have its self-contained `Earthfile`, specific to the programming language and setup of the repository. A third, `frontend`, repository might want to tie everything together in a complete web application.

## Run

To run this build execute

```bash
earthly +all
```

in this directory, or, without cloning the Earthly repo, run this anywhere

```
earthly github.com/earthly/earthly/examples/monorepo:main+all
```

Then, run the resulting container:

```
docker run --rm -p 127.0.0.1:8080:8080 earthly/examples:monorepo
```

and load `http://127.0.0.1:8080` in your browser.

## Explanation

Notice how the build in `frontend/Earthfile` is able to reference other directories and copy artifacts from specific build targets. For example, the line

```Dockerfile
COPY ../html+html/* ./
```

references the `html` target of the directory `../html` and copies all its artifacts in the current build environment. Earthly executes its build for the `html` target and extracts the artifacts to be used here.

## Compare

Compare this example with the example presented in [multirepo](../multirepo).
