# TODO: we must change the DOCKERHUB_USER_SECRET args to be project-based before we can change to 0.7
VERSION --shell-out-anywhere --use-copy-link --no-network 0.6

FROM golang:1.20-alpine3.17

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

# deps downloads and caches all dependencies for earthly. When called directly,
# go.mod and go.sum will be updated locally.
deps:
    FROM +base
    RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.52.2
    COPY go.mod go.sum ./
    COPY ./ast/go.mod ./ast/go.sum ./ast
    COPY ./util/deltautil/go.mod ./util/deltautil/go.sum ./util/deltautil
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

code:
    FROM +deps
    # Use BUILDKIT_PROJECT to point go.mod to a buildkit dir being actively developed. Examples:
    #   --BUILDKIT_PROJECT=../buildkit
    #   --BUILDKIT_PROJECT=github.com/earthly/buildkit:<git-ref>
    ARG BUILDKIT_PROJECT
    IF [ "$BUILDKIT_PROJECT" != "" ]
        COPY --dir "$BUILDKIT_PROJECT"+code/buildkit /buildkit
        RUN go mod edit -replace github.com/moby/buildkit=/buildkit
        RUN go mod download
    END
    # Use CLOUD_API to point go.mod to a cloud API dir being actively developed. Examples:
    #   --CLOUD_API=../cloud/api+proto/api/public/'*'
    #   --CLOUD_API=github.com/earthly/cloud/api:<git-ref>+proto/api/public/'*'
    #   --CLOUD_API=github.com/earthly/cloud-api:<git-ref>+code/'*'
    ARG CLOUD_API
    IF [ "$CLOUD_API" != "" ]
        COPY --dir "$CLOUD_API" /cloud-api/
        RUN go mod edit -replace github.com/earthly/cloud-api=/cloud-api
        RUN go mod download
    END
    COPY ./ast/parser+parser/*.go ./ast/parser/
    COPY --dir analytics autocomplete buildcontext builder logbus cleanup cmd config conslogging debugger \
        dockertar docker2earthly domain features outmon slog cloud states util variables ./
    COPY --dir buildkitd/buildkitd.go buildkitd/settings.go buildkitd/certificates.go buildkitd/
    COPY --dir earthfile2llb/*.go earthfile2llb/
    COPY --dir ast/antlrhandler ast/spec ast/hint ast/command ast/commandflag ast/*.go ast/
    COPY --dir inputgraph/*.go inputgraph/

# update-buildkit updates earthly's buildkit dependency.
update-buildkit:
    FROM +code # if we use deps, go mod tidy will remove a bunch of requirements since it won't have access to our codebase.
    ARG BUILDKIT_GIT_SHA
    ARG BUILDKIT_GIT_BRANCH=earthly-main
    ARG BUILDKIT_GIT_ORG=earthly
    COPY (./buildkitd+buildkit-sha/buildkit_sha --BUILDKIT_GIT_ORG="$BUILDKIT_GIT_ORG" --BUILDKIT_GIT_SHA="$BUILDKIT_GIT_SHA" --BUILDKIT_GIT_BRANCH="$BUILDKIT_GIT_BRANCH") buildkit_sha
    BUILD  ./buildkitd+update-buildkit-earthfile --BUILDKIT_GIT_ORG="$BUILDKIT_GIT_ORG" --BUILDKIT_GIT_SHA="$(cat buildkit_sha)"
    RUN --no-cache go mod edit -replace "github.com/moby/buildkit=github.com/$BUILDKIT_GIT_ORG/buildkit@$(cat buildkit_sha)"
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
        ./buildkitd/dockerd-wrapper.sh ./buildkitd/docker-auto-install.sh ./buildkitd/oom-adjust.sh.template \
        ./release/envcredhelper.sh ./release/ami/cleanup.sh ./release/ami/configure.sh ./release/ami/install.sh \
        ./.buildkite/*.sh \
        ./scripts/tests/*.sh \
        ./scripts/tests/docker-build/*.sh \
        ./scripts/*.sh \
        ./shell_scripts/
    # some scripts need to source /etc/os-release for operating system release information,
    # so -x is needed to let shellcheck read that file.
    RUN shellcheck -x shell_scripts/*

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
    RUN EARTHLY_DISABLE_FRONTEND_DETECTION=true EARTHLY_DISABLE_AUTO_UPDATE=true ./earthly --version > earthly-version-output

    RUN test "$(cat earthly-version-output | wc -l)" = "1"
    RUN grep '^earthly version.*$' earthly-version-output # only --version info should go to stdout

# lint runs basic go linters against the earthly project.
lint:
    FROM +code
    COPY ./.golangci.yaml ./
    RUN golangci-lint run

# lint-newline-ending checks that line endings are unix style and that files end
# with a single newline.
lint-newline-ending:
    FROM alpine:3.15
    WORKDIR /everything
    COPY . .
    # test that line endings are unix-style
    RUN set -e; \
        code=0; \
        for f in $(find . -not -path "./.git/*" -type f \( -iname '*.go' -o -iname 'Earthfile' -o -iname '*.earth' -o -iname '*.md' -o -iname '*.json'\) | grep -v "ast/tests/empty-targets.earth" ); do \
            if ! dos2unix < "$f" | cmp - "$f"; then \
                echo "$f contains windows-style newlines and must be converted to unix-style (use dos2unix to fix)"; \
                code=1; \
            fi; \
        done; \
        exit $code
    # test file ends with a single newline
    RUN set -e; \
        code=0; \
        for f in $(find . -not -path "./.git/*" -type f \( -iname '*.yml' -o -iname '*.go' -o -iname '*.sh' -o -iname '*.template' -o -iname 'Earthfile' -o -iname '*.earth' -o -iname '*.md' -o -iname '*.json' \) | grep -v "ast/tests/empty-targets.earth" | grep -v "tests/version/version-only.earth" ); do \
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
        for f in $(find . -not -path "./.git/*" -type f \( -iname '*.go' -o -iname 'Earthfile' -o -iname '*.earth' -o -iname '*.md' -o -iname '*.json'\) | grep -v "ast/tests/empty-targets.earth" | grep -v "ast/parser/earth_parser.go" | grep -v "ast/parser/earth_lexer.go" ); do \
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

# mocks runs 'go generate' against this module and saves generated mock files
# locally.
mocks:
    FROM +code
    RUN go install git.sr.ht/~nelsam/hel@latest && go install golang.org/x/tools/cmd/goimports@latest
    RUN go generate ./...
    FOR mockfile IN $(find . -name 'helheim*_test.go')
        SAVE ARTIFACT $mockfile AS LOCAL $mockfile
    END

# unit-test runs unit tests (and some integration tests).
unit-test:
    FROM +code
    RUN apk add --no-cache --update podman fuse-overlayfs
    COPY not-a-unit-test.sh .

    ARG testname # when specified, only run specific unit-test, otherwise run all.

    # pkgname determines the package name (or names) that will be tested. The go
    # submodules must be specified explicitly or they will not be run, as
    # "./..." does not match submodules.
    ARG pkgname = ./...

    ARG DOCKERHUB_MIRROR
    ARG DOCKERHUB_MIRROR_INSECURE=false
    ARG DOCKERHUB_MIRROR_HTTP=false
    ARG DOCKERHUB_AUTH=true
    ARG DOCKERHUB_USER_SECRET=+secrets/DOCKERHUB_USER
    ARG DOCKERHUB_TOKEN_SECRET=+secrets/DOCKERHUB_TOKEN
    IF [ -n "$DOCKERHUB_MIRROR" ]
        RUN mkdir -p /etc/docker
        RUN echo "{\"registry-mirrors\": [\"http://$DOCKERHUB_MIRROR\"]" > /etc/docker/daemon.json
        IF [ "$DOCKERHUB_MIRROR_INSECURE" = "true" ] || [ "$DOCKERHUB_MIRROR_HTTP" = "true" ]
          RUN echo ", \"insecure-registries\": [\"$DOCKERHUB_MIRROR\"]" >> /etc/docker/daemon.json
        END
        RUN echo "}" >> /etc/docker/daemon.json
    END
    IF [ "$DOCKERHUB_AUTH" = "true" ]
        WITH DOCKER
            RUN --secret USERNAME=$DOCKERHUB_USER_SECRET \
                --secret TOKEN=$DOCKERHUB_TOKEN_SECRET \
                ./not-a-unit-test.sh
        END
    ELSE
        WITH DOCKER
            RUN testname=$testname pkgname=$pkgname ./not-a-unit-test.sh
        END
    END

    # The following are separate go modules and need to be tested separately.
    # The not-a-unit-test.sh script above actually DOES run unit-tests as well
    BUILD ./ast+unit-test
    BUILD ./util/deltautil+unit-test

# chaos-test runs tests that use chaos and load in order to exercise components
# of earthly. These tests may be more resource-intensive or flaky than typical
# unit or integration tests.
#
# Since the race detector (-race) sets a goroutine limit, these tests are run
# without -race.
chaos-test:
    FROM +code
    RUN go test -tags chaos ./...

offline-test:
    FROM +code
    RUN --network=none go test -run TestOffline ./...

# submodule-decouple-check checks that go submodules within earthly do not
# depend on the core earthly project.
submodule-decouple-check:
    FROM +code
    RUN for submodule in github.com/earthly/earthly/ast github.com/earthly/earthly/util/deltautil; \
    do \
        for dep in $(go list -f '{{range .Deps}}{{.}} {{end}}' $submodule/...); \
        do \
            if [ "$(go list -f '{{if .Module}}{{.Module}}{{end}}' $dep)" == "github.com/earthly/earthly" ]; \
            then \
               echo "FAIL: submodule $submodule imports $dep, which is in the core 'github.com/earthly/earthly' module"; \
               exit 1; \
            fi; \
        done; \
    done

changelog:
    FROM scratch
    COPY CHANGELOG.md .
    SAVE ARTIFACT CHANGELOG.md

lint-changelog:
    FROM python:3
    COPY release/changelogparser.py /usr/bin/changelogparser
    COPY CHANGELOG.md .
    RUN changelogparser --changelog CHANGELOG.md

debugger:
    FROM +code
    ENV CGO_ENABLED=0
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

# earthly builds the earthly CLI and docker image.
earthly:
    FROM +code
    ENV CGO_ENABLED=0
    ARG GOOS=linux
    ARG TARGETARCH
    ARG TARGETVARIANT
    ARG GOARCH=$TARGETARCH
    ARG VARIANT=$TARGETVARIANT
    ARG GO_EXTRA_LDFLAGS="-linkmode external -extldflags -static"
    # GO_GCFLAGS may be used to set the -gcflags parameter to 'go build'. This
    # is particularly useful for disabling optimizations to make the binary work
    # with delve. To disable optimizations:
    #
    #     -GO_GCFLAGS='all=-N -l'
    ARG GO_GCFLAGS
    ARG EXECUTABLE_NAME="earthly"
    ARG DEFAULT_INSTALLATION_NAME="earthly-dev"
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
        printf ' -X main.DefaultInstallationName='"$DEFAULT_INSTALLATION_NAME" >> ./build/ldflags && \
        printf ' '"$GO_EXTRA_LDFLAGS" >> ./build/ldflags && \
        echo "$(cat ./build/ldflags)"
    # Important! If you change the go build options, you may need to also change them
    # in https://github.com/earthly/homebrew-earthly/blob/main/Formula/earthly.rb
    RUN --mount=type=cache,target=$GOCACHE \
        GOARM=${VARIANT#v} go build \
            -tags "$(cat ./build/tags)" \
            -ldflags "$(cat ./build/ldflags)" \
            -gcflags="${GO_GCFLAGS}" \
            -o build/$EXECUTABLE_NAME \
            cmd/earthly/*.go
    SAVE ARTIFACT ./build/tags
    SAVE ARTIFACT ./build/ldflags
    SAVE ARTIFACT build/$EXECUTABLE_NAME AS LOCAL "build/$GOOS/$GOARCH$VARIANT/$EXECUTABLE_NAME"
    SAVE IMAGE --cache-from=earthly/earthly:main

earthly-linux-amd64:
    ARG GO_GCFLAGS
    COPY (+earthly/* \
        --GOARCH=amd64 \
        --VARIANT= \
        --GO_GCFLAGS="${GO_GCFLAGS}" \
        ) ./
    SAVE ARTIFACT ./*

earthly-linux-arm64:
    ARG GO_GCFLAGS
    COPY (+earthly/* \
        --GOARCH=arm64 \
        --VARIANT= \
        --GO_EXTRA_LDFLAGS= \
        --GO_GCFLAGS="${GO_GCFLAGS}" \
        ) ./
    SAVE ARTIFACT ./*

earthly-darwin-amd64:
    ARG GO_GCFLAGS=""
    COPY (+earthly/* \
        --GOOS=darwin \
        --GOARCH=amd64 \
        --VARIANT= \
        --GO_EXTRA_LDFLAGS= \
        --GO_GCFLAGS="${GO_GCFLAGS}" \
        ) ./
    SAVE ARTIFACT ./*

earthly-darwin-arm64:
    ARG GO_GCFLAGS
    COPY (+earthly/* \
        --GOOS=darwin \
        --GOARCH=arm64 \
        --VARIANT= \
        --GO_EXTRA_LDFLAGS= \
        --GO_GCFLAGS="${GO_GCFLAGS}" \
        ) ./
    SAVE ARTIFACT ./*

earthly-windows-amd64:
    ARG GO_GCFLAGS
    COPY (+earthly/* \
        --GOOS=windows \
        --GOARCH=amd64 \
        --VARIANT= \
        --GO_EXTRA_LDFLAGS= \
        --GO_GCFLAGS="${GO_GCFLAGS}" \
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
    ARG EARTHLY_TARGET_TAG_DOCKER
    ARG TAG="dev-$EARTHLY_TARGET_TAG_DOCKER"
    ARG BUILDKIT_PROJECT
    FROM ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT" --TAG="$TAG"
    RUN apk add --update --no-cache docker-cli libcap-ng-utils git
    ENV EARTHLY_IMAGE=true
    COPY earthly-entrypoint.sh /usr/bin/earthly-entrypoint.sh
    ENTRYPOINT ["/usr/bin/earthly-entrypoint.sh"]
    WORKDIR /workspace
    COPY (+earthly/earthly --VERSION=$TAG --DEFAULT_INSTALLATION_NAME="earthly") /usr/bin/earthly
    ARG DOCKERHUB_USER="earthly"
    ARG DOCKERHUB_IMG="earthly"
    SAVE IMAGE --push --cache-from=earthly/earthly:main $DOCKERHUB_USER/$DOCKERHUB_IMG:$TAG

earthly-integration-test-base:
    FROM +earthly-docker
    RUN apk update && apk add pcre-tools curl python3 bash perl findutils
    ENV NO_DOCKER=1
    ENV NETWORK_MODE=host # Note that this breaks access to embedded registry in WITH DOCKER.
    ENV EARTHLY_VERSION_FLAG_OVERRIDES=no-use-registry-for-with-docker # Use tar-based due to above.
    WORKDIR /test

    # The inner buildkit requires Docker hub creds to prevent rate-limiting issues.
    ARG DOCKERHUB_MIRROR
    ARG DOCKERHUB_MIRROR_INSECURE
    ARG DOCKERHUB_MIRROR_HTTP
    ARG DOCKERHUB_AUTH=true
    ARG DOCKERHUB_USER_SECRET=+secrets/DOCKERHUB_USER
    ARG DOCKERHUB_TOKEN_SECRET=+secrets/DOCKERHUB_TOKEN

    IF [ -z "$DOCKERHUB_MIRROR" ]
        # No mirror, easy CI and local use by all
        ENV GLOBAL_CONFIG="{disable_analytics: true}"
        IF [ "$DOCKERHUB_AUTH" = "true" ]
            RUN --secret USERNAME=$DOCKERHUB_USER_SECRET \
                --secret TOKEN=$DOCKERHUB_TOKEN_SECRET \
                (test -n "$USERNAME" || (echo "ERROR: USERNAME not set"; exit 1)) && \
                (test -n "$TOKEN" || (echo "ERROR: TOKEN not set"; exit 1)) && \
                docker login --username="$USERNAME" --password="$TOKEN"
        END
    ELSE
        # Use a mirror, supports mirroring Docker Hub only.
        ENV EARTHLY_ADDITIONAL_BUILDKIT_CONFIG="[registry.\"docker.io\"]
  mirrors = [\"$DOCKERHUB_MIRROR\"]"
        ENV MIRROR_CONFIG="[registry.\"$DOCKERHUB_MIRROR\"]"
        IF [ "$DOCKERHUB_MIRROR_INSECURE" = "true" ]
            ENV EARTHLY_ADDITIONAL_BUILDKIT_CONFIG="$EARTHLY_ADDITIONAL_BUILDKIT_CONFIG
  insecure = true"
            ENV MIRROR_CONFIG="$MIRROR_CONFIG
  insecure = true"
        END
        IF [ "$DOCKERHUB_MIRROR_HTTP" = "true" ]
            ENV EARTHLY_ADDITIONAL_BUILDKIT_CONFIG="$EARTHLY_ADDITIONAL_BUILDKIT_CONFIG
  http = true"
            ENV MIRROR_CONFIG="$MIRROR_CONFIG
  http = true"
        END
        ENV EARTHLY_ADDITIONAL_BUILDKIT_CONFIG="$EARTHLY_ADDITIONAL_BUILDKIT_CONFIG
$MIRROR_CONFIG"

        # NOTE: newlines+indentation is important here, see https://github.com/earthly/earthly/issues/1764 for potential pitfalls
        # yaml will convert newlines to spaces when using regular quoted-strings, therefore we will use the literal-style (denoted by `|`)
        ENV GLOBAL_CONFIG="disable_analytics: true
buildkit_additional_config: |
$(echo "$EARTHLY_ADDITIONAL_BUILDKIT_CONFIG" | sed "s/^/  /g")
"
        IF [ "$DOCKERHUB_AUTH" = "true" ]
            RUN --secret USERNAME=$DOCKERHUB_USER_SECRET \
                --secret TOKEN=$DOCKERHUB_TOKEN_SECRET \
                (test -n "$DOCKERHUB_MIRROR" || (echo "ERROR: DOCKERHUB_MIRROR not set"; exit 1)) && \
                (test -n "$USERNAME" || (echo "ERROR: USERNAME not set"; exit 1)) && \
                (test -n "$TOKEN" || (echo "ERROR: TOKEN not set"; exit 1)) && \
                docker login "$DOCKERHUB_MIRROR" --username="$USERNAME" --password="$TOKEN"
        END
    END

prerelease:
    FROM alpine:3.15
    ARG BUILDKIT_PROJECT
    BUILD \
        --platform=linux/amd64 \
        --platform=linux/arm64 \
        ./buildkitd+buildkitd --TAG=prerelease  --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    COPY (+earthly-all/* --VERSION=prerelease --DEFAULT_INSTALLATION_NAME=earthly) ./
    SAVE IMAGE --push earthly/earthlybinaries:prerelease

prerelease-script:
    FROM alpine:3.15
    COPY ./earthly ./
    # This script is useful in other repos too.
    SAVE ARTIFACT ./earthly

ci-release:
    # TODO: this was multiplatform, but that skyrocketed our build times. #2979
    # may help.
    FROM alpine:3.15
    ARG BUILDKIT_PROJECT
    ARG EARTHLY_GIT_HASH
    ARG --required TAG_SUFFIX
    BUILD \
        --platform=linux/amd64 \
        ./buildkitd+buildkitd --TAG=${EARTHLY_GIT_HASH}-${TAG_SUFFIX} --BUILDKIT_PROJECT="$BUILDKIT_PROJECT" --DOCKERHUB_BUILDKIT_IMG="buildkitd-staging"
    COPY (+earthly/earthly --DEFAULT_BUILDKITD_IMAGE="docker.io/earthly/buildkitd-staging:${EARTHLY_GIT_HASH}-${TAG_SUFFIX}" --VERSION=${EARTHLY_GIT_HASH}-${TAG_SUFFIX} --DEFAULT_INSTALLATION_NAME=earthly) ./earthly-linux-amd64
    SAVE IMAGE --push earthly/earthlybinaries:${EARTHLY_GIT_HASH}-${TAG_SUFFIX}

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

# for-own builds earthly-buildkitd and the earthly CLI for the current system
# and saves the final CLI binary locally.
for-own:
    ARG BUILDKIT_PROJECT
    # GO_GCFLAGS may be used to set the -gcflags parameter to 'go build'. See
    # the documentaation on +earthly for extra detail about this option.
    ARG GO_GCFLAGS
    BUILD ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    COPY (+earthly/earthly --GO_GCFLAGS="${GO_GCFLAGS}") ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/own/earthly

for-linux:
    ARG BUILDKIT_PROJECT
    ARG GO_GCFLAGS
    BUILD --platform=linux/amd64 ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    BUILD ./ast/parser+parser
    COPY (+earthly-linux-amd64/earthly -GO_GCFLAGS="${GO_GCFLAGS}") ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/linux/amd64/earthly

for-darwin:
    ARG BUILDKIT_PROJECT
    ARG GO_GCFLAGS
    BUILD --platform=linux/amd64 ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    BUILD ./ast/parser+parser
    COPY (+earthly-darwin-amd64/earthly -GO_GCFLAGS="${GO_GCFLAGS}") ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/darwin/amd64/earthly

for-darwin-m1:
    ARG BUILDKIT_PROJECT
    ARG GO_GCFLAGS
    BUILD --platform=linux/arm64 ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    BUILD ./ast/parser+parser
    COPY (+earthly-darwin-arm64/earthly -GO_GCFLAGS="${GO_GCFLAGS}") ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/darwin/arm64/earthly

for-windows:
    ARG GO_GCFLAGS
    # BUILD --platform=linux/amd64 ./buildkitd+buildkitd
    BUILD ./ast/parser+parser
    COPY (+earthly-windows-amd64/earthly.exe -GO_GCFLAGS="${GO_GCFLAGS}") ./
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

# lint-all runs all linting checks against the earthly project.
lint-all:
    BUILD +lint
    BUILD +lint-scripts
    BUILD +lint-docs
    BUILD +submodule-decouple-check

lint-docs:
    BUILD +lint-newline-ending
    BUILD +lint-changelog

# TODO: Document qemu vs non-qemu
test-no-qemu:
    BUILD +unit-test
    BUILD +chaos-test
    BUILD +offline-test
    BUILD +earthly-script-no-stdout
    ARG DOCKERHUB_MIRROR
    ARG DOCKERHUB_MIRROR_INSECURE=false
    ARG DOCKERHUB_MIRROR_HTTP=false
    ARG DOCKERHUB_AUTH=true
    ARG DOCKERHUB_USER_SECRET=+secrets/DOCKERHUB_USER
    ARG DOCKERHUB_TOKEN_SECRET=+secrets/DOCKERHUB_TOKEN
    BUILD ./ast/tests+all \
        --DOCKERHUB_AUTH=$DOCKERHUB_AUTH \
        --DOCKERHUB_USER_SECRET=$DOCKERHUB_USER_SECRET \
        --DOCKERHUB_TOKEN_SECRET=$DOCKERHUB_TOKEN_SECRET \
        --DOCKERHUB_MIRROR=$DOCKERHUB_MIRROR \
        --DOCKERHUB_MIRROR_INSECURE=$DOCKERHUB_MIRROR_INSECURE \
        --DOCKERHUB_MIRROR_HTTP=$DOCKERHUB_MIRROR_HTTP
    ARG GLOBAL_WAIT_END="false"
    BUILD ./tests+ga-no-qemu \
        --DOCKERHUB_AUTH=$DOCKERHUB_AUTH \
        --DOCKERHUB_USER_SECRET=$DOCKERHUB_USER_SECRET \
        --DOCKERHUB_TOKEN_SECRET=$DOCKERHUB_TOKEN_SECRET \
        --DOCKERHUB_MIRROR=$DOCKERHUB_MIRROR \
        --DOCKERHUB_MIRROR_INSECURE=$DOCKERHUB_MIRROR_INSECURE \
        --DOCKERHUB_MIRROR_HTTP=$DOCKERHUB_MIRROR_HTTP \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

test-qemu:
    ARG DOCKERHUB_MIRROR
    ARG DOCKERHUB_MIRROR_INSECURE=false
    ARG DOCKERHUB_MIRROR_HTTP=false
    ARG DOCKERHUB_AUTH=true
    ARG DOCKERHUB_USER_SECRET=+secrets/DOCKERHUB_USER
    ARG DOCKERHUB_TOKEN_SECRET=+secrets/DOCKERHUB_TOKEN
    ARG GLOBAL_WAIT_END="false"
    BUILD ./tests+ga-qemu \
        --DOCKERHUB_AUTH=$DOCKERHUB_AUTH \
        --DOCKERHUB_USER_SECRET=$DOCKERHUB_USER_SECRET \
        --DOCKERHUB_TOKEN_SECRET=$DOCKERHUB_TOKEN_SECRET \
        --DOCKERHUB_MIRROR=$DOCKERHUB_MIRROR \
        --DOCKERHUB_MIRROR_INSECURE=$DOCKERHUB_MIRROR_INSECURE \
        --DOCKERHUB_MIRROR_HTTP=$DOCKERHUB_MIRROR_HTTP \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

test:
    ARG DOCKERHUB_MIRROR
    ARG DOCKERHUB_MIRROR_INSECURE=false
    ARG DOCKERHUB_MIRROR_HTTP=false
    ARG DOCKERHUB_AUTH=true
    ARG DOCKERHUB_USER_SECRET=+secrets/DOCKERHUB_USER
    ARG DOCKERHUB_TOKEN_SECRET=+secrets/DOCKERHUB_TOKEN
    BUILD +test-no-qemu \
        --DOCKERHUB_AUTH=$DOCKERHUB_AUTH \
        --DOCKERHUB_USER_SECRET=$DOCKERHUB_USER_SECRET \
        --DOCKERHUB_TOKEN_SECRET=$DOCKERHUB_TOKEN_SECRET \
        --DOCKERHUB_MIRROR=$DOCKERHUB_MIRROR \
        --DOCKERHUB_MIRROR_INSECURE=$DOCKERHUB_MIRROR_INSECURE \
        --DOCKERHUB_MIRROR_HTTP=$DOCKERHUB_MIRROR_HTTP
    BUILD +test-qemu \
        --DOCKERHUB_AUTH=$DOCKERHUB_AUTH \
        --DOCKERHUB_USER_SECRET=$DOCKERHUB_USER_SECRET \
        --DOCKERHUB_TOKEN_SECRET=$DOCKERHUB_TOKEN_SECRET \
        --DOCKERHUB_MIRROR=$DOCKERHUB_MIRROR \
        --DOCKERHUB_MIRROR_INSECURE=$DOCKERHUB_MIRROR_INSECURE \
        --DOCKERHUB_MIRROR_HTTP=$DOCKERHUB_MIRROR_HTTP

test-all:
    BUILD +examples
    ARG DOCKERHUB_MIRROR
    ARG DOCKERHUB_MIRROR_INSECURE=false
    ARG DOCKERHUB_MIRROR_HTTP=false
    ARG DOCKERHUB_AUTH=true
    ARG DOCKERHUB_USER_SECRET=+secrets/DOCKERHUB_USER
    ARG DOCKERHUB_TOKEN_SECRET=+secrets/DOCKERHUB_TOKEN
    BUILD +test-no-qemu \
        --DOCKERHUB_AUTH=$DOCKERHUB_AUTH \
        --DOCKERHUB_USER_SECRET=$DOCKERHUB_USER_SECRET \
        --DOCKERHUB_TOKEN_SECRET=$DOCKERHUB_TOKEN_SECRET \
        --DOCKERHUB_MIRROR=$DOCKERHUB_MIRROR \
        --DOCKERHUB_MIRROR_INSECURE=$DOCKERHUB_MIRROR_INSECURE \
        --DOCKERHUB_MIRROR_HTTP=$DOCKERHUB_MIRROR_HTTP
    BUILD +test-qemu \
        --DOCKERHUB_AUTH=$DOCKERHUB_AUTH \
        --DOCKERHUB_USER_SECRET=$DOCKERHUB_USER_SECRET \
        --DOCKERHUB_TOKEN_SECRET=$DOCKERHUB_TOKEN_SECRET \
        --DOCKERHUB_MIRROR=$DOCKERHUB_MIRROR \
        --DOCKERHUB_MIRROR_INSECURE=$DOCKERHUB_MIRROR_INSECURE \
        --DOCKERHUB_MIRROR_HTTP=$DOCKERHUB_MIRROR_HTTP
    BUILD ./tests+experimental  \
        --DOCKERHUB_AUTH=$DOCKERHUB_AUTH \
        --DOCKERHUB_USER_SECRET=$DOCKERHUB_USER_SECRET \
        --DOCKERHUB_TOKEN_SECRET=$DOCKERHUB_TOKEN_SECRET \
        --DOCKERHUB_MIRROR=$DOCKERHUB_MIRROR \
        --DOCKERHUB_MIRROR_INSECURE=$DOCKERHUB_MIRROR_INSECURE \
        --DOCKERHUB_MIRROR_HTTP=$DOCKERHUB_MIRROR_HTTP

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
    BUILD ./examples/cloud-secrets+base

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
    BUILD ./examples/typescript-node+docker
    BUILD ./examples/bazel+run
    BUILD ./examples/bazel+image
    BUILD ./examples/aws-sso+base
    BUILD ./examples/mkdocs+build

license:
    COPY LICENSE ./
    SAVE ARTIFACT LICENSE

# npm-update-all helps keep all node package-lock.json files up to date.
npm-update-all:
    FROM node:16.16.0-alpine3.15
    COPY . /code
    WORKDIR /code
    FOR nodepath IN \
            contrib/earthfile-syntax-highlighting \
            examples/cache-command/npm \
            examples/js \
            examples/react \
            examples/ruby-on-rails \
            examples/tutorial/js/part3 \
            examples/tutorial/js/part4 \
            examples/tutorial/js/part5/services/service-one \
            examples/tutorial/js/part6/api \
            examples/tutorial/js/part6/app \
            tests/remote-cache/test2
        RUN cd $nodepath && npm update
        SAVE ARTIFACT --if-exists $nodepath/package-lock.json AS LOCAL $nodepath/package-lock.json
    END
