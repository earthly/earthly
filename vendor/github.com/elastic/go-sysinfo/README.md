# go-sysinfo

[![Build Status](https://beats-ci.elastic.co/job/Library/job/go-sysinfo-mbp/job/main/badge/icon)](https://beats-ci.elastic.co/job/Library/job/go-sysinfo-mbp/job/main/)
[![Go Documentation](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)][godocs]

[godocs]: http://godoc.org/github.com/elastic/go-sysinfo

go-sysinfo is a library for collecting system information. This includes
information about the host machine and processes running on the host.

The available features vary based on what has been implemented by the "provider"
for the operating system. At runtime you check to see if additional interfaces
are implemented by the returned `Host` or `Process`. For example:

```go
process, err := sysinfo.Self()
if err != nil {
	return err
}

if handleCounter, ok := process.(types.OpenHandleCounter); ok {
	count, err := handleCounter.OpenHandleCount()
	if err != nil {
		return err
	}
	log.Printf("%d open handles", count)
}
```

These tables show what methods are implemented as well as the extra interfaces
that are implemented.

| `Host` Features  | Darwin | Linux | Windows | AIX |
|------------------|--------|-------|---------|-----|
| `Info()`         | x      | x     | x       | x   |
| `Memory()`       | x      | x     | x       | x   |
| `CPUTimer`       | x      | x     | x       | x   |
| `LoadAverage`    | x      | x     |         |     |
| `VMStat`         |        | x     |         |     |
| `NetworkCounters`|        | x     |         |     |

| `Process` Features     | Darwin | Linux | Windows | AIX |
|------------------------|--------|-------|---------|-----|
| `Info()`               | x      | x     | x       | x   |
| `Memory()`             | x      | x     | x       | x   |
| `User()`               | x      | x     | x       | x   |
| `Parent()`             | x      | x     | x       | x   |
| `CPUTimer`             | x      | x     | x       | x   |
| `Environment`          | x      | x     |         | x   |
| `OpenHandleEnumerator` |        | x     |         |     |
| `OpenHandleCounter`    |        | x     |         |     |
| `Seccomp`              |        | x     |         |     |
| `Capabilities`         |        | x     |         |     |
| `NetworkCounters`      |        | x     |         |     |

### GOOS / GOARCH Pairs

This table lists the OS and architectures for which a "provider" is implemented.

| GOOS / GOARCH  | Requires CGO | Tested |
|----------------|--------------|--------|
| aix/ppc64      | x            |        |
| darwin/amd64   | optional *   | x      |
| darwin/arm64   | optional *   | x      |
| linux/386      |              |        |
| linux/amd64    |              | x      |
| linux/arm      |              |        |
| linux/arm64    |              |        |
| linux/mips     |              |        |
| linux/mips64   |              |        |
| linux/mips64le |              |        |
| linux/mipsle   |              |        |
| linux/ppc64    |              |        |
| linux/ppc64le  |              |        |
| linux/riscv64  |              |        |
| linux/s390x    |              |        |
| windows/amd64  |              | x      |
| windows/arm64  |              |        |
| windows/arm    |              |        |

* On darwin (macOS) host information like machineid and process information like memory, cpu, user and starttime require cgo.
