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
    SAVE IMAGE

code:
    FROM +deps
    COPY --dir autocomplete buildcontext builder cleanup cmd config conslogging debugger dockertar \
        domain llbutil logging ./
    COPY --dir buildkitd/buildkitd.go buildkitd/settings.go buildkitd/
    COPY --dir earthfile2llb/antlrhandler earthfile2llb/dedup earthfile2llb/image \
        earthfile2llb/imr earthfile2llb/variables earthfile2llb/*.go earthfile2llb/
    COPY ./earthfile2llb/parser+parser/*.go ./earthfile2llb/parser/
    SAVE IMAGE

lint-scripts:
    FROM +deps
    COPY ./earth ./buildkitd/entrypoint.sh ./earth-buildkitd-wrapper.sh \
        ./buildkitd/dockerd-wrapper.sh ./release/envcredhelper.sh \
        ./.buildkite/*.sh \
        ./shell_scripts/
    RUN shellcheck shell_scripts/*

lint:
    FROM +code
    RUN output="$(ineffassign . | grep -v '/earthly/earthfile2llb/parser/.*\.go')" ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi
    RUN output="$(goimports -d . 2>&1)" ; \
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
    ARG GO_EXTRA_LDFLAGS="-linkmode external -extldflags -static"
    RUN test -n "$GOOS" && test -n "$GOARCH"
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
    SAVE ARTIFACT build/earth AS LOCAL "build/$GOOS/$GOARCH/earth"

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
    SAVE ARTIFACT ./*

earth-docker:
    FROM ./buildkitd+buildkitd
    RUN apk add --update --no-cache docker-cli
    ENV ENABLE_LOOP_DEVICE=false
    ENV FORCE_LOOP_DEVICE=false
    COPY earth-buildkitd-wrapper.sh /usr/bin/earth-buildkitd-wrapper.sh
    ENTRYPOINT ["/usr/bin/earth-buildkitd-wrapper.sh"]
    ARG EARTHLY_TARGET_TAG_DOCKER
    ARG TAG=$EARTHLY_TARGET_TAG_DOCKER
    COPY --build-arg VERSION=$TAG +earth/earth /usr/bin/earth
    SAVE IMAGE --push earthly/earth:$TAG

# we abuse docker here to distribute our binaries
prerelease-docker:
    FROM alpine:3.11
    BUILD --build-arg TAG=prerelease ./buildkitd+buildkitd
    COPY --build-arg VERSION=prerelease +earth-all/* ./
    SAVE IMAGE --push earthly/earthlybinaries:prerelease

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
    BUILD +prerelease-docker

test:
    BUILD +lint
    BUILD +lint-scripts
    BUILD +unit-test
    BUILD ./examples/tests+ga

test-all:
    BUILD +examples
    BUILD +test
    BUILD ./examples/tests+experimental

examples:
    BUILD ./examples/go+docker
    BUILD ./examples/java+docker
    BUILD ./examples/js+docker
    BUILD ./examples/cpp+docker
    BUILD ./examples/scala+docker
    BUILD ./examples/dotnet+docker
    BUILD ./examples/python+docker
    BUILD ./examples/monorepo+all
    BUILD ./examples/multirepo+docker
    BUILD ./examples/integration-test+integration-test
    BUILD ./examples/readme/go1+all
    BUILD ./examples/readme/go2+all
    BUILD ./examples/readme/go3+build
    BUILD ./examples/readme/proto+docker
    BUILD github.com/earthly/hello-world+hello

test-fail:
    RUN false
