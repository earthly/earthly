# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Deprecated

### Removed

### Fixed

## [1.9.0]

### Added

- Replace pkg/errors with Go 1.13 native errors. [#123](https://github.com/elastic/go-sysinfo/pull/123)
- Add OS family mappings for `rocky`, `openEuler`, and `almalinux`. [#143](https://github.com/elastic/go-sysinfo/pull/143)

### Changed

- Remove custom sysctl implementation and partial cgo requirement. [#135](https://github.com/elastic/go-sysinfo/pull/135)
- Changes on the `Host` and `LoadAverage` interfaces, now implemented by default on Linux and Darwin platforms. [#140](https://github.com/elastic/go-sysinfo/pull/140)

## [1.8.1]

### Fixed

- Report OS name as Windows 11 when version is >= 10.0.22000. [#118](https://github.com/elastic/go-sysinfo/issues/118) [#121](https://github.com/elastic/go-sysinfo/pull/121)

## [1.8.0]

### Added

- Added the Oracle Linux ("ol") platform to the "redhat" OS family. [#54](https://github.com/elastic/go-sysinfo/issues/54) [#115](https://github.com/elastic/go-sysinfo/pull/115)
- Added the Linux Mint ("linuxmint") platform to the "debian" OS family. [#52](https://github.com/elastic/go-sysinfo/issues/52)

### Changed

- Updated module to require Go 1.17. [#111](https://github.com/elastic/go-sysinfo/pull/111)
- The boot time value for Windows is now rounded to the nearest second to provide a more stable value. [#53](https://github.com/elastic/go-sysinfo/issues/53) [#114](https://github.com/elastic/go-sysinfo/pull/114)

### Fixed

- Fix handling of environment variables without values on macOS. [#94](https://github.com/elastic/go-sysinfo/pull/94)
- Fix build tags on AIX provider such that CGO is required. [#106](https://github.com/elastic/go-sysinfo/issues/106)

## [1.7.1] - 2021-10-11

### Fixed

- Fixed getting OS info when an unsupported file or directory is found matching /etc/\*-release [#102](https://github.com/elastic/go-sysinfo/pull/102)

## [1.7.0] - 2021-02-22

### Added

- Add per-process network stats [#96](https://github.com/elastic/go-sysinfo/pull/96)

## [1.6.0] - 2021-02-09

### Added

- Add darwin/arm64 support (Apple M1). [#91](https://github.com/elastic/go-sysinfo/pull/91)

## [1.5.0] - 2021-01-14

### Added

- Added os.type field to host info. [#87](https://github.com/elastic/go-sysinfo/pull/87)

## [1.4.0] - 2020-07-21

### Added

- Add AIX support [#77](https://github.com/elastic/go-sysinfo/pull/77)
- Added detection of containerized cgroup in Kubernetes [#80](https://github.com/elastic/go-sysinfo/pull/80)

## [1.3.0] - 2020-01-13

### Changed

- Convert NetworkCountersInfo maps to uint64 [#75](https://github.com/elastic/go-sysinfo/pull/75)

## [1.2.1] - 2020-01-03

### Fixed

- Create a `sidToString` function to deal with API changes in various versions of golang.org/x/sys/windows. [#74](https://github.com/elastic/go-sysinfo/pull/74)

## [1.2.0] - 2019-12-09

### Added

- Added detection of systemd cgroups to the `IsContainerized` check. [#71](https://github.com/elastic/go-sysinfo/pull/71)
- Added networking counters for Linux hosts. [#72](https://github.com/elastic/go-sysinfo/pull/72)

## [1.1.1] - 2019-10-29

### Fixed

- Fixed an issue determining the Linux distribution for Fedora 30. [#69](https://github.com/elastic/go-sysinfo/pull/69)

## [1.1.0] - 2019-08-22

### Added

- Add `VMStat` interface for Linux. [#59](https://github.com/elastic/go-sysinfo/pull/59)

## [1.0.2] - 2019-07-09

### Fixed

- Fixed a leak when calling the CommandLineToArgv function. [#51](https://github.com/elastic/go-sysinfo/pull/51)
- Fixed a crash when calling the CommandLineToArgv function. [#58](https://github.com/elastic/go-sysinfo/pull/58)

## [1.0.1] - 2019-05-08

### Fixed

- Add support for new prometheus/procfs API. [#49](https://github.com/elastic/go-sysinfo/pull/49)

## [1.0.0] - 2019-05-03

### Added

- Add Windows provider implementation. [#22](https://github.com/elastic/go-sysinfo/pull/22)
- Add Windows process provider. [#26](https://github.com/elastic/go-sysinfo/pull/26)
- Add `OpenHandleEnumerator` and `OpenHandleCount` and implement these for Windows. [#27](https://github.com/elastic/go-sysinfo/pull/27)
- Add user info to Process. [#34](https://github.com/elastic/go-sysinfo/pull/34)
- Implement `Processes` for Darwin. [#35](https://github.com/elastic/go-sysinfo/pull/35)
- Add `Parent()` to `Process`. [#46](https://github.com/elastic/go-sysinfo/pull/46)

### Fixed

- Fix Windows registry handle leak. [#33](https://github.com/elastic/go-sysinfo/pull/33)
- Fix Linux host ID by search for older locations for the machine-id file. [#44](https://github.com/elastic/go-sysinfo/pull/44)

### Changed

- Changed the host containerized check to reduce false positives. [#42](https://github.com/elastic/go-sysinfo/pull/42) [#43](https://github.com/elastic/go-sysinfo/pull/43)

[Unreleased]: https://github.com/elastic/go-sysinfo/compare/v1.8.1...HEAD
[1.8.1]: https://github.com/elastic/go-sysinfo/releases/tag/v1.8.1
[1.8.0]: https://github.com/elastic/go-sysinfo/releases/tag/v1.8.0
[1.7.1]: https://github.com/elastic/go-sysinfo/releases/tag/v1.7.1
[1.7.0]: https://github.com/elastic/go-sysinfo/releases/tag/v1.7.0
[1.6.0]: https://github.com/elastic/go-sysinfo/releases/tag/v1.6.0
[1.5.0]: https://github.com/elastic/go-sysinfo/releases/tag/v1.5.0
[1.4.0]: https://github.com/elastic/go-sysinfo/releases/tag/v1.4.0
[1.3.0]: https://github.com/elastic/go-sysinfo/releases/tag/v1.3.0
[1.2.1]: https://github.com/elastic/go-sysinfo/releases/tag/v1.2.1
[1.2.0]: https://github.com/elastic/go-sysinfo/releases/tag/v1.2.0
[1.1.1]: https://github.com/elastic/go-sysinfo/releases/tag/v1.1.0
[1.1.0]: https://github.com/elastic/go-sysinfo/releases/tag/v1.1.0
[1.0.2]: https://github.com/elastic/go-sysinfo/releases/tag/v1.0.2
[1.0.1]: https://github.com/elastic/go-sysinfo/releases/tag/v1.0.1
[1.0.0]: https://github.com/elastic/go-sysinfo/releases/tag/v1.0.0
