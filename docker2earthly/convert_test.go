package docker2earthly_test

import (
	"errors"
	"testing"

	"github.com/earthly/earthly/docker2earthly"
)

func TestGenerateEarthfile(t *testing.T) {
	type args struct {
		buildContextPath string
		dockerfilePath   string
		imageTags        []string
		buildArgs        []string
		platforms        []string
		target           string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr error
	}{
		{
			name: "all fields are populated",
			args: args{
				buildContextPath: "/my/build/context",
				dockerfilePath:   "./dir/../MyDockerfile",
				imageTags:        []string{"test-image:v1.2.3", "test-image:v1.2.3.4"},
				buildArgs:        []string{"arg1", "arg2"},
				platforms:        []string{"linux/amd64", "linux/arm64"},
				target:           "target1",
			},
			want: `
VERSION --use-docker-ignore 0.7
# This Earthfile was generated using docker-build command
docker:
	ARG arg1
	ARG arg2
	FROM DOCKERFILE \
	--build-arg arg1=$arg1 \
	--build-arg arg2=$arg2 \
	--target target1 \
	-f /my/build/context/MyDockerfile \
	/my/build/context
	SAVE IMAGE --push test-image:v1.2.3 test-image:v1.2.3.4

build:
	BUILD --platform linux/amd64 --platform linux/arm64 +docker
`,
		},
		{
			name: "Dockerfile has absolute path",
			args: args{
				buildContextPath: "/my/build/context",
				dockerfilePath:   "/my/build/context/dir/MyDockerfile",
				imageTags:        []string{"test-image:v1.2.3"},
				buildArgs:        []string{"arg1", "arg2"},
				platforms:        []string{"linux/amd64", "linux/arm64"},
				target:           "target1",
			},
			want: `
VERSION --use-docker-ignore 0.7
# This Earthfile was generated using docker-build command
docker:
	ARG arg1
	ARG arg2
	FROM DOCKERFILE \
	--build-arg arg1=$arg1 \
	--build-arg arg2=$arg2 \
	--target target1 \
	-f /my/build/context/dir/MyDockerfile \
	/my/build/context
	SAVE IMAGE --push test-image:v1.2.3

build:
	BUILD --platform linux/amd64 --platform linux/arm64 +docker
`,
		},
		{
			name: "no args",
			args: args{
				buildContextPath: "/build-context",
				dockerfilePath:   "./dir/MyDockerfile",
				imageTags:        []string{"test-image:v1.2.3"},
				platforms:        []string{"linux/amd64"},
				target:           "target1",
			},
			want: `
VERSION --use-docker-ignore 0.7
# This Earthfile was generated using docker-build command
docker:
	FROM DOCKERFILE \
	--target target1 \
	-f /build-context/dir/MyDockerfile \
	/build-context
	SAVE IMAGE --push test-image:v1.2.3

build:
	BUILD --platform linux/amd64 +docker
`,
		},
		{
			name: "no target",
			args: args{
				buildContextPath: "/build-context",
				dockerfilePath:   "./dir/MyDockerfile",
				imageTags:        []string{"test-image:v1.2.3"},
				buildArgs:        []string{"arg1"},
				platforms:        []string{"linux/amd64"},
			},
			want: `
VERSION --use-docker-ignore 0.7
# This Earthfile was generated using docker-build command
docker:
	ARG arg1
	FROM DOCKERFILE \
	--build-arg arg1=$arg1 \
	-f /build-context/dir/MyDockerfile \
	/build-context
	SAVE IMAGE --push test-image:v1.2.3

build:
	BUILD --platform linux/amd64 +docker
`,
		},
		{
			name: "no platform",
			args: args{
				buildContextPath: "/build-context",
				dockerfilePath:   "./dir/MyDockerfile",
				imageTags:        []string{"test-image:v1.2.3"},
				buildArgs:        []string{"arg1"},
				target:           "target1",
			},
			want: `
VERSION --use-docker-ignore 0.7
# This Earthfile was generated using docker-build command
docker:
	ARG arg1
	FROM DOCKERFILE \
	--build-arg arg1=$arg1 \
	--target target1 \
	-f /build-context/dir/MyDockerfile \
	/build-context
	SAVE IMAGE --push test-image:v1.2.3

build:
	BUILD +docker
`,
		},
		{
			name: "no tags",
			args: args{
				buildContextPath: "/my/build/context",
				dockerfilePath:   "./dir/../MyDockerfile",
				buildArgs:        []string{"arg1", "arg2"},
				platforms:        []string{"linux/amd64", "linux/arm64"},
				target:           "target1",
			},
			want: `
VERSION --use-docker-ignore 0.7
# This Earthfile was generated using docker-build command
docker:
	ARG arg1
	ARG arg2
	FROM DOCKERFILE \
	--build-arg arg1=$arg1 \
	--build-arg arg2=$arg2 \
	--target target1 \
	-f /my/build/context/MyDockerfile \
	/my/build/context

build:
	BUILD --platform linux/amd64 --platform linux/arm64 +docker
`,
		},
		{
			name: "no optional values",
			args: args{
				buildContextPath: "/build-context",
				dockerfilePath:   "./dir/MyDockerfile",
			},
			want: `
VERSION --use-docker-ignore 0.7
# This Earthfile was generated using docker-build command
docker:
	FROM DOCKERFILE \
	-f /build-context/dir/MyDockerfile \
	/build-context

build:
	BUILD +docker
`,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := docker2earthly.GenerateEarthfile(tt.args.buildContextPath, tt.args.dockerfilePath, tt.args.imageTags, tt.args.buildArgs, tt.args.platforms, tt.args.target)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("GenerateEarthfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateEarthfile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
