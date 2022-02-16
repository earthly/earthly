# Running the build

In the example `Earthfile` we have defined two explicit targets: `+build` and `+docker`. When we run Earthly, we can tell it to execute either target by passing the target name. In this case our docker target calls on our build target, so we can run both with:

```bash
earthly +docker
```

The output might look like this:

![Earthly build output](../guides/img/go-example.png)

Notice how to the left of `|`, within the output, we can see some targets like `+base`, `+build` and `+docker` . Notice how the output is interleaved between `+docker` and `+build`. This is because the system executes independent build steps in parallel. The reason this is possible effortlessly is because only very few things are shared between the builds of the recipes and those things are declared and obvious. The rest is completely isolated.

In addition, notice how even though the base is used as part of both `build` and `docker`, it is only executed once. This is because the system deduplicates execution, where possible.

Furthermore, the fact that the `docker` target depends on the `build` target is visible within the command `COPY +build/...`. Through this command, the system knows that it also needs to build the target `+build`, in order to satisfy the dependency on the artifact.

Finally, notice how the output of the build (the docker image and the files) are only written after the build is declared a success. This is due to another isolation principle of Earthly: a build either succeeds completely or it fails altogether.

Once the build has executed, we can run the resulting docker image to try it out:

{% method %}
{% sample lang="Go" %}

```
docker run --rm go-example:latest
```

{% sample lang="JavaScript" %}

```
docker run --rm js-example:latest
```

{% sample lang="Java" %}

```
docker run --rm java-example:latest
```

{% sample lang="Python" %}

```
docker run --rm python-example:latest
```

{% endmethod %}
