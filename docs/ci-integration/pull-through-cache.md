# Pull Through Cache

## Introduction

Docker Hub, Quay, and other registry providers all have pull limits, and costs associated with running them. Running large builds (or many small builds, frequently) may incur excess costs, rate limiting, or both. This guide will help you set up your own "pull-through" cache to optimize traffic, and bypass the limitations imposed by registry providers.

## What Is A Pull Through Cache?

A pull through cache is a registry mirror that contains no images. When your client checks the registry for an image, the registry will either:

- Give an existing response from its cache; thereby avoiding egress (or a pull) from your registry,
- Or pull the image and its metadata from the registry on your behalf; caching it for later use.

## Running A Pull-Through Cache

To run a cache, you'll need the ability to deploy a persistent service, somewhere. This could be a dedicated instance with Docker installed, or a container in your Kubernetes cluster. While we won't be giving details for *how* to set up either of these ways to run a service, we will be sharing configuration and usage details, and how you can use it with Earthly. 

Docker has a [guide for getting a pull-through cache up and running](https://docs.docker.com/registry/recipes/mirror/#run-a-registry-as-a-pull-through-cache), and good [documentation of the available options](https://docs.docker.com/registry/configuration/). You can get the registry image (and details) [here](https://hub.docker.com/_/registry).

### Configuration & Tips

Pull-through caches run _unsecured_ by default. Add an `htpasswd` file for basic authentication, at a minimum:
```yaml
auth:
  htpasswd:
    realm: basic-realm
    path: /auth/htpasswd
```

Adding TLS is also highly recommended. you can bring your own certificates, or use the built-in LetsEncrypt support:
```yaml
http:
  tls:
    letsencrypt:
      cachefile: /certs/cachefile
      email: me@example.com
      hosts: [my.cool.mirror.horse]
```

The currently shipping `library/registry` image does not support the DNS-01 challenge yet, and [some of the LetsEncrypt challenge support is getting out of date](https://github.com/distribution/distribution/issues/3041). If you need this, there is a tracking issue; We have had success by [building the binary ourselves](https://github.com/earthly/registry/blob/3f06d1fc5d7f456b63b870b2851fd18cd2098dcf/Earthfile#L3-L11) and replacing it in the image that Docker ships.