# earthfile-syntax-highlighting

<div align="center"><img alt="Earthly" width="700px" src="https://github.com/earthly/earthly/raw/main/img/logo-banner-white-bg.png" /></div>

Syntax highlighting for [Earthly](https://earthly.dev) Earthfiles.

For an introduction of Earthly see the [Earthly GitHub repository](https://github.com/earthly/earthly) or the [Earthly documentation](https://docs.earthly.dev).

## Release Notes

### 0.0.9

* Add highlighting for IF.

### 0.0.8

* Properly handle comments in conjunction with line continuation.

### 0.0.7

* Fix README image.

### 0.0.6

* Fix quote escaping.

### 0.0.5

* Fix `FROM DOCKERFILE`
* Fix highlighting for target and artifact refs in edge cases (eg `g++`)
* Make case-sensitive
* Add highlighting for `WITH DOCKER` ... `END`

### 0.0.4

* Add highlighting for `FROM DOCKERFILE`

### 0.0.3

* Add highlighting for `HEALTHCHECK`.
* Switch from `build.earth` to `Earthfile`.

### 0.0.2

Add screenshot in the README.

### 0.0.1

Initial release of earthfile-syntax-highlighting.
