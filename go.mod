module github.com/earthly/earthly

go 1.13

require (
	github.com/alessio/shellescape v0.0.0-00010101000000-000000000000
	github.com/antlr/antlr4 v0.0.0-20200225173536-225249fdaef5
	github.com/armon/circbuf v0.0.0-20190214190532-5111143e8da2
	github.com/containerd/containerd v1.4.1-0.20200903181227-d4e78200d6da
	github.com/creack/pty v1.1.11
	github.com/docker/cli v20.10.0-beta1.0.20201029214301-1d20b15adc38+incompatible
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v20.10.0-beta1.0.20201110211921-af34b94a78a1+incompatible
	github.com/fatih/color v1.9.0
	github.com/golang/protobuf v1.4.2
	github.com/joho/godotenv v1.3.0
	github.com/moby/buildkit v0.7.1-0.20200708233707-488130002abb
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/image-spec v1.0.1
	github.com/otiai10/copy v1.1.1
	github.com/pkg/errors v0.9.1
	github.com/seehuhn/password v0.0.0-20131211191456-9ed6612376fa
	github.com/segmentio/backo-go v0.0.0-20200129164019-23eae7c10bd3 // indirect
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.5.1
	github.com/tonistiigi/fsutil v0.0.0-20201103201449-0834f99b7b85
	github.com/urfave/cli/v2 v2.1.1
	github.com/xtgo/uuid v0.0.0-20140804021211-a0b114877d4c // indirect
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/sync v0.0.0-20200625203802-6e8e738ad208
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.24.0
	gopkg.in/segmentio/analytics-go.v3 v3.1.0
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/alessio/shellescape => github.com/alexcb/shellescape v0.0.0-20200921195046-bf72418e9bfb
	github.com/docker/docker => github.com/docker/docker v17.12.0-ce-rc1.0.20200310163718-4634ce647cf2+incompatible
	github.com/hashicorp/go-immutable-radix => github.com/tonistiigi/go-immutable-radix v0.0.0-20170803185627-826af9ccf0fe
	github.com/jaguilar/vt100 => github.com/tonistiigi/vt100 v0.0.0-20190402012908-ad4c4a574305
	github.com/moby/buildkit => github.com/earthly/buildkit v0.7.1-0.20201117194031-9d476009bb3b
	github.com/urfave/cli/v2 => github.com/alexcb/cli/v2 v2.2.1-0.20200824212017-2ae03fa69ce7
)
