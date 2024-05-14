# Pull Through Cache

## Introduction

Docker Hub, Quay, and other registry providers have pull limits, and costs associated with using them.
Running large builds (or many small builds, frequently) may incur costs, rate limiting, or both.
This guide will help you set up your own "pull-through" cache to reduce network traffic, and bypass the limitations imposed by registry providers.

## What Is A Pull Through Cache?

A pull through cache is a registry mirror that contains no images. When your client checks the registry for an image, the registry will either:

- Give an existing response from its cache; thereby avoiding egress (or a pull) from your registry,
- Or pull the image and its metadata from the registry on your behalf; caching it for later use.

## Running A Pull-Through Cache

To run a cache, you'll need the ability to deploy a persistent service, somewhere. This could be a dedicated instance with Docker installed, or a container in your Kubernetes cluster.

There are multiple ways to setup a registry -- Docker, for example, has a [guide for using the registry as a pull through cache](https://docs.docker.com/docker-hub/mirror/),
as well as documentation for the [available options](https://distribution.github.io/distribution/about/configuration/), and other details under the [registry image](https://hub.docker.com/_/registry).

Documenting all the possible ways to set up a pull through cache is beyond the scope of this document; however, it does include a [quick getting-started section](#insecure-docker-hub-cache-example) for those who wish
to run an insecure pull through cache.

### Configuration & Tips

####  Set Up Mirror Authentication

Pull-through caches run _unsecured_ by default. Add an `htpasswd` file for basic authentication, at a minimum:
```yaml
auth:
  htpasswd:
    realm: basic-realm
    path: /auth/htpasswd
```

#### Set Up Mirror TLS

Adding TLS is also highly recommended. you can bring your own certificates, or use the built-in LetsEncrypt support:
```yaml
http:
  tls:
    letsencrypt:
      cachefile: /certs/cachefile
      email: me@example.com
      hosts: [my.cool.mirror.horse]
```

The currently shipping `library/registry` image does not support the DNS-01 challenge yet, and [some of the LetsEncrypt challenge support is getting out of date](https://github.com/distribution/distribution/issues/3041). If you need this, there is a [tracking issue](https://github.com/docker/distribution-library-image/issues/96); We have had success by [building the binary ourselves](https://github.com/earthly/registry/blob/3f06d1fc5d7f456b63b870b2851fd18cd2098dcf/Earthfile#L3-L11) and replacing it in the image that Docker ships.

#### Use An Insecure Mirror

By default, Earthly expects your mirror to be using TLS. While this is not recommended, you can use an unsecured mirror by specifying the following config in the `buildkit_additional_config` setting:

```yaml
global:
  buildkit_additional_config: |
    [registry."<upstream>"]
      mirrors = ["<mirror>"]

    [registry."<mirror>"]
      insecure = true
```

Where `<mirror>` is the host/port of your mirror, and `<upstream>` is the address of the registry you are intending to mirror.

## Insecure Docker Hub Cache Example

This section contains a quick-start guide for running an insecure pull through cache using docker's `registry` container.

This guide assumes you are running on a trusted network with one computer acting as a server,
and the other as your development workstation.

In these examples, we will assume the server has an IP set to `192.168.0.80`.

### Setting up the pull through cache

First connect to the server where you will be running the cache (e.g. `ssh 192.168.0.80`),
and create a file under `~/.docker-registry-config.yml` containing:

```yaml
version: 0.1
log:
  fields:
    service: registry
storage:
  cache:
    blobdescriptor: inmemory
  filesystem:
    rootdirectory: /var/lib/registry
http:
  addr: :5000
  headers:
    X-Content-Type-Options: [nosniff]
health:
  storagedriver:
    enabled: true
    interval: 10s
    threshold: 3
proxy:
  remoteurl: https://registry-1.docker.io
  username: [username]
  password: [password]
```

*Note that you'll need to replace `[username]` and `[password]` with your dockerhub credentials.*

Next, start the registry container with:

```bash
docker run --rm --network host -d --name docker-registry -v $HOME/.docker-registry-config.yml:/root/config.yml registry.hub.docker.com/library/registry:2 registry serve /root/config.yml
```

You can then verify the registry is running by tailing logs with:

```bash
docker logs --follow docker-registry
```

{% hint style='info' %}
You may want to leave a second terminal window open to display the logs while you work on the following sections;
this will make it more obvious when the cache is being used.
{% endhint %}

The rest of the guide focus on configuring your workstation to use this cache.

### Configuring Docker to Use the Cache

To configure docker to use this cache as a mirror, edit the `/etc/docker/daemon.json` file, and add:

```json
{
  "registry-mirrors" : ["http://192.168.0.80:5000"]
}
```

Then restart docker:

```bash
sudo service docker restart
```

Next you should be able to pull an image (e.g. `docker pull alpine:3.18`), which should use the cache.

#### Verifying the Cache is Actually Working (Optional)

If you want to verify the cache is working, you can block access to dockerhub on your workstation by adding

```
0.0.0.0 index.docker.io auth.docker.io registry-1.docker.io dseasb33srnrn.cloudfront.net production.cloudflare.docker.com
```

to your `/etc/hosts` file.

If `/etc/docker/daemon.json` was not correctly configured, you should see an error such as:

```
Error response from daemon: Get "https://registry-1.docker.io/v2/": dial tcp 0.0.0.0:443: connect: connection refused
```

If the cache is correctly configured, the pull command should work, and you should see logs on your server under the docker-registry container:

```
192.168.0.126 - - [22/Mar/2022:19:10:39 +0000] "HEAD /v2/library/alpine/manifests/3.15 HTTP/1.1" 200 1638 "" "docker/20.10.12 go/go1.16.12 git-commit/459d0df kernel/5.13.0-35-generic os/linux arch/amd64 UpstreamClient(Docker-Client/20.10.12 \\(linux\\))"
```

### Configuring Earthly to Use the Cache

To configure earthly to use the cache, you must edit `~/.earthly/config.yml` to include:

```yaml
global:
  buildkit_additional_config: |
    [registry."docker.io"]
      mirrors = ["192.168.0.80:5000"]
    [registry."192.168.0.80:5000"]
      insecure = true
```

The next time earthly is run, it will detect the configuration change and will restart the `earthly-buildkitd` container to reflect these settings.

You can force these settings to be applied, and verify the mirror appears in the BuildKit config by running:

```bash
earthly bootstrap && docker exec earthly-buildkitd cat /etc/buildkitd.toml
```
