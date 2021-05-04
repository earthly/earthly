FROM golang:1.16-alpine3.13

RUN apk add --update --no-cache \
    bash \
    bash-completion \
    binutils \
    ca-certificates \
    coreutils \
    curl \
    findutils \
    g++ \
    git \
    grep \
    less \
    make \
    openssl \
    util-linux

WORKDIR /earthly

deps:
    FROM +base
    RUN go get golang.org/x/tools/cmd/goimports
    RUN go get golang.org/x/lint/golint
    RUN go get github.com/gordonklaus/ineffassign
    COPY go.mod go.sum ./
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

code:
    FROM +deps
    COPY --platform=linux/amd64 ./ast/parser+parser/*.go ./ast/parser/
    COPY --dir analytics autocomplete buildcontext builder cleanup cmd config conslogging debugger dockertar \
        docker2earthly domain logging secretsclient states util variables ./
    COPY --dir buildkitd/buildkitd.go buildkitd/settings.go buildkitd/
    COPY --dir earthfile2llb/*.go earthfile2llb/
    COPY --dir ast/antlrhandler ast/spec ast/*.go ast/

update-buildkit:
    FROM +code # if we use deps, go mod tidy will remove a bunch of requirements since it won't have access to our codebase.
    ARG BUILDKIT_BRANCH=earthly-main
    BUILD ./buildkitd+update-buildkit --BUILDKIT_BRANCH=$BUILDKIT_BRANCH
    RUN --no-cache go mod edit -replace "github.com/moby/buildkit=github.com/earthly/buildkit@$BUILDKIT_BRANCH"
    RUN --no-cache go mod tidy
    SAVE ARTIFACT go.mod AS LOCAL go.mod-fixme  # this is a bug since we can't save to go.mod which was already saved in +deps
    SAVE ARTIFACT go.sum AS LOCAL go.sum-fixme  # this is a bug since we can't save to go.sum which was already saved in +deps


lint-scripts:
    FROM --platform=linux/amd64 alpine:3.13
    RUN apk add --update --no-cache shellcheck
    COPY ./earthly ./scripts/install-all-versions.sh ./buildkitd/entrypoint.sh ./earthly-buildkitd-wrapper.sh \
        ./buildkitd/dockerd-wrapper.sh ./buildkitd/docker-auto-install.sh \
        ./release/envcredhelper.sh ./.buildkite/*.sh \
        ./scripts/tests/*.sh \
        ./shell_scripts/
    RUN shellcheck shell_scripts/*.sh

lint:
    FROM +code
    RUN output="$(ineffassign ./... 2>&1 | grep -v '/earthly/ast/parser/.*\.go')" ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi
    RUN output="$(goimports -d $(find . -type f -name '*.go' | grep -v \.pb\.go) 2>&1)"  ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi
    RUN golint -set_exit_status ./...
    RUN output="$(go vet ./... 2>&1)" ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi

unit-test:
    FROM +code
    RUN go test ./...

shellrepeater:
    FROM +code
    ARG GOCACHE=/go-cache
    ARG EARTHLY_TARGET_TAG
    ARG VERSION=$EARTHLY_TARGET_TAG
    ARG EARTHLY_GIT_HASH
    RUN --mount=type=cache,target=$GOCACHE \
        go build \
            -ldflags "-d -X main.Version=$VERSION $GO_EXTRA_LDFLAGS -X main.GitSha=$EARTHLY_GIT_HASH $GO_EXTRA_LDFLAGS" \
            -tags netgo -installsuffix netgo \
            -o build/shellrepeater \
            cmd/shellrepeater/*.go
    SAVE ARTIFACT build/shellrepeater

debugger:
    FROM +code
    ARG GOCACHE=/go-cache
    ARG EARTHLY_TARGET_TAG
    ARG VERSION=$EARTHLY_TARGET_TAG
    ARG EARTHLY_GIT_HASH
    RUN --mount=type=cache,target=$GOCACHE \
        go build \
            -ldflags "-d -X main.Version=$VERSION $GO_EXTRA_LDFLAGS -X main.GitSha=$EARTHLY_GIT_HASH $GO_EXTRA_LDFLAGS" \
            -tags netgo -installsuffix netgo \
            -o build/earth_debugger \
            cmd/debugger/*.go
    SAVE ARTIFACT build/earth_debugger

earthly:
    FROM +code
    ARG GOOS=linux
    ARG TARGETARCH
    ARG TARGETVARIANT
    ARG GOARCH=$TARGETARCH
    ARG VARIANT=$TARGETVARIANT
    ARG GO_EXTRA_LDFLAGS="-linkmode external -extldflags -static"
    RUN test -n "$GOOS" && test -n "$GOARCH"
    RUN test "$GOARCH" != "arm" || test -n "$VARIANT"
    ARG EARTHLY_TARGET_TAG_DOCKER
    ARG VERSION=$EARTHLY_TARGET_TAG_DOCKER
    ARG EARTHLY_GIT_HASH
    ARG DEFAULT_BUILDKITD_IMAGE=earthly/buildkitd:$VERSION
    ARG BUILD_TAGS=dfrunmount dfrunsecurity dfsecrets dfssh dfrunnetwork
    ARG GOCACHE=/go-cache
    RUN mkdir -p build
    RUN printf "$BUILD_TAGS" > ./build/tags && echo "$(cat ./build/tags)"
    RUN printf '-X main.DefaultBuildkitdImage='"$DEFAULT_BUILDKITD_IMAGE" > ./build/ldflags && \
        printf ' -X main.Version='"$VERSION" >> ./build/ldflags && \
        printf ' -X main.GitSha='"$EARTHLY_GIT_HASH" >> ./build/ldflags && \
        printf ' '"$GO_EXTRA_LDFLAGS" >> ./build/ldflags && \
        echo "$(cat ./build/ldflags)"
    # Important! If you change the go build options, you may need to also change them
    # in https://github.com/Homebrew/homebrew-core/blob/master/Formula/earthly.rb.
    RUN --mount=type=cache,target=$GOCACHE \
        GOARM=${VARIANT#v} go build \
            -tags "$(cat ./build/tags)" \
            -ldflags "$(cat ./build/ldflags)" \
            -o build/earthly \
            cmd/earthly/*.go
    SAVE ARTIFACT ./build/tags
    SAVE ARTIFACT ./build/ldflags
    SAVE ARTIFACT build/earthly AS LOCAL "build/$GOOS/$GOARCH$VARIANT/earthly"
    SAVE IMAGE --cache-from=earthly/earthly:main

earthly-linux-amd64:
    COPY \
        --build-arg GOARCH=amd64 \
        --build-arg VARIANT= \
        +earthly/* ./
    SAVE ARTIFACT ./*

earthly-linux-arm7:
    COPY \
        --build-arg GOARCH=arm \
        --build-arg VARIANT=v7 \
        --build-arg GO_EXTRA_LDFLAGS= \
        +earthly/* ./
    SAVE ARTIFACT ./*

earthly-linux-arm64:
    COPY \
        --build-arg GOARCH=arm64 \
        --build-arg VARIANT= \
        --build-arg GO_EXTRA_LDFLAGS= \
        +earthly/* ./
    SAVE ARTIFACT ./*

earthly-darwin-amd64:
    COPY \
        --build-arg GOOS=darwin \
        --build-arg GOARCH=amd64 \
        --build-arg VARIANT= \
        --build-arg GO_EXTRA_LDFLAGS= \
        +earthly/* ./
    SAVE ARTIFACT ./*

earthly-darwin-arm64:
    COPY \
        --build-arg GOOS=darwin \
        --build-arg GOARCH=arm64 \
        --build-arg VARIANT= \
        --build-arg GO_EXTRA_LDFLAGS= \
        +earthly/* ./
    SAVE ARTIFACT ./*

earthly-all:
    COPY +earthly-linux-amd64/earthly ./earthly-linux-amd64
    COPY +earthly-linux-arm7/earthly ./earthly-linux-arm7
    COPY +earthly-linux-arm64/earthly ./earthly-linux-arm64
    COPY +earthly-darwin-amd64/earthly ./earthly-darwin-amd64
    COPY +earthly-darwin-arm64/earthly ./earthly-darwin-arm64
    SAVE ARTIFACT ./*

earthly-docker:
    FROM ./buildkitd+buildkitd
    RUN apk add --update --no-cache docker-cli
    ENV NETWORK_MODE=host
    ENV EARTHLY_IMAGE=true
    COPY earthly-buildkitd-wrapper.sh /usr/bin/earthly-buildkitd-wrapper.sh
    ENTRYPOINT ["/usr/bin/earthly-buildkitd-wrapper.sh"]
    ARG EARTHLY_TARGET_TAG_DOCKER
    ARG TAG=$EARTHLY_TARGET_TAG_DOCKER
    COPY --build-arg VERSION=$TAG +earthly/earthly /usr/bin/earthly
    SAVE IMAGE --push --cache-from=earthly/earthly:main earthly/earthly:$TAG

earthly-integration-test-base:
    FROM +earthly-docker
    ENV EARTHLY_CONVERSION_PARALLELISM=5
    RUN earthly config global.disable_analytics true
    # The inner buildkit requires Docker hub creds to prevent rate-limiting issues.
    ARG DOCKERHUB_AUTH=true
    ARG DOCKERHUB_USER_SECRET=+secrets/earthly-technologies/dockerhub/user
    ARG DOCKERHUB_TOKEN_SECRET=+secrets/earthly-technologies/dockerhub/token
    IF $DOCKERHUB_AUTH
        RUN --secret USERNAME=$DOCKERHUB_USER_SECRET \
            --secret TOKEN=$DOCKERHUB_TOKEN_SECRET \
            docker login --username="$USERNAME" --password="$TOKEN"
    END

prerelease:
    FROM alpine:3.13
    BUILD --build-arg TAG=prerelease \
        --platform=linux/amd64 \
        --platform=linux/arm/v7 \
        --platform=linux/arm64 \
        ./buildkitd+buildkitd
    COPY --build-arg VERSION=prerelease +earthly-all/* ./
    SAVE IMAGE --push earthly/earthlybinaries:prerelease

prerelease-script:
    FROM alpine:3.13
    COPY ./earthly ./
    # This script is useful in other repos too.
    SAVE ARTIFACT ./earthly

dind:
    BUILD +dind-alpine
    BUILD +dind-ubuntu

dind-alpine:
    FROM docker:dind
    RUN apk add --update --no-cache docker-compose
    ARG EARTHLY_TARGET_TAG_DOCKER
    ARG DIND_ALPINE_TAG=alpine-$EARTHLY_TARGET_TAG_DOCKER
    SAVE IMAGE --push --cache-from=earthly/dind:main earthly/dind:$DIND_ALPINE_TAG

dind-ubuntu:
    FROM ubuntu:latest
    COPY ./buildkitd/docker-auto-install.sh /usr/local/bin/docker-auto-install.sh
    RUN docker-auto-install.sh
    ARG DIND_UBUNTU_TAG=ubuntu-$EARTHLY_TARGET_TAG_DOCKER
    SAVE IMAGE --push --cache-from=earthly/dind:ubuntu-main earthly/dind:$DIND_UBUNTU_TAG

for-own:
    BUILD ./buildkitd+buildkitd
    COPY +earthly/earthly ./
    SAVE ARTIFACT ./earthly

for-linux:
    BUILD --platform=linux/amd64 ./buildkitd+buildkitd
    COPY +earthly-linux-amd64/earthly ./
    SAVE ARTIFACT ./earthly

for-darwin:
    BUILD --platform=linux/amd64 ./buildkitd+buildkitd
    COPY +earthly-darwin-amd64/earthly ./
    SAVE ARTIFACT ./earthly

for-darwin-m1:
    BUILD --platform=linux/arm64 ./buildkitd+buildkitd
    COPY +earthly-darwin-arm64/earthly ./
    SAVE ARTIFACT ./earthly

all-buildkitd:
    BUILD \
        --platform=linux/amd64 \
        --platform=linux/arm/v7 \
        --platform=linux/arm64 \
        ./buildkitd+buildkitd

all-dind:
    BUILD \
        --platform=linux/amd64 \
        --platform=linux/arm64 \
        +dind
    BUILD \
        --platform=linux/arm/v7 \
        +dind-alpine

all:
    BUILD +all-buildkitd
    BUILD +earthly-all
    BUILD +earthly-docker
    BUILD +prerelease
    BUILD +all-dind

test:
    BUILD +lint
    BUILD +lint-scripts
    BUILD +unit-test
    BUILD ./ast/tests+all
    ARG DOCKERHUB_AUTH=true
    BUILD ./examples/tests+ga --DOCKERHUB_AUTH=$DOCKERHUB_AUTH

test-all:
    BUILD +examples
    BUILD +test
    ARG DOCKERHUB_AUTH=true
    BUILD ./examples/tests+experimental --DOCKERHUB_AUTH=$DOCKERHUB_AUTH

examples:
    BUILD +examples1
    BUILD +examples2

examples1:
    ARG TARGETARCH
    BUILD ./examples/cpp+docker
    IF [ "$TARGETARCH" = "amd64" ]
        # This only works on amd64 for now.
        BUILD ./examples/dotnet+docker
    END
    BUILD ./examples/elixir+docker
    BUILD ./examples/go+docker
    BUILD ./examples/grpc+test
    IF [ "$TARGETARCH" = "amd64" ]
        # This only works on amd64 for now.
        BUILD ./examples/integration-test+integration-test
    END
    BUILD ./examples/java+docker
    BUILD ./examples/js+docker
    BUILD ./examples/monorepo+all
    BUILD ./examples/multirepo+docker
    BUILD ./examples/python+docker
    BUILD ./examples/cutoff-optimization+run

examples2:
    BUILD ./examples/readme/go1+all
    BUILD ./examples/readme/go2+build
    BUILD ./examples/readme/proto+docker
    # TODO: This example is flaky for some reason.
    #BUILD ./examples/terraform+localstack
    BUILD ./examples/ruby+docker
    BUILD ./examples/ruby-on-rails+docker
    IF [ "$TARGETARCH" = "amd64" ]
        # This crashes randomly on arm.
        BUILD ./examples/scala+docker
    END
    BUILD ./examples/cobol+docker
    BUILD ./examples/rust+docker
    BUILD ./examples/multiplatform+all
    BUILD ./examples/multiplatform-cross-compile+build-all-platforms
    BUILD github.com/earthly/hello-world:main+hello
