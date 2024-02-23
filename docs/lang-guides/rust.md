# Rust

This page will help you use Earthly if you are using Rust.

## Step 1: Import the Rust library

To get started, import the rust library as shown below. Note that the `--global-cache` flag is currently required to allow for adequate caching of the Rust toolchain. This flag will be implied in a future release.

```Dockerfile
VERSION 0.8

IMPORT github.com/earthly/lib/rust:2.2.11 AS rust
```

## Step 2: Initialize the Rust toolchain

Next, initialize the Rust toolchain via `rust+INIT`. This will install any necessary dependencies that `rust+CARGO` needs underneath.

```Dockerfile
install:
  FROM rust:1.73.0-bookworm
  RUN apt-get update -qq
  RUN apt-get install --no-install-recommends -qq autoconf autotools-dev libtool-bin clang cmake bsdmainutils
  RUN rustup component add clippy
  RUN rustup component add rustfmt
  # Call +INIT before copying the source file to avoid installing depencies every time source code changes. 
  # This parametrization will be used in future calls to functions of the library
  DO rust+INIT --keep_fingerprints=true
```

## Step 3: Build your Rust project

Now you can build your Rust project. Collect the necessary sources and call `rust+CARGO` to build your project.

```Dockerfile
source:
  FROM +install
  COPY --keep-ts Cargo.toml Cargo.lock ./
  COPY --keep-ts --dir package1 package2  ./

build:
  FROM +source
  DO rust+CARGO --args="build --release" --output="release/[^/\.]+"
  SAVE ARTIFACT ./target/release/*
```

Notice the need for the `--keep-ts` flag when copying the source files. This is necessary to ensure that the timestamps of the source files are preserved such that Rust's incremental compilation works correctly.

Additionally, because cargo does not make a good distinction between intermediate and final artifacts, we use the `--output` flag to specify which files should be extracted from the cache at the end of the operation.

## Finally

For a complete Earthfile example on how to use Rust in Earthly, visit the [rust example directory on GitHub](https://github.com/earthly/earthly/tree/main/examples/rust).

See also the reference documentation for [lib/rust](https://github.com/earthly/lib/tree/main/rust), to understand the different parameters used with `rust+INIT` and `rust+CARGO`.
