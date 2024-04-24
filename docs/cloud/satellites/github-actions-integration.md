# GitHub actions integration (**experimental**)

This new feature allows triggering Earthly satellites directly from GHA (GitHub actions) without the need of an intermediate runner. 

Earthly Satellites embed a GitHub self-hosted runner, so they can directly pull jobs from GHA. The runner comes with the Earthly CLI preinstalled, and it's configured to use the Satellite Buildkit instance, so GHA jobs will share the same Satellite cache than the traditional Satellite builds.

Notice that this self-hosted runner can run any arbitrary GHA job, not necessarily an Earthly command, so given that it runs within the Satellite, it benefits from its persistent local storage (see ["Persistent Folders" section below](#persistent-folders)).

## Configuration

### GitHub token

### CLI

`$ earthly github add --help`

``` 
NAME:
earthly github add - Add GHA integration

USAGE:
earthly github add --org <org> [--repo <repo>] --token <token>

DESCRIPTION:
This command sets the configuration to create a new GitHub-Earthly integration, to trigger satellite builds from GHA (GitHub Actions).
From the GitHub side, integration can be done at two levels: organization-wide and per repository.
The provided token must have enough permissions to register webhooks and to create GitHub self hosted runners in those two scenarios.

OPTIONS:
--org value       The name of the Earthly organization to set an integration with. Defaults to selected organization
--gh-org value    The name of the GitHub organization to set an integration with
--gh-repo value   The name of the GitHub repository to set an integration with
--gh-token value  The GitHub token used for the integration
--help, -h        show help
```

## Satellite configuration
This feature is currently disabled by default. It is enabled in a per-satellite basis.

### Managed (Earthly Cloud) Satellites
 You can request enabling it for any of your managed satellites by sending an email to support@earthly.dev.

### Self-hosted satellites
To enable the GH runner for a self-hosted satellite, just set this environment entry when launching it: 
```
RUNNER_GHA_ENABLED=true
```
#### Example
```shell
docker run --privileged \
    -v satellite-cache:/tmp/earthly:rw \
    -p 8372:8372 \
    -e EARTHLY_TOKEN=GuFna*****nve7e \ 
    -e EARTHLY_ORG=my-org \
    -e SATELLITE_NAME=my-satellite \
    -e SATELLITE_HOST=153.65.8.0 \
    -e RUNNER_GHA_ENABLED=true \
  earthly/satellite:v0.8.9
```
{% hint style='info' %}
##### Required version
Use at least earthly/satellite:v0.8.9
{% endhint %}

#### Logs
You should see a log message like this, when the GHA runner is enabled: 
```
{...,"msg":"starting GHA job polling loop",...}
```

## GHA job definition
Job configuration is performed through [runs-on](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idruns-on) labels. In particular:

### Satellite name
```
earthly-satellite#<satellite-name>
```
This label marks the job to run on the referenced satellite, belonging to the [integrated](#cli) organization. 

Only one label starting with `earthly-satellite#` is allowed per job.
### Persistent folders
```
earthly-cache-folder#<absolute-path>
```
These labels allow defining folders whose contents will be shared across multiple builds.
This is specially useful for defining persistent caches for tools external to Earthly. 

Notice that multiple labels starting with `earthly-cache-folder#` can be set for a given job. One per persistent folder.

### Examples
#### Running an earthly job
The following example runs the +build target in the Satellite. Given that the GH runner is configured to use the Satellite Buildkit instance, the persistent satellite cache is implicitly used here.
```yml
earthly-job:
  runs-on: [earthly-satellite#my-gha-satellite]
  env:
    FORCE_COLOR: 1
  steps:
    - uses: actions/checkout@v2
    - name: Earthly build
      run: earthly -ci +build
```
#### Caching non-Earthly jobs
The following example runs maven externally to Earthly, but benefiting from the satellite storage to mount a persistent local cache for the maven artifacts:  
```yml
maven-job:
  runs-on: [earthly-satellite#my-gha-satellite, earthly-cache-folder#/root/.m2]
  env:
    FORCE_COLOR: 1
  steps:
    - uses: actions/checkout@v4
    - uses: actions/setup-java@v4
      with:
        java-version: '17'
        distribution: 'temurin'
    - name: Run the Maven verify phase
      run: mvn --batch-mode --update-snapshots verify
```
## Early access
This feature is in closed-beta at the moment. You can request early access through support@earthly.dev