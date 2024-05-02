# Self-signed certificates

This guide will demonstrate the use of a private registry using self-signed certificates in conjunction with Earthly.

For information about configuring the registry itself, see the [Docker Registry deployment documentation](https://docs.docker.com/registry/deploying/).

## Create an Earthfile

No special considerations are needed in the Earthfile itself. You can use `SAVE IMAGE` just like any other repository.

```
FROM alpine:3.18

build:
    RUN echo "Hello from Earthly!" > motd
    ENTRYPOINT cat motd
    SAVE IMAGE --push <registry-hostname>/hello-earthly:with-love
```

## Add certificates to Earthly

Set the following configuration options in your [Earthly config](../../earthly-config/earthly-config.md).

```yaml
global:
  buildkit_additional_args: ["-v", "<absolute-path-to-ca-file>:/etc/config/add.ca"]
  buildkit_additional_config: |
    [registry."<registry-hostname>"]
      ca=["/etc/config/add.ca"]
```

Where `<absolute-path-to-ca-file>` is the location of the CA certificate you wish to add and `<registry-hostname>` is the hostname of the registry. The quotes are not a mistake, and should be left in.

## Insecure registries

For testing purposes, you can also define insecure registries for Earthly to access. Note that the non-test use of insecure registries is strongly discouraged due to the risk of man-in-the-middle (MITM) attacks.

To configure Earthly to use an insecure registry, use the following [Earthly config](../../earthly-config/earthly-config.md) settings.

```yaml
global:
  buildkit_additional_config: |
    [registry."<registry-hostname>"]
      insecure = true
```

In addition, you will need to specify the `--insecure` flag in any `SAVE IMAGE` command.  Again, the quotes are not a mistake, and should be left in.

```
FROM alpine:3.18

build:
    RUN echo "Hello from Earthly!" > motd
    ENTRYPOINT cat motd
    SAVE IMAGE --push --insecure <registry-hostname>/hello-earthly:with-love
```

{% hint style='danger' %}
##### Note

The `http` and `insecure` settings are typically mutually exclusive. Setting `insecure=true` should only be used when the registry is https and is configured with an insecure certificate.
Setting `http=true` is only for the case where a standard http-based registry is used (i.e. no SSL encryption). If both are set BuildKit will attempt to connect to the registry using either http (port 80), or https (port 443).

{% endhint %}


## Other BuildKit options

Other settings for configuring registries in Earthly via [BuildKit options](https://github.com/moby/buildkit/blob/master/docs/buildkitd.toml.md) can be seen below.

```yaml
global:
  buildkit_additional_config: |
    [registry."<registry-hostname>"]
      mirrors = ["<mirror>"]
      http = true|false
      insecure = true|false
      ca=["<ca-path-pem>"]
      [[registry."<registry-hostname>".keypair]]
        key="<key-path-pem>"
        cert="<cert-path-pem>"
```
