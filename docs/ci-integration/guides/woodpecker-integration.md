
# Woodpecker CI integration

This example uses [Woodpecker CI](https://woodpecker-ci.org/) to build the Earthly target `+build`.


## Configuration

The project needs to be [trusted](https://woodpecker-ci.org/docs/usage/project-settings#trusted) to grant the capabilities like mounting volumes (required for the docker socket). We also need to include the `earthly/earthly` image in the list of images that are allowed to run in [privileged mode](https://woodpecker-ci.org/docs/administration/server-config#woodpecker_escalate)



```yml
#.woodpecker.yml
pipeline:
  earthly:
    image: earthly/earthly:v0.8.11
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
    environment:
      - FORCE_COLOR=1
      - EARTHLY_EXEC_CMD="/bin/sh" 
    secrets: [REGISTRY, REGISTRY_USER, REGISTRY_PASSWORD]
    commands:
     - docker login -u $${REGISTRY_USER} -p $${REGISTRY_PASSWORD} $${REGISTRY}
     - earthly --ci --push +build
```

For a complete guide on CI integration see the [CI integration guide](../overview.md).
