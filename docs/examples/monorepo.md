# Monorepo example

In this example, we will walk through a simple monorepo setup that can be used with Earthly. The entire code of this exercise is available in the [examples/monorepo directory on GitHub](https://github.com/vladaionescu/earthly/tree/master/examples/monorepo).

In this example, let's assume we've organized our monorepo such that each root-level directory is a sub-project.

```
.
├── proj1
└── proj2
```

As such, each sub-project might have its own self-contained `build.earth` file, specific to the programming language and setup of the sub-project. A root-level `build.earth` can tie everything together or surface important targets of the sub-projects.

```
.
├── build.earth
├── proj1
│   └── build.earth
└── proj2
    └── build.earth
```

Here is an example of a possible root-level `build.earth` file, where an `all` target simply calls the `+docker` target of each sub-project.

```Dockerfile
# build.earth

FROM alpine:3.11

all:
    BUILD ./proj1+docker
    BUILD ./proj2+docker
```

Note that the directory hierarchy may be as vast and deeply-nested as appropriate for your setup. In addition, build targets within projects may depend on targets from other projects. As an example, consider the case where one project builds a library and another takes the built library as an artifact and imports it in order to use it in its own build.

Further, throught the use of caching, the build setup is able to infer automatically which sub-projects to rebuild because of local changes, and which ones to reuse cache for.

To review this example with its complete code, check out the [examples/monorepo directory on GitHub](https://github.com/vladaionescu/earthly/tree/master/examples/monorepo).
