VERSION --pass-args --no-network --arg-scope-and-set 0.7
PROJECT earthly-technologies/core

# TODO update to 3.18; however currently "podman login" (used under not-a-unit-test.sh) will error with
# "Error: default OCI runtime "crun" not found: invalid argument".
FROM golang:1.21-alpine3.17

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
    RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.1
    COPY go.mod go.sum ./
    COPY ./ast/go.mod ./ast/go.sum ./ast
    COPY ./util/deltautil/go.mod ./util/deltautil/go.sum ./util/deltautil
    RUN go mod download
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

# code downloads and caches all dependencies for earthly and then copies the go code
# directories into the image.
# If BUILDKIT_PROJECT or CLOUD_API environment variables are set it will also update the go mods
# for the local versions
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
    COPY --dir analytics autocomplete buildcontext builder logbus cleanup cloud cmd config conslogging debugger \
        dockertar docker2earthly domain features internal outmon slog states util variables ./
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
    FROM alpine:3.18

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

# lint-scripts runs the shellcheck package to detect potential errors in shell scripts
lint-scripts:
    BUILD +lint-scripts-auth-test
    BUILD +lint-scripts-misc

# earthly-script-no-stdout validates the ./earthly script doesn't print anything to stdout (stderr only)
# This is to ensure commands such as: MYSECRET="$(./earthly secrets get -n /user/my-secret)" work
earthly-script-no-stdout:
    # This validates the ./earthly script doesn't print anything to stdout (it should print to stderr)
    # This is to ensure commands such as: MYSECRET="$(./earthly secrets get -n /user/my-secret)" work
    FROM earthly/dind:alpine-3.18
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

lint-newline-ending:
    FROM alpine:3.18
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

# markdown-spellcheck runs vale against md files
markdown-spellcheck:
    FROM --platform=linux/amd64 +vale
    WORKDIR /everything
    COPY . .
    # TODO figure out a way to ignore this pattern in vale (doesn't seem to be working under spelling's filter option)
    RUN find . -type f -iname '*.md' | xargs -n 1 sed -i 's/{[^}]*}//g'
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
    ARG DOCKERHUB_MIRROR_AUTH=false
    ARG DOCKERHUB_MIRROR_AUTH_FROM_CLOUD_SECRETS=false

    IF [ -n "$DOCKERHUB_MIRROR" ]
        RUN mkdir -p /etc/docker
        RUN echo "{\"registry-mirrors\": [\"http://$DOCKERHUB_MIRROR\"]" > /etc/docker/daemon.json
        IF [ "$DOCKERHUB_MIRROR_INSECURE" = "true" ] || [ "$DOCKERHUB_MIRROR_HTTP" = "true" ]
          RUN echo ", \"insecure-registries\": [\"$DOCKERHUB_MIRROR\"]" >> /etc/docker/daemon.json
        END
        RUN echo "}" >> /etc/docker/daemon.json
    END
    IF [ "$DOCKERHUB_MIRROR_AUTH_FROM_CLOUD_SECRETS" = "true" ]
        RUN if [ "$DOCKERHUB_MIRROR_AUTH" = "true" ]; then echo "ERROR: DOCKERHUB_MIRROR_AUTH_FROM_CLOUD_SECRETS and DOCKERHUB_MIRROR_AUTH are mutually exclusive" && exit 1; fi
        WITH DOCKER
            RUN --secret DOCKERHUB_MIRROR_USER=dockerhub-mirror/user \
                --secret DOCKERHUB_MIRROR_PASS=dockerhub-mirror/pass \
                USE_EARTHLY_MIRROR=true ./not-a-unit-test.sh
        END
    ELSE IF [ "$DOCKERHUB_MIRROR_AUTH" = "true" ]
        WITH DOCKER
            RUN --secret DOCKERHUB_MIRROR_USER \
                --secret DOCKERHUB_MIRROR_PASS \
                ./not-a-unit-test.sh
        END
    ELSE
        RUN ./not-a-unit-test.sh
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

# offline-test runs offline tests with network set to none
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

# changelog saves the CHANGELOG.md as an artifact
changelog:
    FROM scratch
    COPY CHANGELOG.md .
    SAVE ARTIFACT CHANGELOG.md

lint-changelog:
    FROM python:3
    COPY release/changelogparser.py /usr/bin/changelogparser
    COPY CHANGELOG.md .
    RUN changelogparser --changelog CHANGELOG.md

