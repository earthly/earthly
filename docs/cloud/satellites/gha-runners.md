# Satellites as GitHub Actions runners

{% hint style='warning' %}
This feature is experimental.

Not recommended for production usage since it might introduce breaking changes in the future.

Feedback is welcome and much appreciated!

{% endhint %}

Earthly satellites are now bundled with a GitHub Actions runner, so they can directly pull jobs from GitHub Actions without the need of an intermediate runner.

These runners come with the Earthly CLI preinstalled and configured to use the satellite BuildKit instance, so GitHub Actions jobs will share the same satellite cache than the traditional satellite builds.

## Getting started

Satellite-based GitHub Actions runners can be enabled for a particular repository or for all repositories of a GitHub organization at once.

The integration process requires you to provide us with a GitHub token, so we can:
- register a webhook in your GitHub repository/organization to receive events associated to GitHub Actions jobs
- create GitHub self-hosted runners on demand, to process your repository/organization jobs

Follow the next steps to create such integrations:

### 1. Create a GitHub token
Both GitHub classic and fine-grained tokens are supported, but depending on the type of installation (organization-wide or single-repository), the provided token requires different scopes:

| Integration type | User type          | Classic token scopes          | Fine-grained token permissions                                       | 
|------------------|--------------------|-------------------------------|----------------------------------------------------------------------|
| Organization     | Organization admin | `admin:org_hook`, `admin:org` | `organization_hooks:write`, `organization_self_hosted_runners:write` |
| Repository       | Repository admin   | `admin:repo_hook`, `repo`     | `repository_hooks:write`, `administration:write`                     |

{% hint style='info' %}
Follow the [official docs](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens) for detailed information on how to create a GitHub token, and make sure to set an expiration long enough, since the integration won't work after the token expires.
{% endhint %}

### 2. Register the integration via CLI
Create the integration using the `earthly github add` CLI command, passing the token created in the previous step.

#### Organization integration
``` 
earthly github add \
  --org <earthly_organization> \
  --gh-org <github_organization> \
  --gh-token <github_token>
``` 

#### Single repository integration
``` 
earthly github add \
  --org <earthly_organization> \
  --gh-org <github_organization> \
  --gh-repo <github_repo> \
  --gh-token <github_token>
``` 

### 3. Configure your satellites

This feature needs to be enabled during satellite creation to be able to use it.

#### Earthly-Cloud satellites
Launch the satellite with the `enable-gha-runner` [feature-flag](https://docs.earthly.dev/earthly-cloud/satellites/managing#changing-feature-flags) enabled.
```
earthly satellite launch --feature-flag enable-gha-runner <satellite-name>
``` 

#### Self-hosted satellites
To enable the GH runner for a self-hosted satellite, set this environment entry when launching it:
```
-e RUNNER_GHA_ENABLED=true
```
also note that the satellite container must have access to the docker daemon in order to run the GitHub Actions jobs in containers:
```
-v /var/run/docker.sock:/var/run/docker.sock
```

##### Example
```shell
docker run --privileged \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v satellite-cache:/tmp/earthly:rw \
    -p 8372:8372 \
    -e EARTHLY_TOKEN=<earthly_token> \ 
    -e EARTHLY_ORG=<earthly_org_name>  \
    -e SATELLITE_NAME=<satellite_name> \
    -e SATELLITE_HOST=<satellite_host> \
    -e RUNNER_GHA_ENABLED=true \
  earthly/satellite:v0.8.13
```
{% hint style='info' %}
**Required version:** Use at least `earthly/satellite:v0.8.13
{% endhint %}

##### Logs
You should see a log message like this, when the GitHub Actions runner is enabled:
```
{...,"msg":"starting GitHub Actions job polling loop",...}
```

### 4. Configure your GitHub Actions jobs
In order to make a job run into the satellite, you'll need to set its [runs-on](https://docs.github.com/en/actions/using-workflows/workflow-syntax-for-github-actions#jobsjob_idruns-on) label as follows:

```
runs-on: [earthly-satellite#<satellite-name>]
```

#### Example
The following example runs the `+build` target in the satellite. Given that the GH runner is configured to use the satellite BuildKit instance, the persistent satellite cache is implicitly used here.
```yml
earthly-job:
  runs-on: [earthly-satellite#my-gha-satellite]
  env:
    FORCE_COLOR: 1
    EARTHLY_TOKEN: "${{ secrets.EARTHLY_TOKEN }}"
  steps:
    - uses: actions/checkout@v2
    - name: Earthly build
      run: earthly -ci +build
```

{% hint style='warning' %}
For Earthly-Cloud satellites make sure you have an [EARTHLY_TOKEN](https://docs.earthly.dev/docs/earthly-command#earthly-account-create-token) available in your [GitHub Actions secrets](https://docs.github.com/en/actions/security-guides/using-secrets-in-github-actions) store, and add it to the job environment, as shown in the previous example. Future versions will remove this requirement.

{% endhint %}
