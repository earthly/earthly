FROM golang:1.13-alpine3.11

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
    shellcheck \
    util-linux

WORKDIR /earthly

deps:
    RUN go get golang.org/x/tools/cmd/goimports
    RUN go get golang.org/x/lint/golint
    RUN go get github.com/gordonklaus/ineffassign
    COPY go.mod go.sum ./
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

code:
    FROM +deps
    COPY --dir analytics autocomplete buildcontext builder cleanup cmd config conslogging debugger dockertar \
        docker2earth domain fileutils llbutil logging secretsclient states ./
    COPY --dir buildkitd/buildkitd.go buildkitd/settings.go buildkitd/
    COPY --dir earthfile2llb/antlrhandler \
        earthfile2llb/variables earthfile2llb/*.go earthfile2llb/
    COPY ./earthfile2llb/parser+parser/*.go ./earthfile2llb/parser/

lint-scripts:
    FROM +deps
    COPY ./earth ./buildkitd/entrypoint.sh ./earth-buildkitd-wrapper.sh \
        ./buildkitd/dockerd-wrapper.sh ./buildkitd/docker-auto-install.sh \
        ./release/envcredhelper.sh ./.buildkite/*.sh \
        ./.tests/private-repo.sh \
        ./shell_scripts/
    RUN shellcheck shell_scripts/*

lint:
    FROM +code
    RUN output="$(ineffassign . | grep -v '/earthly/earthfile2llb/parser/.*\.go')" ; \
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

buildkitd:
    BUILD ./buildkitd+buildkitd

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

earth:
    FROM +code
    ARG GOOS=linux
    ARG GOARCH=amd64
    ARG GOARM
    ARG GO_EXTRA_LDFLAGS="-linkmode external -extldflags -static"
    RUN test -n "$GOOS" && test -n "$GOARCH"
    RUN test "$GOARCH" != "ARM" || test -n "$GOARM"
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
        go build \
            -tags "$(cat ./build/tags)" \
            -ldflags "$(cat ./build/ldflags)" \
            -o build/earth \
            cmd/earth/*.go
    SAVE ARTIFACT ./build/tags
    SAVE ARTIFACT ./build/ldflags
    SAVE ARTIFACT build/earth AS LOCAL "build/$GOOS/$GOARCH$GOARM/earth"

earth-arm5:
    COPY \
        --build-arg GOARCH=arm \
        --build-arg GOARM=5 \
        --build-arg GO_EXTRA_LDFLAGS= \
        +earth/* ./
    SAVE ARTIFACT ./*

earth-arm6:
    COPY \
        --build-arg GOARCH=arm \
        --build-arg GOARM=6 \
        --build-arg GO_EXTRA_LDFLAGS= \
        +earth/* ./
    SAVE ARTIFACT ./*

earth-arm7:
    COPY \
        --build-arg GOARCH=arm \
        --build-arg GOARM=7 \
        --build-arg GO_EXTRA_LDFLAGS= \
        +earth/* ./
    SAVE ARTIFACT ./*

earth-darwin:
    COPY \
        --build-arg GOOS=darwin \
        --build-arg GOARCH=amd64 \
        --build-arg GO_EXTRA_LDFLAGS= \
        +earth/* ./
    SAVE ARTIFACT ./*

earth-all:
    COPY +earth/earth ./earth-linux-amd64
    COPY +earth-darwin/earth ./earth-darwin-amd64
    COPY +earth-arm5/earth ./earth-linux-arm5
    COPY +earth-arm6/earth ./earth-linux-arm6
    COPY +earth-arm7/earth ./earth-linux-arm7
    SAVE ARTIFACT ./*

earth-docker:
    FROM ./buildkitd+buildkitd
    RUN apk add --update --no-cache docker-cli
    ENV NETWORK_MODE=host
    COPY earth-buildkitd-wrapper.sh /usr/bin/earth-buildkitd-wrapper.sh
    ENTRYPOINT ["/usr/bin/earth-buildkitd-wrapper.sh"]
    ARG EARTHLY_TARGET_TAG_DOCKER
    ARG TAG=$EARTHLY_TARGET_TAG_DOCKER
    COPY --build-arg VERSION=$TAG +earth/earth /usr/bin/earth
    SAVE IMAGE --push earthly/earth:$TAG

prerelease:
    FROM alpine:3.11
    BUILD --build-arg TAG=prerelease ./buildkitd+buildkitd
    COPY --build-arg VERSION=prerelease +earth-all/* ./
    SAVE IMAGE --push earthly/earthlybinaries:prerelease

dind:
    BUILD +dind-alpine
    BUILD +dind-ubuntu

dind-alpine:
    FROM docker:dind
    RUN apk add --update --no-cache docker-compose
    ARG EARTHLY_TARGET_TAG_DOCKER
    ARG DIND_ALPINE_TAG=alpine-$EARTHLY_TARGET_TAG_DOCKER
    SAVE IMAGE --push earthly/dind:$DIND_ALPINE_TAG

dind-ubuntu:
    FROM ubuntu:latest
    COPY ./buildkitd/docker-auto-install.sh /usr/local/bin/docker-auto-install.sh
    RUN docker-auto-install.sh
    ARG DIND_UBUNTU_TAG=ubuntu-$EARTHLY_TARGET_TAG_DOCKER
    SAVE IMAGE --push earthly/dind:$DIND_UBUNTU_TAG

for-linux:
    BUILD +buildkitd
    COPY +earth/earth ./
    SAVE ARTIFACT ./earth

for-darwin:
    BUILD +buildkitd
    COPY +earth-darwin/earth ./
    SAVE ARTIFACT ./earth

all:
    BUILD +buildkitd
    BUILD +earth-all
    BUILD +earth-docker
    BUILD +prerelease
    BUILD +dind

test:
    BUILD +lint
    BUILD +lint-scripts
    BUILD +unit-test
    ARG DOCKERHUB_USER_SECRET
    ARG DOCKERHUB_TOKEN_SECRET
    BUILD \
        --build-arg DOCKERHUB_USER_SECRET \
        --build-arg DOCKERHUB_TOKEN_SECRET \
        ./examples/tests+ga

test-all:
    BUILD +examples
    BUILD +test
    ARG DOCKERHUB_USER_SECRET
    ARG DOCKERHUB_TOKEN_SECRET
    BUILD \
        --build-arg DOCKERHUB_USER_SECRET \
        --build-arg DOCKERHUB_TOKEN_SECRET \
        ./examples/tests+experimental

examples:
    BUILD ./examples/cpp+docker
    BUILD ./examples/dotnet+docker
    BUILD ./examples/elixir+docker
    BUILD ./examples/go+docker
    BUILD ./examples/grpc+test
    BUILD ./examples/integration-test+integration-test
    BUILD ./examples/java+docker
    BUILD ./examples/js+docker
    BUILD ./examples/monorepo+all
    BUILD ./examples/multirepo+docker
    BUILD ./examples/python+docker
    BUILD ./examples/readme/go1+all
    BUILD ./examples/readme/go2+all
    BUILD ./examples/readme/go3+build
    BUILD ./examples/readme/proto+docker
    # TODO: This example is flaky for some reason.
    #BUILD ./examples/terraform+localstack
    BUILD ./examples/ruby+docker
    BUILD ./examples/ruby-on-rails+docker
    BUILD ./examples/scala+docker
    BUILD github.com/earthly/hello-world:main+hello

test-fail:
    RUN false

test-interactive:
    FROM earthly/dind:alpine
    RUN apk add --update --no-cache \
        jq
    RUN cat /run/secrets/earthly_debugger_settings | jq -r .repeaterAddr | xargs nc
