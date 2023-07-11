*The code of this example is available in the [earthly/earthly-example-ci repository on GitHub](https://github.com/earthly/earthly-example-ci).*

An alternative to using a traditional CI system (see [part 8a](./part-8a-using-earthly-in-your-current-ci.md)) is to use Earthly CI. Earthly CI is a full-fledged continuous integration system that uses Earthly Satellites under the hood, which means that it is dramatically faster than traditional CI systems through the use of Satellites' instantly-available cache. With Earthly CI, you don't need any build specification other than a simple Earthfile.

Earthly CI is part of Earthly Cloud. You can use it for free as part of our free tier. Get started with Earthly Cloud by visiting the [sign up](https://cloud.earthly.dev/login) page, and get 6,000 build minutes/month for free.

## Getting Started with Earthly CI

To get started with Earthly CI, you need to first [sign up for free to Earthly Cloud](https://cloud.earthly.dev/login).

## Defining an Earthly CI Pipeline

An Earthly CI pipeline is defined as an Earthly target that begins with the command `PIPELINE`. Here is an example:

```Earthfile
VERSION 0.7
PROJECT my-org/my-project

FROM alpine:3.17

hello-world-pipeline:
    PIPELINE
    TRIGGER push main
    TRIGGER pr main
    BUILD +hello-world

hello-world:
    RUN echo Hello world
```

In this example, we have declared a pipeline called `hello-world-pipeline`. This pipeline will be triggered on every push to the `main` branch, and on the commits within every pull request intending to merge into the `main` branch. The pipeline will run the target `+hello-world`, which simply prints out `Hello world`.

To get this pipeline to run, you will need to add the Earthfile to a GitHub repository in the default branch (typically `main`).

After you've committed your Earthfile, go back to the Earthly Web UI to create a new project, and add the repository to the project. Please note that the `PROJECT` declaration in the Earthfile needs to match the organization name and project name in the Earthly Web UI.

The pipeline should start executing automatically. Its status can be seen in the Earthly Web UI.

![Earthly CI pipeline summary](./img/status-1.png)

![Earthly CI pipeline detail](./img/status-2.png)

Earthly CI pipelines can also be executed locally the same way as any other Earthly target. This is useful for testing the pipeline locally before pushing it to the repository.

```bash
earthly +hello-world-pipeline
```

To read more about Earthly CI, check out the [Earthly CI page](../cloud/earthly-ci.md).
