module github.com/earthly/earthly

go 1.16

require (
	github.com/alessio/shellescape v1.4.1
	github.com/antlr/antlr4 v0.0.0-20200225173536-225249fdaef5
	github.com/armon/circbuf v0.0.0-20190214190532-5111143e8da2
	github.com/containerd/containerd v1.6.6
	github.com/creack/pty v1.1.11
	github.com/docker/cli v20.10.14+incompatible
	github.com/docker/distribution v2.8.1+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/dustin/go-humanize v1.0.0
	github.com/earthly/cloud-api v1.0.1-0.20220728160826-b477d75c55f6
	github.com/elastic/go-sysinfo v1.7.1
	github.com/fatih/color v1.9.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/hashicorp/go-multierror v1.1.1
	github.com/jdxcode/netrc v0.0.0-20210204082910-926c7f70242a
	github.com/jessevdk/go-flags v1.5.0
	github.com/joho/godotenv v1.3.0
	github.com/mattn/go-colorable v0.1.8
	github.com/mattn/go-isatty v0.0.12
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/moby/buildkit v0.8.2-0.20210129065303-6b9ea0c202cf
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/image-spec v1.0.3-0.20211202183452-c5a74bcca799
	github.com/otiai10/copy v1.1.1
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/tonistiigi/fsutil v0.0.0-20220510150904-0dbf3a8a7d58
	github.com/urfave/cli/v2 v2.3.0
	golang.org/x/crypto v0.0.0-20220315160706-3147a52a75dd
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211
	google.golang.org/grpc v1.47.0
	gopkg.in/yaml.v3 v3.0.1
)

replace (
	github.com/docker/docker => github.com/docker/docker v20.10.3-0.20220414164044-61404de7df1a+incompatible
	github.com/jessevdk/go-flags => github.com/alexcb/go-flags v0.0.0-20210722203016-f11d7ecb5ee5

	github.com/moby/buildkit => github.com/earthly/buildkit v0.0.1-0.20220818002838-fd74555d900a
	github.com/tonistiigi/fsutil => github.com/earthly/fsutil v0.0.0-20220719234708-392da7eebeb5
)
