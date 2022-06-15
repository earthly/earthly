VERSION --shell-out-anywhere --use-copy-link 0.6

FROM golang:1.17-alpine3.14

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
    RUN go install golang.org/x/tools/cmd/goimports@latest
    RUN go install golang.org/x/lint/golint@latest
    RUN go install github.com/gordonklaus/ineffassign@latest
    COPY go.mod go.sum ./
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

code:
    FROM +deps
    # Use BUILDKIT_PROJECT to point go.mod to a buildkit dir being actively developed.
    # --BUILDKIT_PROJECT=../buildkit or --BUILDKIT_PROJECT=github.com/earthly/buildkit:4f1c968c3778a7140c56e3e6a755b7f2d38f2156
    ARG BUILDKIT_PROJECT
    IF [ "$BUILDKIT_PROJECT" != "" ]
        COPY --dir "$BUILDKIT_PROJECT"+code/buildkit /buildkit
        RUN go mod edit -replace github.com/moby/buildkit=/buildkit
        RUN go mod download
    END
    COPY ./ast/parser+parser/*.go ./ast/parser/
    COPY --dir analytics autocomplete buildcontext builder cleanup cmd config conslogging debugger dockertar \
        docker2earthly domain features outmon slog cloud states util variables ./
    COPY --dir buildkitd/buildkitd.go buildkitd/settings.go buildkitd/certificates.go buildkitd/
    COPY --dir earthfile2llb/*.go earthfile2llb/
    COPY --dir ast/antlrhandler ast/spec ast/*.go ast/

update-buildkit:
    FROM +code # if we use deps, go mod tidy will remove a bunch of requirements since it won't have access to our codebase.
    ARG BUILDKIT_GIT_SHA
    ARG BUILDKIT_GIT_BRANCH=earthly-main
    COPY (./buildkitd+buildkit-sha/buildkit_sha --BUILDKIT_GIT_SHA="$BUILDKIT_GIT_SHA" --BUILDKIT_GIT_BRANCH="$BUILDKIT_GIT_BRANCH") buildkit_sha
    BUILD  ./buildkitd+update-buildkit --BUILDKIT_GIT_SHA="$(cat buildkit_sha)"
    RUN --no-cache go mod edit -replace "github.com/moby/buildkit=github.com/earthly/buildkit@$(cat buildkit_sha)"
    RUN --no-cache go mod tidy
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

lint-scripts-base:
    FROM alpine:3.15

    ARG TARGETARCH

    IF [ $TARGETARCH == "arm64" ]
        RUN echo "Downloading, and manually installing shellcheck for ARM" && \
            wget https://github.com/koalaman/shellcheck/releases/download/stable/shellcheck-stable.linux.aarch64.tar.xz && \
            tar -xf shellcheck-stable.linux.aarch64.tar.xz && \
            mv shellcheck-stable/shellcheck /usr/bin/shellcheck
    ELSE
        RUN echo "Installing shellcheck from Alpine repos" && \
            apk add --update --no-cache shellcheck
    END

    WORKDIR /shell_scripts

lint-scripts-misc:
    FROM +lint-scripts-base
    COPY ./earthly ./scripts/install-all-versions.sh ./buildkitd/entrypoint.sh ./earthly-entrypoint.sh \
        ./buildkitd/dockerd-wrapper.sh ./buildkitd/docker-auto-install.sh \
        ./release/envcredhelper.sh ./.buildkite/*.sh \
        ./scripts/tests/*.sh \
        ./scripts/*.sh \
        ./shell_scripts/
    RUN shellcheck shell_scripts/*

lint-scripts-auth-test:
    FROM +lint-scripts-base
    COPY ./scripts/tests/auth/*.sh ./
    # the auth test script make use of a common setup.sh which contain unused variables
    # when run directly; so we must exclude checking this directly, and make use of the -x
    # flag to source setup.sh during analysis.
    RUN shellcheck -x test-*.sh

lint-scripts:
    BUILD +lint-scripts-auth-test
    BUILD +lint-scripts-misc

earthly-script-no-stdout:
    # This validates the ./earthly script doesn't print anything to stdout (it should print to stderr)
    # This is to ensure commands such as: MYSECRET="$(./earthly secrets get -n /user/my-secret)" work
    FROM earthly/dind:alpine
    RUN apk add --no-cache --update bash
    COPY earthly .earthly_version_flag_overrides .

    # This script performs an explicit "docker pull earthlybinaries:prerelease" which can cause rate-limiting
    # to work-around this, we will copy an earthly binary in, and disable auto-updating (and therefore don't require a WITH DOCKER)
    COPY +earthly/earthly /root/.earthly/earthly-prerelease
    RUN EARTHLY_DISABLE_AUTO_UPDATE=true ./earthly --version > earthly-version-output

    RUN test "$(cat earthly-version-output | wc -l)" = "1"
    RUN grep '^earthly version.*$' earthly-version-output # only --version info should go to stdout

lint:
    FROM +code
    RUN output="$(ineffassign ./... 2>&1 | grep -v '/earthly/ast/parser/.*\.go')" ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi
    RUN output="$(goimports -d $(find . -type f -name '*.go' | grep -v \./ast/parser/.*\.go) 2>&1)"  ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi
    # names.go defines some very obvious consts that do not need comments; however golint doesn't support disabling rules: https://github.com/golang/lint/issues/263
    # therefore, we will hide this file from golint, and restore it after.
    RUN mv variables/reserved/names.go variables/reserved/names.skip-go-lint && echo "package reserved" > variables/reserved/names.go
    RUN golint -set_exit_status ./...
    RUN mv variables/reserved/names.skip-go-lint variables/reserved/names.go
    RUN output="$(go vet ./... 2>&1)" ; \
        if [ -n "$output" ]; then \
            echo "$output" ; \
            exit 1 ; \
        fi
    RUN if find . -type f -name \*.go | xargs grep '"io/ioutil"'; then echo "io/ioutil is deprecated: https://go.dev/doc/go1.16#ioutil"; exit 1; fi

lint-newline-ending:
    FROM alpine:3.15
    WORKDIR /everything
    COPY . .
    # test that line endings are unix-style
    RUN set -e; \
        code=0; \
        for f in $(find . -type f \( -iname '*.go' -o -iname 'Earthfile' -o -iname '*.earth' -o -iname '*.md' \) | grep -v "ast/tests/empty-targets.earth" ); do \
            if ! dos2unix < "$f" | cmp - "$f"; then \
                echo "$f contains windows-style newlines and must be converted to unix-style (use dos2unix to fix)"; \
                code=1; \
            fi; \
        done; \
        exit $code
    # test file ends with a single newline
    RUN set -e; \
        code=0; \
        for f in $(find . -type f \( -iname '*.yml' -o -iname '*.go' -o -iname '*.sh' -o -iname '*.template' -o -iname 'Earthfile' -o -iname '*.earth' -o -iname '*.md' \) | grep -v "ast/tests/empty-targets.earth" | grep -v "tests/version/version-only.earth" ); do \
            if [ "$(tail -c 1 $f)" != "$(printf '\n')" ]; then \
                echo "$f does not end with a newline"; \
                code=1; \
            fi; \
        done; \
        exit $code
    RUN if [ "$(tail -c 1 ast/tests/empty-targets.earth)" = "$(printf '\n')" ]; then \
            echo "$f is a special-case test which must not end with a newline."; \
            exit 1; \
        fi
    # check for files with trailing newlines
    RUN set -e; \
        code=0; \
        for f in $(find . -type f \( -iname '*.go' -o -iname 'Earthfile' -o -iname '*.earth' -o -iname '*.md' \) | grep -v "ast/tests/empty-targets.earth" | grep -v "ast/parser/earth_parser.go" | grep -v "ast/parser/earth_lexer.go" ); do \
            if [ "$(tail -c 2 $f)" == "$(printf '\n\n')" ]; then \
                echo "$f has trailing newlines"; \
                code=1; \
            fi; \
        done; \
        exit $code

vale:
    WORKDIR /
    RUN curl -sfL https://install.goreleaser.com/github.com/ValeLint/vale.sh | sh -s v2.10.3
    WORKDIR /etc/vale
    COPY .vale/ .

markdown-spellcheck:
    FROM --platform=linux/amd64 +vale
    WORKDIR /everything
    COPY . .
    # TODO figure out a way to ignore this pattern in vale (doesn't seem to be working under spelling's filter option)
    RUN find . -type f -iname '*.md' |  xargs -n 1 sed -i 's/{[^}]*}//g'
    # TODO remove the greps once the corresponding markdown files have spelling fixed (or techterms added to .vale/styles/HouseStyle/tech-terms/...
    RUN find . -type f -iname '*.md' | xargs vale --config /etc/vale/vale.ini --output line --minAlertLevel error

unit-test:
    FROM +code
    COPY podman-setup.sh .
    WITH DOCKER
        RUN ./podman-setup.sh && \
            go test ./...
    END

changelog:
    FROM scratch
    COPY CHANGELOG.md .
    SAVE ARTIFACT CHANGELOG.md

lint-changelog:
    FROM python:3
    COPY release/changelogparser.py /usr/bin/changelogparser
    COPY CHANGELOG.md .
    RUN changelogparser --changelog CHANGELOG.md

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
    ARG GO_EXTRA_LDFLAGS="-linkmode external -extldflags -static"
    ARG EARTHLY_TARGET_TAG
    ARG VERSION=$EARTHLY_TARGET_TAG
    ARG EARTHLY_GIT_HASH
    RUN --mount=type=cache,target=$GOCACHE \
        go build \
            -ldflags "-X main.Version=$VERSION -X main.GitSha=$EARTHLY_GIT_HASH $GO_EXTRA_LDFLAGS" \
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
    ARG EXECUTABLE_NAME="earthly"
    RUN test -n "$GOOS" && test -n "$GOARCH"
    RUN test "$GOARCH" != "arm" || test -n "$VARIANT"
    ARG EARTHLY_TARGET_TAG_DOCKER
    ARG VERSION="dev-$EARTHLY_TARGET_TAG_DOCKER"
    ARG EARTHLY_GIT_HASH
    ARG DEFAULT_BUILDKITD_IMAGE=docker.io/earthly/buildkitd:$VERSION # The image needs to be fully qualified for alternative frontend support.
    ARG BUILD_TAGS=dfrunmount dfrunsecurity dfsecrets dfssh dfrunnetwork dfheredoc forceposix
    ARG GOCACHE=/go-cache
    RUN mkdir -p build
    RUN printf "$BUILD_TAGS" > ./build/tags && echo "$(cat ./build/tags)"
    RUN printf '-X main.DefaultBuildkitdImage='"$DEFAULT_BUILDKITD_IMAGE" > ./build/ldflags && \
        printf ' -X main.Version='"$VERSION" >> ./build/ldflags && \
        printf ' -X main.GitSha='"$EARTHLY_GIT_HASH" >> ./build/ldflags && \
        printf ' '"$GO_EXTRA_LDFLAGS" >> ./build/ldflags && \
        echo "$(cat ./build/ldflags)"
    # Important! If you change the go build options, you may need to also change them
    # in https://github.com/earthly/homebrew-earthly/blob/main/Formula/earthly.rb
    RUN --mount=type=cache,target=$GOCACHE \
        GOARM=${VARIANT#v} go build \
            -tags "$(cat ./build/tags)" \
            -ldflags "$(cat ./build/ldflags)" \
            -o build/$EXECUTABLE_NAME \
            cmd/earthly/*.go
    SAVE ARTIFACT ./build/tags
    SAVE ARTIFACT ./build/ldflags
    SAVE ARTIFACT build/$EXECUTABLE_NAME AS LOCAL "build/$GOOS/$GOARCH$VARIANT/$EXECUTABLE_NAME"
    SAVE IMAGE --cache-from=earthly/earthly:main

earthly-linux-amd64:
    COPY (+earthly/* \
        --GOARCH=amd64 \
        --VARIANT= \
        ) ./
    SAVE ARTIFACT ./*

earthly-linux-arm64:
    COPY (+earthly/* \
        --GOARCH=arm64 \
        --VARIANT= \
        --GO_EXTRA_LDFLAGS= \
        ) ./
    SAVE ARTIFACT ./*

earthly-darwin-amd64:
    COPY (+earthly/* \
        --GOOS=darwin \
        --GOARCH=amd64 \
        --VARIANT= \
        --GO_EXTRA_LDFLAGS= \
        ) ./
    SAVE ARTIFACT ./*

earthly-darwin-arm64:
    COPY (+earthly/* \
        --GOOS=darwin \
        --GOARCH=arm64 \
        --VARIANT= \
        --GO_EXTRA_LDFLAGS= \
        ) ./
    SAVE ARTIFACT ./*

earthly-windows-amd64:
    COPY (+earthly/* \
        --GOOS=windows \
        --GOARCH=amd64 \
        --VARIANT= \
        --GO_EXTRA_LDFLAGS= \
        --EXECUTABLE_NAME=earthly.exe \
        ) ./
    SAVE ARTIFACT ./*

earthly-all:
    COPY +earthly-linux-amd64/earthly ./earthly-linux-amd64
    COPY +earthly-linux-arm64/earthly ./earthly-linux-arm64
    COPY +earthly-darwin-amd64/earthly ./earthly-darwin-amd64
    COPY +earthly-darwin-arm64/earthly ./earthly-darwin-arm64
    COPY +earthly-windows-amd64/earthly.exe ./earthly-windows-amd64.exe
    SAVE ARTIFACT ./*

earthly-docker:
    ARG BUILDKIT_PROJECT
    FROM ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    RUN apk add --update --no-cache docker-cli libcap-ng-utils
    ENV EARTHLY_IMAGE=true
    COPY earthly-entrypoint.sh /usr/bin/earthly-entrypoint.sh
    ENTRYPOINT ["/usr/bin/earthly-entrypoint.sh"]
    WORKDIR /workspace
    ARG EARTHLY_TARGET_TAG_DOCKER
    ARG TAG="dev-$EARTHLY_TARGET_TAG_DOCKER"
    COPY (+earthly/earthly --VERSION=$TAG) /usr/bin/earthly
    SAVE IMAGE --push --cache-from=earthly/earthly:main earthly/earthly:$TAG

earthly-integration-test-base:
    FROM +earthly-docker
    ENV NO_DOCKER=1
    ENV NETWORK_MODE=host
    WORKDIR /test

    # The inner buildkit requires Docker hub creds to prevent rate-limiting issues.
    ARG DOCKERHUB_MIRROR
    ARG DOCKERHUB_MIRROR_INSECURE
    ARG DOCKERHUB_AUTH=true
    ARG DOCKERHUB_USER_SECRET=+secrets/DOCKERHUB_USER
    ARG DOCKERHUB_TOKEN_SECRET=+secrets/DOCKERHUB_TOKEN

    IF [ -z $DOCKERHUB_MIRROR ]
        # No mirror, easy CI and local use by all
        ENV GLOBAL_CONFIG="{disable_analytics: true, local_registry_host: 'tcp://127.0.0.1:8371', conversion_parallelism: 5}"
        IF [ "$DOCKERHUB_AUTH" = "true" ]
            RUN --secret USERNAME=$DOCKERHUB_USER_SECRET \
                --secret TOKEN=$DOCKERHUB_TOKEN_SECRET \
                docker login --username="$USERNAME" --password="$TOKEN"
        END
    ELSE
        # Use a mirror, supports mirroring Docker Hub only.
        ENV EARTHLY_ADDITIONAL_BUILDKIT_CONFIG="[registry.\"docker.io\"]
  mirrors = [\"$DOCKERHUB_MIRROR\"]"

        IF [ "$DOCKERHUB_MIRROR_INSECURE" = "true" ]
            ENV EARTHLY_ADDITIONAL_BUILDKIT_CONFIG="$EARTHLY_ADDITIONAL_BUILDKIT_CONFIG
[registry.\"$DOCKERHUB_MIRROR\"]
  http = true
  insecure = true"
        END

        # NOTE: newlines+indentation is important here, see https://github.com/earthly/earthly/issues/1764 for potential pitfalls
        # yaml will convert newlines to spaces when using regular quoted-strings, therefore we will use the literal-style (denoted by `|`)
        ENV GLOBAL_CONFIG="disable_analytics: true
local_registry_host: 'tcp://127.0.0.1:8371'
conversion_parallelism: 5
buildkit_additional_config: |
$(echo "$EARTHLY_ADDITIONAL_BUILDKIT_CONFIG" | sed "s/^/  /g")
"
        IF [ "$DOCKERHUB_AUTH" = "true" ]
            RUN --secret USERNAME=$DOCKERHUB_USER_SECRET \
                --secret TOKEN=$DOCKERHUB_TOKEN_SECRET \
                docker login $DOCKERHUB_MIRROR --username="$USERNAME" --password="$TOKEN"
        END
    END

prerelease:
    FROM alpine:3.15
    ARG BUILDKIT_PROJECT
    BUILD \
        --platform=linux/amd64 \
        --platform=linux/arm64 \
        ./buildkitd+buildkitd --TAG=prerelease  --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    COPY (+earthly-all/* --VERSION=prerelease) ./
    SAVE IMAGE --push earthly/earthlybinaries:prerelease

prerelease-script:
    FROM alpine:3.15
    COPY ./earthly ./
    # This script is useful in other repos too.
    SAVE ARTIFACT ./earthly

dind:
    BUILD +dind-alpine
    BUILD +dind-ubuntu

dind-alpine:
    FROM docker:20.10.14-dind
    COPY ./buildkitd/docker-auto-install.sh /usr/local/bin/docker-auto-install.sh
    RUN docker-auto-install.sh
    ARG EARTHLY_TARGET_TAG_DOCKER
    ARG DIND_ALPINE_TAG=alpine-$EARTHLY_TARGET_TAG_DOCKER
    ARG DOCKERHUB_USER=earthly
    SAVE IMAGE --push --cache-from=earthly/dind:main $DOCKERHUB_USER/dind:$DIND_ALPINE_TAG

dind-ubuntu:
    FROM ubuntu:20.04
    COPY ./buildkitd/docker-auto-install.sh /usr/local/bin/docker-auto-install.sh
    RUN docker-auto-install.sh
    ARG EARTHLY_TARGET_TAG_DOCKER
    ARG DIND_UBUNTU_TAG=ubuntu-$EARTHLY_TARGET_TAG_DOCKER
    ARG DOCKERHUB_USER=earthly
    SAVE IMAGE --push --cache-from=earthly/dind:ubuntu-main $DOCKERHUB_USER/dind:$DIND_UBUNTU_TAG

for-own:
    ARG BUILDKIT_PROJECT
    BUILD ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    COPY +earthly/earthly ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/own/earthly

for-linux:
    ARG BUILDKIT_PROJECT
    BUILD --platform=linux/amd64 ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    BUILD ./ast/parser+parser
    COPY +earthly-linux-amd64/earthly ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/linux/amd64/earthly

for-darwin:
    ARG BUILDKIT_PROJECT
    BUILD --platform=linux/amd64 ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    BUILD ./ast/parser+parser
    COPY +earthly-darwin-amd64/earthly ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/darwin/amd64/earthly

for-darwin-m1:
    ARG BUILDKIT_PROJECT
    BUILD --platform=linux/arm64 ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    BUILD ./ast/parser+parser
    COPY +earthly-darwin-arm64/earthly ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/darwin/arm64/earthly

for-windows:
    # BUILD --platform=linux/amd64 ./buildkitd+buildkitd
    BUILD ./ast/parser+parser
    COPY +earthly-windows-amd64/earthly.exe ./
    SAVE ARTIFACT ./earthly.exe AS LOCAL ./build/windows/amd64/earthly.exe

all-buildkitd:
    ARG BUILDKIT_PROJECT
    BUILD \
        --platform=linux/amd64 \
        --platform=linux/arm64 \
        ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"

all-dind:
    BUILD \
        --platform=linux/amd64 \
        --platform=linux/arm64 \
        +dind

all:
    BUILD +all-buildkitd
    BUILD +earthly-all
    BUILD +earthly-docker
    BUILD +prerelease
    BUILD +all-dind

test:
    BUILD +lint
    BUILD +lint-scripts
    BUILD +lint-newline-ending
    BUILD +lint-changelog
    BUILD +unit-test
    BUILD +earthly-script-no-stdout
    ARG DOCKERHUB_MIRROR
    ARG DOCKERHUB_MIRROR_INSECURE=false
    ARG DOCKERHUB_AUTH=true
    ARG DOCKERHUB_USER_SECRET=+secrets/DOCKERHUB_USER
    ARG DOCKERHUB_TOKEN_SECRET=+secrets/DOCKERHUB_TOKEN
    BUILD ./ast/tests+all \
        --DOCKERHUB_AUTH=$DOCKERHUB_AUTH \
        --DOCKERHUB_USER_SECRET=$DOCKERHUB_USER_SECRET \
        --DOCKERHUB_TOKEN_SECRET=$DOCKERHUB_TOKEN_SECRET \
        --DOCKERHUB_MIRROR=$DOCKERHUB_MIRROR \
        --DOCKERHUB_MIRROR_INSECURE=$DOCKERHUB_MIRROR_INSECURE
    BUILD ./tests+ga \
        --DOCKERHUB_AUTH=$DOCKERHUB_AUTH \
        --DOCKERHUB_USER_SECRET=$DOCKERHUB_USER_SECRET \
        --DOCKERHUB_TOKEN_SECRET=$DOCKERHUB_TOKEN_SECRET \
        --DOCKERHUB_MIRROR=$DOCKERHUB_MIRROR \
        --DOCKERHUB_MIRROR_INSECURE=$DOCKERHUB_MIRROR_INSECURE

test-all:
    BUILD +examples
    ARG DOCKERHUB_MIRROR
    ARG DOCKERHUB_MIRROR_INSECURE=false
    ARG DOCKERHUB_AUTH=true
    ARG DOCKERHUB_USER_SECRET=+secrets/DOCKERHUB_USER
    ARG DOCKERHUB_TOKEN_SECRET=+secrets/DOCKERHUB_TOKEN
    BUILD +test \
        --DOCKERHUB_AUTH=$DOCKERHUB_AUTH \
        --DOCKERHUB_USER_SECRET=$DOCKERHUB_USER_SECRET \
        --DOCKERHUB_TOKEN_SECRET=$DOCKERHUB_TOKEN_SECRET \
        --DOCKERHUB_MIRROR=$DOCKERHUB_MIRROR \
        --DOCKERHUB_MIRROR_INSECURE=$DOCKERHUB_MIRROR_INSECURE
    BUILD ./tests+experimental  \
        --DOCKERHUB_AUTH=$DOCKERHUB_AUTH \
        --DOCKERHUB_USER_SECRET=$DOCKERHUB_USER_SECRET \
        --DOCKERHUB_TOKEN_SECRET=$DOCKERHUB_TOKEN_SECRET \
        --DOCKERHUB_MIRROR=$DOCKERHUB_MIRROR \
        --DOCKERHUB_MIRROR_INSECURE=$DOCKERHUB_MIRROR_INSECURE

examples:
    BUILD +examples1
    BUILD +examples2

examples1:
    ARG TARGETARCH
    BUILD ./examples/c+docker
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
    BUILD ./examples/react+docker
    BUILD ./examples/cutoff-optimization+run
    BUILD ./examples/import+build
    BUILD ./examples/secrets+base

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
    BUILD ./examples/cache-command/npm+docker
    BUILD ./examples/cache-command/mvn+docker

license:
    COPY LICENSE ./
    SAVE ARTIFACT LICENSE
