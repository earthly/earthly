# Release checklist

* Make sure you have a GITHUB_TOKEN set.
  ```bash
  export GITHUB_TOKEN="..."
  ```
* Choose the next [release tag](https://github.com/vladaionescu/earthly/releases).
  ```bash
  export RELEASE_TAG="v..."
  ```
* Run
  ```bash
  earth --build-arg RELEASE_TAG --secret GITHUB_TOKEN --push -P +release
  ```

### Updating buildkitd

Update the vladaionescu/buildkit fork, then run this locally to rebuild and push to docker:

```bash
DOCKER_BUILDKIT=1 docker build -t earthly/buildkit:latest --target buildkit-buildkitd-linux .
docker push earthly/buildkit:latest
```

Then edit buildkitd/build.earth in this repo to point to the new sha256 of the earthly/buildkit image.
