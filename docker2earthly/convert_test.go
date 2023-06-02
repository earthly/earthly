package docker2earthly

import (
	"testing"
)

func TestGenerateEarthfileContent(t *testing.T) {
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
		wantErr bool
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
			wantErr: false,
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
			wantErr: false,
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
			wantErr: false,
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
			wantErr: false,
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
			wantErr: false,
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
			wantErr: false,
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateEarthfileContent(tt.args.buildContextPath, tt.args.dockerfilePath, tt.args.imageTags, tt.args.buildArgs, tt.args.platforms, tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateEarthfileContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateEarthfileContent() got = %v, want %v", got, tt.want)
			}
		})
	}
}
