module github.com/earthly/earthly

go 1.13

require (
	github.com/antlr/antlr4 v0.0.0-20200225173536-225249fdaef5
	github.com/containerd/containerd v1.4.0-0
	github.com/docker/cli v0.0.0-20200227165822-2298e6a3fe24
	github.com/docker/distribution v2.7.1-0.20190205005809-0d3efadf0154+incompatible
	github.com/docker/docker v1.14.0-0.20190319215453-e7b5f7dbe98c
	github.com/fatih/color v1.9.0
	github.com/godbus/dbus v4.1.0+incompatible // indirect
	github.com/golang/protobuf v1.3.3
	github.com/google/pprof v0.0.0-20191025152101-a8b9f9d2d3ce // indirect
	github.com/ianlancetaylor/demangle v0.0.0-20181102032728-5e5cf60278f6 // indirect
	github.com/moby/buildkit v0.7.1-0.20200331041515-6e725bdbb6ce
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/opencontainers/image-spec v1.0.1
	github.com/otiai10/copy v1.1.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.4.2
	github.com/urfave/cli/v2 v2.1.1
	golang.org/x/arch v0.0.0-20190927153633-4e8777c89be4 // indirect
	golang.org/x/crypto v0.0.0-20200302210943-78000ba7a073
	golang.org/x/lint v0.0.0-20190930215403-16217165b5de // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
)

replace github.com/hashicorp/go-immutable-radix => github.com/tonistiigi/go-immutable-radix v0.0.0-20170803185627-826af9ccf0fe

replace github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305

replace github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200227195959-4d242818bf55

replace github.com/docker/docker => github.com/docker/docker v1.4.2-0.20200227233006-38f52c9fec82
