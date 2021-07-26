# Other Platforms

While Earthly works best on Linux-based hosts, you can use macOS and Windows hosts to do your builds too. These builds will _still_ run a Linux environment within Earthly to provide services. Below are the special considerations you may need to take for your host platform.  

## macOS

### Dependencies

To install `git`, the easiest way is to install the "XCode Command Line Tools". If you open up `Terminal`, and type:

```go
git --version
```

Then macOS will prompt you to install these tools. You can also use the `git` provided installer or Homebrew, if you prefer. [Details can be found here](https://git-scm.com/download/mac).

To install `docker`, [download and install Docker CE](https://hub.docker.com/editions/community/docker-ce-desktop-mac). Be sure to grab the correct installer depending on your CPU architecture.

### Installing Earthly

If the computer is persistent (i.e. one that will live longer than a couple builds), it is probably best to [install Earthly directly](../alt-installation.md#macos-users), and pin it to a speciic version. You can do this by:

```bash
brew install earthly/earthly/earthly@v0.5.19 && earthly bootstrap
```


## Windows (WSL2)

On Windows, it is possible to use the Linux installation within WSL2. The only unique thing you will need to do is [install the Windows version of Docker]((https://hub.docker.com/editions/community/docker-ce-desktop-windows)).

After that, you can follow the [regular CI integration instructions](overview.md).

## Windows (Native, Experimental)

### Dependencies

To install `git`, use the  [MSI installer](https://gitforwindows.org/). This will provide `git`, and a Bash shell; which may prove more natural for using Earthly. You may also use your package manager of choice.

To install `docker`, [download and install Docker CE](https://hub.docker.com/editions/community/docker-ce-desktop-windows). Both the HyperV and WSL2 backends are supported, but the WSL2 one is very likely to be faster.

### Installing Earthly

Fow now, you will need to [manually download and install it.](../alt-installation.md#native-windows), and ensure it is within your `$PATH`. 