# Building from source

Earthly utilizes [dogfooding](https://en.wikipedia.org/wiki/Eating_your_own_dog_food) to build itself.
To understand the structure of the source code you may utilize this [page](../examples/earthly).
So before building from source you need to install Earthly from binary [see here](installation.md) .
After installing Earthly you can clone the repository using git.
Once cloned you can execute earth +all within the root directory to build the programm.
Note that if you are running under windows you should clone within the linux file system to prevent potential access problems.
Once your build is done, you are able to find the binaries within the *build* directory.
That directory contains directories for your target platform (e.g. linux) and your target architecture (e.g. amd64).
An example output of running tree within the build directory may return by default:
```bash
├── darwin
│   └── amd64
│       └── earth
└── linux
    └── amd64
        └── earth

4 directories, 2 files
```

After building you may want to remove the original binary and replace it with the binary just generated.
If no errors occur you should now have a fully functioning executable compiled from source.
If however you are facing any issue, the following aggregates a list of common problems (e.g. on installation within a new environment such as WSL2):

## Cloning a public repository fails

By default, cloning a git repository in Earthly utilizes ssh and this might be problematic if unconfigured.
You have two options:

### Configuring ssh
Therefore you need to configure ssh-agent (see [here](../guides/auth) before usage.
The following commands might help accordingly:

```bash
# Starting the ssh agent
eval $(ssh-agent)

# Generate a key. The following is a mere example
ssh-keygen  -t rsa -b 4096 -C "my-email@example.org"

# Adds the key to the agent
ssh-add

# Copies public key to clipboard so you can add it to your github account manually
# You could use any ssh (public) key here though
# Note that clip.exe assumes WSL2, you might want to change it to (e.g.) xclip 
cat ~/.ssh/id_rsa.pub | clip.exe
```

Note that adding a SSH Key to your account on github is also explained in more details for more platforms within the [Github Docs](https://docs.github.com/en/free-pro-team@latest/github/authenticating-to-github/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent).

### Utilizing HTTP
Both the environment variable GIT_URL_INSTEAD_OF and the [git configuration](../earth-config).
Allow rewriting ssh access of github to https.
As demonstrated in [the CI integration guide](../guides/ci-integration) the variable can be utilized as:
```bash
export GIT_URL_INSTEAD_OF="https://github.com/=git@github.com:"
```
The configuration option works similar and is documented [here](../earth-config).

## Running a full build stops after "SUCCESS", but does not exit
After starting a run of "earth +all" the console displays a green line which states
```=========================== SUCCESS ===========================```
It does not denote the end of the build process however and some additional work might still be required.
However it might take a while to be displayed.
For example some docker layers may still need to be loaded, which takes a while.
Depending on the available systems capacity, this may take a while and requires patience.
