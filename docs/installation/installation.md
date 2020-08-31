# Installation


## Pre-requisites (all platforms)

* [Docker](https://docs.docker.com/install/)
* [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

## Install earth

### Linux

```bash
sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/latest/download/earth-linux-amd64 -O /usr/local/bin/earth && chmod +x /usr/local/bin/earth'
```

Alternatively, you may also download the linux binary from [the releases page](https://github.com/earthly/earthly/releases), rename it to `earth` and place it in your `bin`.

### Mac

#### Homebrew

```bash
brew install earthly
```

#### Script

```bash
wget https://github.com/earthly/earthly/releases/latest/download/earth-darwin-amd64 -O /usr/local/bin/earth && chmod +x /usr/local/bin/earth
```

Alternatively, you may also download the darwin binary from [the releases page](https://github.com/earthly/earthly/releases), rename it to `earth` and place it in your `bin`.

### Windows via WSL (**beta**)

Earthly on Windows requires [Docker Desktop WSL2 backend](https://docs.docker.com/docker-for-windows/wsl/). Under `wsl`, run the following to install `earth`.

```bash
sudo /bin/sh -c 'wget https://github.com/earthly/earthly/releases/latest/download/earth-linux-amd64 -O /usr/local/bin/earth && chmod +x /usr/local/bin/earth'
```

### CI

For instructions on how to install `earth` for CI use, see the [CI integration guide](../guides/ci-integration.md).

## Configuration

If you use SSH-based git authentication, then your git credentials will just work with Earthly. Read more about [git auth](../guides/auth).

For a full list of configuration options, see the [Configuration reference](../earth-config/earth-config.md)

## Verify installation

To verify that the installation works correctly, you can issue a simple build of an existing hello-world project

```bash
earth github.com/earthly/hello-world+hello
```

You should see the output

```
github.com/earthly/hello-world:master+hello | --> RUN [echo 'Hello, world!']
github.com/earthly/hello-world:master+hello | Hello, world!
github.com/earthly/hello-world:master+hello | Target github.com/earthly/hello-world:master+hello built successfully
=========================== SUCCESS ===========================
```

## Syntax highlighting

### VS Code extension

[<img src="./img/vscode-plugin.png" alt="Earthfile Syntax Highlighting" width="457" />](https://marketplace.visualstudio.com/items?itemName=earthly.earthfile-syntax-highlighting)

Add [Earthfile Syntax Highlighting](https://marketplace.visualstudio.com/items?itemName=earthly.earthfile-syntax-highlighting) to VS Code.

```
ext install earthly.earthfile-syntax-highlighting
```