# debugger builds the earthly debugger and saves the artifact in build/earth_debugger
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
    # as well as https://github.com/Homebrew/homebrew-core/blob/master/Formula/earthly.rb 
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

# earthly-linux-amd64 builds the earthly artifact  for linux amd64
earthly-linux-amd64:
    ARG GO_GCFLAGS
    COPY --platform=linux/amd64 (+earthly/* \
        --GOARCH=amd64 \
        --VARIANT= \
        --GO_GCFLAGS="${GO_GCFLAGS}" \
        ) ./
    SAVE ARTIFACT ./*

# earthly-linux-arm64 builds the earthly artifact  for linux arm64
earthly-linux-arm64:
    ARG GO_GCFLAGS
    COPY (+earthly/* \
        --GOARCH=arm64 \
        --VARIANT= \
        --GO_EXTRA_LDFLAGS= \
        --GO_GCFLAGS="${GO_GCFLAGS}" \
        ) ./
    SAVE ARTIFACT ./*

# earthly-darwin-amd64 builds the earthly artifact  for darwin amd64
earthly-darwin-amd64:
    ARG GO_GCFLAGS=""
    COPY --platform=linux/amd64 (+earthly/* \
        --GOOS=darwin \
        --GOARCH=amd64 \
        --VARIANT= \
        --GO_EXTRA_LDFLAGS= \
        --GO_GCFLAGS="${GO_GCFLAGS}" \
        ) ./
    SAVE ARTIFACT ./*

# earthly-darwin-arm64 builds the earthly artifact for darwin arm64
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

# earthly-windows-arm64 builds the earthly artifact  for windows arm64
earthly-windows-amd64:
    ARG GO_GCFLAGS
    COPY --platform=linux/amd64 (+earthly/* \
        --GOOS=windows \
        --GOARCH=amd64 \
        --VARIANT= \
        --GO_EXTRA_LDFLAGS= \
        --GO_GCFLAGS="${GO_GCFLAGS}" \
        --EXECUTABLE_NAME=earthly.exe \
        ) ./
    SAVE ARTIFACT ./*

# earthly-all builds earthly for all supported environments
# This includes:
# linux amd64 and linux arm64
# Darwin amd64 and arm64
# Windows amd64
earthly-all:
    COPY +earthly-linux-amd64/earthly ./earthly-linux-amd64
    COPY +earthly-linux-arm64/earthly ./earthly-linux-arm64
    COPY +earthly-darwin-amd64/earthly ./earthly-darwin-amd64
    COPY +earthly-darwin-arm64/earthly ./earthly-darwin-arm64
    COPY +earthly-windows-amd64/earthly.exe ./earthly-windows-amd64.exe
    SAVE ARTIFACT ./*

# earthly-docker builds earthly as a docker image and pushes
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

# earthly-integration-test-base builds earthly docker and then
# if no dockerhub mirror is not set it will attempt to login to dockerhub using the provided docker hub username and token.
# Otherwise, it will attempt to login to the docker hub mirror using the provided username and password
earthly-integration-test-base:
    FROM +earthly-docker
    RUN apk update && apk add pcre-tools curl python3 bash perl findutils expect yq
    COPY scripts/acbtest/acbtest scripts/acbtest/acbgrep /bin/
    ENV NO_DOCKER=1
    ENV NETWORK_MODE=host # Note that this breaks access to embedded registry in WITH DOCKER.
    ENV EARTHLY_VERSION_FLAG_OVERRIDES=no-use-registry-for-with-docker # Use tar-based due to above.
    WORKDIR /test

    # The inner buildkit requires Docker hub creds to prevent rate-limiting issues.
    ARG DOCKERHUB_MIRROR
    ARG DOCKERHUB_MIRROR_INSECURE=false
    ARG DOCKERHUB_MIRROR_HTTP=false
    ARG DOCKERHUB_MIRROR_AUTH=false
    ARG DOCKERHUB_MIRROR_AUTH_FROM_CLOUD_SECRETS=false

    # DOCKERHUB_AUTH will login to docker hub (and pull from docker hub rather than a mirror)
    ARG DOCKERHUB_AUTH=false

    COPY setup-registry.sh .
    IF [ "$DOCKERHUB_MIRROR_AUTH_FROM_CLOUD_SECRETS" = "true" ]
        RUN if [ "$DOCKERHUB_MIRROR_AUTH" = "true" ]; then echo "ERROR: DOCKERHUB_MIRROR_AUTH_FROM_CLOUD_SECRETS and DOCKERHUB_MIRROR_AUTH are mutually exclusive" && exit 1; fi
        RUN --secret DOCKERHUB_MIRROR_USER=dockerhub-mirror/user --secret DOCKERHUB_MIRROR_PASS=dockerhub-mirror/pass USE_EARTHLY_MIRROR=true ./setup-registry.sh
    ELSE IF [ "$DOCKERHUB_MIRROR_AUTH" = "true" ]
        RUN --secret DOCKERHUB_MIRROR_USER --secret DOCKERHUB_MIRROR_PASS ./setup-registry.sh
    ELSE IF [ "$DOCKERHUB_AUTH" = "true" ]
        RUN --secret DOCKERHUB_USER --secret DOCKERHUB_PASS ./setup-registry.sh
    ELSE
        RUN ./setup-registry.sh
    END

    # pull out buildkit_additional_config from the earthly config, for the special case of earthly-in-earthly testing
    # which runs earthly-entrypoint.sh, which calls buildkitd/entrypoint, which requires EARTHLY_VERSION_FLAG_OVERRIDES to be set
    # NOTE: yq will print out `null` if the key does not exist, this will cause a literal null to be inserted into /etc/buildkit.toml, which will
    # cause buildkit to crash -- this is why we first assign it to a tmp variable, followed by an if.
    ENV EARTHLY_ADDITIONAL_BUILDKIT_CONFIG="$(export tmp=$(cat /etc/.earthly/config.yml | yq .global.buildkit_additional_config); if [ "$tmp" != "null" ]; then echo "$tmp"; fi)"

# prerelease builds and pushes the prerelease version of earthly.
# Tagged as prerelease
prerelease:
    FROM alpine:3.18
    ARG BUILDKIT_PROJECT
    BUILD \
        --platform=linux/amd64 \
        --platform=linux/arm64 \
        ./buildkitd+buildkitd --TAG=prerelease  --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    COPY (+earthly-all/* --VERSION=prerelease --DEFAULT_INSTALLATION_NAME=earthly) ./
    SAVE IMAGE --push earthly/earthlybinaries:prerelease

# prerelease-script copies the earthly folder and saves it as an artifact
prerelease-script:
    FROM alpine:3.18
    COPY ./earthly ./
    # This script is useful in other repos too.
    SAVE ARTIFACT ./earthly

# ci-release builds earthly for linux/amd64 in a container and pushes wtth the tag
# EARTHLY_GIT_HASH-TAG_SUFFIX Where TAG_SUFFIX must be provided
ci-release:
    # TODO: this was multiplatform, but that skyrocketed our build times. #2979
    # may help.
    FROM alpine:3.18
    ARG BUILDKIT_PROJECT
    ARG EARTHLY_GIT_HASH
    ARG --required TAG_SUFFIX
    BUILD \
        --platform=linux/amd64 \
        ./buildkitd+buildkitd --TAG=${EARTHLY_GIT_HASH}-${TAG_SUFFIX} --BUILDKIT_PROJECT="$BUILDKIT_PROJECT" --DOCKERHUB_BUILDKIT_IMG="buildkitd-staging"
    COPY (+earthly/earthly --DEFAULT_BUILDKITD_IMAGE="docker.io/earthly/buildkitd-staging:${EARTHLY_GIT_HASH}-${TAG_SUFFIX}" --VERSION=${EARTHLY_GIT_HASH}-${TAG_SUFFIX} --DEFAULT_INSTALLATION_NAME=earthly) ./earthly-linux-amd64
    SAVE IMAGE --push earthly/earthlybinaries:${EARTHLY_GIT_HASH}-${TAG_SUFFIX}

# dind builds both the alpine and ubuntu dind containers for earthly
dind:
    # OS_IMAGE is the base image to use, e.g. alpine, ubuntu
    ARG --required OS_IMAGE
    # OS_VERSION is the version of the base OS to use, e.g. 3.18.0, 23.04
    ARG --required OS_VERSION
    # DOCKER_VERSION is the version of docker to use, e.g. 20.10.14
    ARG --required DOCKER_VERSION
    FROM $OS_IMAGE:$OS_VERSION
    COPY ./buildkitd/docker-auto-install.sh /usr/local/bin/docker-auto-install.sh
    RUN docker-auto-install.sh
    LET DOCKER_VERSION_TAG=$DOCKER_VERSION
    IF [ "$OS_IMAGE" = "ubuntu" ]
        # the docker ce repo contains packages such as "5:24.0.4-1~ubuntu.20.04~focal", we will remove the the epoch and debian-revision values,
        # in order to display the upstream-version, e.g. "24.0.5-1".
        SET DOCKER_VERSION_TAG="$(echo $DOCKER_VERSION | sed 's/^[0-9]*:\([^~]*\).*$/\1/')"
        RUN if echo $DOCKER_VERSION_TAG | grep "[^0-9.-]"; then echo "DOCKER_VERSION_TAG looks bad; got $DOCKER_VERSION_TAG" && exit 1; fi
    END
    LET TAG=$OS_IMAGE-$OS_VERSION-docker-$DOCKER_VERSION_TAG
    ARG INCLUDE_TARGET_TAG_DOCKER=true
    IF [ "$INCLUDE_TARGET_TAG_DOCKER" = "true" ]
      ARG EARTHLY_TARGET_TAG_DOCKER
      SET TAG=$TAG-$EARTHLY_TARGET_TAG_DOCKER
    END
    ARG DOCKERHUB_USER=earthly
    IF [ "$LATEST" = "true" ]
      # latest means the version is ommitted (for historical reasons we initially just called it earthly/dind:alpine-3.18 or earthly/dind:ubuntu)
      SAVE IMAGE --push --cache-from=earthly/dind:$OS_IMAGE-main $DOCKERHUB_USER/dind:$OS_IMAGE
    END
    ARG DATETIME="$(date --utc +%Y%m%d%H%M%S)" # note this must be overriden when building a multi-platform image (otherwise the values wont match)
    SAVE IMAGE --push --cache-from=earthly/dind:$OS_IMAGE-main $DOCKERHUB_USER/dind:$TAG "$DOCKERHUB_USER/dind:$TAG-$DATETIME"


# for-own builds earthly-buildkitd and the earthly CLI for the current system
# and saves the final CLI binary locally at ./build/own/earthly
for-own:
    ARG BUILDKIT_PROJECT
    # GO_GCFLAGS may be used to set the -gcflags parameter to 'go build'. See
    # the documentation on +earthly for extra detail about this option.
    ARG GO_GCFLAGS
    BUILD ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    COPY (+earthly/earthly --GO_GCFLAGS="${GO_GCFLAGS}") ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/own/earthly

# for-linux builds earthly-buildkitd and the earthly CLI for the a linux amd64 system
# and saves the final CLI binary locally in the ./build/linux folder.
for-linux:
    ARG BUILDKIT_PROJECT
    ARG GO_GCFLAGS
    BUILD --platform=linux/amd64 ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    BUILD ./ast/parser+parser
    COPY (+earthly-linux-amd64/earthly -GO_GCFLAGS="${GO_GCFLAGS}") ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/linux/amd64/earthly

# for-linux-arm64 builds earthly-buildkitd and the earthly CLI for the a linux arm64 system
# and saves the final CLI binary locally in the ./build/linux folder.
for-linux-arm64:
    ARG BUILDKIT_PROJECT
    ARG GO_GCFLAGS
    BUILD --platform=linux/arm64 ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    BUILD ./ast/parser+parser
    COPY (+earthly-linux-arm64/earthly -GO_GCFLAGS="${GO_GCFLAGS}") ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/linux/arm64/earthly

# for-darwin builds earthly-buildkitd and the earthly CLI for the a darwin amd64 system
# and saves the final CLI binary locally in the ./build/darwin folder.
# For arm64 use +for-darwin-m1
for-darwin:
    ARG BUILDKIT_PROJECT
    ARG GO_GCFLAGS
    BUILD --platform=linux/amd64 ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    BUILD ./ast/parser+parser
    COPY (+earthly-darwin-amd64/earthly -GO_GCFLAGS="${GO_GCFLAGS}") ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/darwin/amd64/earthly

# for-darwin-m1 builds earthly-buildkitd and the earthly CLI for the a darwin m1 system
# and saves the final CLI binary locally.
for-darwin-m1:
    ARG BUILDKIT_PROJECT
    ARG GO_GCFLAGS
    BUILD --platform=linux/arm64 ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    BUILD ./ast/parser+parser
    COPY (+earthly-darwin-arm64/earthly -GO_GCFLAGS="${GO_GCFLAGS}") ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/darwin/arm64/earthly

# for-windows builds earthly-buildkitd and the earthly CLI for the a windows system
# and saves the final CLI binary locally in the ./build/windows folder.
for-windows:
    ARG GO_GCFLAGS
    # BUILD --platform=linux/amd64 ./buildkitd+buildkitd
    BUILD ./ast/parser+parser
    COPY (+earthly-windows-amd64/earthly.exe -GO_GCFLAGS="${GO_GCFLAGS}") ./
    SAVE ARTIFACT ./earthly.exe AS LOCAL ./build/windows/amd64/earthly.exe

# all-buildkitd builds buildkitd for both linux amd64 and linux arm64
all-buildkitd:
    ARG BUILDKIT_PROJECT
    BUILD \
        --platform=linux/amd64 \
        --platform=linux/arm64 \
        ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"

dind-alpine:
    BUILD +dind --OS_IMAGE=alpine --OS_VERSION=3.18 --DOCKER_VERSION=23.0.6-r5

dind-ubuntu:
    BUILD +dind --OS_IMAGE=ubuntu --OS_VERSION=20.04 --DOCKER_VERSION=5:24.0.5-1~ubuntu.20.04~focal
    BUILD +dind --OS_IMAGE=ubuntu --OS_VERSION=23.04 --DOCKER_VERSION=5:24.0.5-1~ubuntu.23.04~lunar

# all-dind builds alpine and ubuntu dind containers for both linux amd64 and linux arm64
all-dind:
    RUN --no-cache date --utc +%Y%m%d%H%M%S > datetime
    ARG DATETIME="$(cat datetime)"
    BUILD \
        --platform=linux/amd64 \
        --platform=linux/arm64 \
        +dind-alpine --DATETIME=$DATETIME
    BUILD \
        --platform=linux/amd64 \
        --platform=linux/arm64 \
        +dind-ubuntu --DATETIME=$DATETIME

# all builds all of the following:
# - Buildkitd for both linux amd64 and linux arm64
# - Earthly for all supported environments linux amd64 and linux arm64, Darwin amd64 and arm64, and Windos amd64
# - Earthly as a container image
# - Prerelease version of earthly as a container image
# - Dind alpine and ubuntu for both linux amd64 and linux arm64 as container images
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

# lint-docs runs lint against changelog and checks that line endings are unix style and files end
# with a single newline.
lint-docs:
    BUILD +lint-newline-ending
    BUILD +lint-changelog

# test-no-qemu runs tests without qemu virtualization by passing in dockerhub authentication and 
# using secure docker hub mirror configurations
test-no-qemu:
    BUILD --pass-args +test-quick
    BUILD --pass-args +test-no-qemu-quick
    BUILD --pass-args +test-no-qemu-normal
    BUILD --pass-args +test-no-qemu-slow

# test-quick runs the unit, chaos, offline, and go tests and ensures the earthly script does not write to stdout
test-quick:
    BUILD +unit-test
    BUILD +chaos-test
    BUILD +offline-test
    BUILD +earthly-script-no-stdout
    BUILD --pass-args ./ast/tests+all

# test-no-qemu-quick runs the tests from ./tests+ga-no-qemu-quick
test-no-qemu-quick:
    BUILD --pass-args ./tests+ga-no-qemu-quick \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-quick runs the tests from ./tests+ga-no-qemu-normal
test-no-qemu-normal:
    BUILD --pass-args ./tests+ga-no-qemu-normal \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-quick runs the tests from ./tests+ga-no-qemu-slow
test-no-qemu-slow:
    BUILD --pass-args ./tests+ga-no-qemu-slow \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-quick runs the tests from ./tests+ga-no-qemu-slow
test-no-qemu-kind:
    BUILD --pass-args ./tests+ga-no-qemu-kind \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-quick runs the tests from ./tests+ga-qemu
test-qemu:
    ARG GLOBAL_WAIT_END="false"
    BUILD --pass-args ./tests+ga-qemu \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test runs both no-qemu tests and qemu tests
test:
    BUILD --pass-args +test-no-qemu
    BUILD --pass-args +test-qemu

# test runs examples, no-qemu, qemu, and experimental tests
test-all:
    BUILD +examples
    BUILD --pass-args +test-no-qemu
    BUILD --pass-args +test-qemu
    BUILD --pass-args ./tests+experimental

# examples runs both sets of examples
examples:
    BUILD +examples1
    BUILD +examples2

# examples1 runs set 1 of examples
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

# examples2 runs set 2 of examples
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
    BUILD ./examples/clojure+docker
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
    BUILD ./examples/zig+docker

# license copies the license file and saves it as an artifact
license:
    COPY LICENSE ./
    SAVE ARTIFACT LICENSE

# npm-update-all helps keep all node package-lock.json files up to date.
npm-update-all:
    FROM node:16.16.0-alpine3.15
    COPY . /code
    WORKDIR /code
    FOR nodepath IN \
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
