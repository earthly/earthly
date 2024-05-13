VERSION 0.8
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
    openssh \
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
    COPY --dir analytics autocomplete billing buildcontext builder logbus cleanup cloud cmd config conslogging debugger \
        dockertar docker2earthly domain features internal outmon slog states util variables regproxy ./
    COPY --dir buildkitd/buildkitd.go buildkitd/settings.go buildkitd/certificates.go buildkitd/
    COPY --dir earthfile2llb/*.go earthfile2llb/
    COPY --dir ast/antlrhandler ast/spec ast/hint ast/command ast/commandflag ast/*.go ast/
    COPY --dir inputgraph/*.go inputgraph/testdata inputgraph/

# update-buildkit updates earthly's buildkit dependency.
update-buildkit:
    FROM +code # if we use deps, go mod tidy will remove a bunch of requirements since it won't have access to our codebase.
    ARG BUILDKIT_GIT_SHA
    ARG BUILDKIT_GIT_BRANCH=earthly-main
    ARG BUILDKIT_GIT_ORG=earthly
    ARG BUILDKIT_GIT_REPO=buildkit
    COPY (./buildkitd+buildkit-sha/buildkit_sha --BUILDKIT_GIT_ORG="$BUILDKIT_GIT_ORG" --BUILDKIT_GIT_SHA="$BUILDKIT_GIT_SHA" --BUILDKIT_GIT_BRANCH="$BUILDKIT_GIT_BRANCH") buildkit_sha
    BUILD  ./buildkitd+update-buildkit-earthfile --BUILDKIT_GIT_ORG="$BUILDKIT_GIT_ORG" --BUILDKIT_GIT_SHA="$(cat buildkit_sha)" --BUILDKIT_GIT_REPO="$BUILDKIT_GIT_REPO"
    RUN --no-cache go mod edit -replace "github.com/moby/buildkit=github.com/$BUILDKIT_GIT_ORG/$BUILDKIT_GIT_REPO@$(cat buildkit_sha)"
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
    FROM earthly/dind:alpine-3.19-docker-25.0.5-r0
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
        for f in $(find . -not -path "./.git/*" -type f \( -iname '*.go' -o -iname 'Earthfile' -o -iname '*.earth' -o -iname '*.md' -o -iname '*.json' \) | grep -v "ast/tests/empty-targets.earth" ); do \
            if ! dos2unix < "$f" | cmp - "$f"; then \
                echo "$f contains windows-style newlines and must be converted to unix-style (use dos2unix to fix)"; \
                code=1; \
            fi; \
        done; \
        exit $code
    # test file ends with a single newline
    RUN set -e; \
        code=0; \
        for f in $(find . -not -path "./.git/*" -type f \( -iname '*.yml' -o -iname '*.go' -o -iname '*.sh' -o -iname '*.template' -o -iname 'Earthfile' -o -iname '*.earth' -o -iname '*.md' -o -iname '*.json' \) | grep -v "ast/tests/empty-targets.earth" | grep -v "tests/version/version-only.earth" | grep -v "examples/mkdocs" ); do \
            if [ "$(tail -c 1 $f)" != "$(printf '\n')" ]; then \
                echo "$f does not end with a newline"; \
                code=1; \
            fi; \
        done; \
        exit $code
    RUN export f=ast/tests/empty-targets.earth && \
    if [ "$(tail -c 1 $f)" = "$(printf '\n')" ]; then \
            echo "$f is a special-case test which must not end with a newline."; \
            exit 1; \
        fi
    # check for files with trailing newlines
    RUN set -e; \
        code=0; \
        for f in $(find . -not -path "./.git/*" -type f \( -iname '*.go' -o -iname 'Earthfile' -o -iname '*.earth' -o -iname '*.md' -o -iname '*.json' \) | grep -v "ast/tests/empty-targets.earth" | grep -v "ast/parser/earth_parser.go" | grep -v "ast/parser/earth_lexer.go" ); do \
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

unit-test-parser:
    FROM +deps
    COPY scripts/unit-test-parser/main.go .
    RUN go build -o testparser main.go
    SAVE ARTIFACT testparser

# unit-test runs unit tests (and some integration tests).
unit-test:
    FROM +code
    RUN apk add --no-cache --update podman fuse-overlayfs
    COPY +unit-test-parser/testparser .
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

# offline-test runs offline tests with network set to none
offline-test:
    FROM +code
    RUN --network=none (go test -run TestOffline ./cloud || kill $$) | tee test.log
    RUN if grep 'no tests to run' test.log; then echo "error: no test found" && exit 1; fi

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
    SAVE ARTIFACT CHANGELOG.md

changelog-parser:
    FROM python:3
    RUN pip install packaging
    COPY release/changelogparser.py /usr/bin/changelogparser
    WORKDIR /changelog
    COPY CHANGELOG.md .

lint-changelog:
    FROM +changelog-parser
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
    #     --GO_GCFLAGS='all=-N -l'
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
    ARG PUSH_LATEST_TAG="false"
    ARG PUSH_PRERELEASE_TAG="false"
    FROM ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT" --TAG="$TAG"
    RUN apk add --update --no-cache docker-cli libcap-ng-utils git
    ENV EARTHLY_IMAGE=true
    # When Earthly is run from a container, the registry proxy networking setup
    # will fail as the registry is meant to be run on a dynamic localhost port
    # (which won't be exposed by the container). Let's fall back to tar-based
    # image transfer until this can be addressed further.
    ENV EARTHLY_DISABLE_REMOTE_REGISTRY_PROXY=true
    COPY earthly-entrypoint.sh /usr/bin/earthly-entrypoint.sh
    ENTRYPOINT ["/usr/bin/earthly-entrypoint.sh"]
    WORKDIR /workspace
    COPY (+earthly/earthly --VERSION=$TAG --DEFAULT_INSTALLATION_NAME="earthly") /usr/bin/earthly
    ARG DOCKERHUB_USER="earthly"
    ARG DOCKERHUB_IMG="earthly"
    # Multiple SAVE IMAGE's lead to differing image digests, but multiple
    # arguments to the save SAVE IMAGE do not. Using variables here doesn't work
    # either, unfortunately, as the names are quoted and treated as a single arg.
    IF [ "$PUSH_LATEST_TAG" == "true" ]
       SAVE IMAGE --push --cache-from=earthly/earthly:main $DOCKERHUB_USER/$DOCKERHUB_IMG:$TAG $DOCKERHUB_USER/$DOCKERHUB_IMG:latest
    ELSE IF [ "$PUSH_PRERELEASE_TAG" == "true" ]
       SAVE IMAGE --push --cache-from=earthly/earthly:main $DOCKERHUB_USER/$DOCKERHUB_IMG:$TAG $DOCKERHUB_USER/$DOCKERHUB_IMG:prerelease
    ELSE
       SAVE IMAGE --push --cache-from=earthly/earthly:main $DOCKERHUB_USER/$DOCKERHUB_IMG:$TAG
    END

# earthly-integration-test-base builds earthly docker and then
# if no dockerhub mirror is not set it will attempt to login to dockerhub using the provided docker hub username and token.
# Otherwise, it will attempt to login to the docker hub mirror using the provided username and password
earthly-integration-test-base:
    FROM --pass-args +earthly-docker
    RUN apk update && apk add pcre-tools curl python3 bash perl findutils expect yq && apk add --upgrade sed
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
    RUN rm ./setup-registry.sh

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
    COPY (+earthly-linux-amd64/earthly --GO_GCFLAGS="${GO_GCFLAGS}") ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/linux/amd64/earthly

# for-linux-arm64 builds earthly-buildkitd and the earthly CLI for the a linux arm64 system
# and saves the final CLI binary locally in the ./build/linux folder.
for-linux-arm64:
    ARG BUILDKIT_PROJECT
    ARG GO_GCFLAGS
    BUILD --platform=linux/arm64 ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    BUILD ./ast/parser+parser
    COPY (+earthly-linux-arm64/earthly --GO_GCFLAGS="${GO_GCFLAGS}") ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/linux/arm64/earthly

# for-darwin builds earthly-buildkitd and the earthly CLI for the a darwin amd64 system
# and saves the final CLI binary locally in the ./build/darwin folder.
# For arm64 use +for-darwin-m1
for-darwin:
    ARG BUILDKIT_PROJECT
    ARG GO_GCFLAGS
    BUILD --platform=linux/amd64 ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    BUILD ./ast/parser+parser
    COPY (+earthly-darwin-amd64/earthly --GO_GCFLAGS="${GO_GCFLAGS}") ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/darwin/amd64/earthly

# for-darwin-m1 builds earthly-buildkitd and the earthly CLI for the a darwin m1 system
# and saves the final CLI binary locally.
for-darwin-m1:
    ARG BUILDKIT_PROJECT
    ARG GO_GCFLAGS
    BUILD --platform=linux/arm64 ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"
    BUILD ./ast/parser+parser
    COPY (+earthly-darwin-arm64/earthly --GO_GCFLAGS="${GO_GCFLAGS}") ./
    SAVE ARTIFACT ./earthly AS LOCAL ./build/darwin/arm64/earthly

# for-windows builds earthly-buildkitd and the earthly CLI for the a windows system
# and saves the final CLI binary locally in the ./build/windows folder.
for-windows:
    ARG GO_GCFLAGS
    # BUILD --platform=linux/amd64 ./buildkitd+buildkitd
    BUILD ./ast/parser+parser
    COPY (+earthly-windows-amd64/earthly.exe --GO_GCFLAGS="${GO_GCFLAGS}") ./
    SAVE ARTIFACT ./earthly.exe AS LOCAL ./build/windows/amd64/earthly.exe

# all-buildkitd builds buildkitd for both linux amd64 and linux arm64
all-buildkitd:
    ARG BUILDKIT_PROJECT
    BUILD \
        --platform=linux/amd64 \
        --platform=linux/arm64 \
        ./buildkitd+buildkitd --BUILDKIT_PROJECT="$BUILDKIT_PROJECT"

# all builds all of the following:
# - Buildkitd for both linux amd64 and linux arm64
# - Earthly for all supported environments linux amd64 and linux arm64, Darwin amd64 and arm64, and Windos amd64
# - Earthly as a container image
# - Prerelease version of earthly as a container image
all:
    BUILD +all-buildkitd
    BUILD +earthly-all
    BUILD +earthly-docker
    BUILD +prerelease

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
    BUILD --pass-args +test-misc
    BUILD --pass-args +test-no-qemu-group1
    BUILD --pass-args +test-no-qemu-group2
    BUILD --pass-args +test-no-qemu-group3
    BUILD --pass-args +test-no-qemu-group4
    BUILD --pass-args +test-no-qemu-group5
    BUILD --pass-args +test-no-qemu-group6
    BUILD --pass-args +test-no-qemu-group7
    BUILD --pass-args +test-no-qemu-group8
    BUILD --pass-args +test-no-qemu-group9
    BUILD --pass-args +test-no-qemu-group10
    BUILD --pass-args +test-no-qemu-group11
    BUILD --pass-args +test-no-qemu-group12
    BUILD --pass-args +test-no-qemu-slow

# test-misc runs misc (non earthly-in-earthly) tests
test-misc:
    BUILD +test-misc-group1
    BUILD +test-misc-group2
    BUILD +test-misc-group3
    BUILD +test-ast

test-misc-group1:
    BUILD +unit-test

test-misc-group2:
    BUILD +offline-test

test-misc-group3:
    BUILD +earthly-script-no-stdout

test-ast:
    BUILD +test-ast-group1
    BUILD +test-ast-group2
    BUILD +test-ast-group3

test-ast-group1:
    BUILD --pass-args ./ast/tests+group1

test-ast-group2:
    BUILD --pass-args ./ast/tests+group2

test-ast-group3:
    BUILD --pass-args ./ast/tests+group3

# test-no-qemu-group1 runs the tests from ./tests+ga-no-qemu-group1
test-no-qemu-group1:
    BUILD --pass-args ./tests+ga-no-qemu-group1 \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-group2 runs the tests from ./tests+ga-no-qemu-group2
test-no-qemu-group2:
    BUILD --pass-args ./tests+ga-no-qemu-group2 \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-group3 runs the tests from ./tests+ga-no-qemu-group3
test-no-qemu-group3:
    BUILD --pass-args ./tests+ga-no-qemu-group3 \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-group4 runs the tests from ./tests+ga-no-qemu-group4
test-no-qemu-group4:
    BUILD --pass-args ./tests+ga-no-qemu-group4 \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-group5 runs the tests from ./tests+ga-no-qemu-group5
test-no-qemu-group5:
    BUILD --pass-args ./tests+ga-no-qemu-group5 \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-group6 runs the tests from ./tests+ga-no-qemu-group6
test-no-qemu-group6:
    BUILD --pass-args ./tests+ga-no-qemu-group6 \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-group7 runs the tests from ./tests+ga-no-qemu-group7
test-no-qemu-group7:
    BUILD --pass-args ./tests+ga-no-qemu-group7 \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-group8 runs the tests from ./tests+ga-no-qemu-group8
test-no-qemu-group8:
    BUILD --pass-args ./tests+ga-no-qemu-group8 \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-group9 runs the tests from ./tests+ga-no-qemu-group9
test-no-qemu-group9:
    BUILD --pass-args ./tests+ga-no-qemu-group9 \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-group10 runs the tests from ./tests+ga-no-qemu-group10
test-no-qemu-group10:
    BUILD --pass-args ./tests+ga-no-qemu-group10 \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-group11 runs the tests from ./tests+ga-no-qemu-group11
test-no-qemu-group11:
    BUILD --pass-args ./tests+ga-no-qemu-group11 \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-group12 runs the tests from ./tests+ga-no-qemu-group12
test-no-qemu-group12:
    BUILD --pass-args ./tests+ga-no-qemu-group12 \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-slow runs the tests from ./tests+ga-no-qemu-slow
test-no-qemu-slow:
    BUILD --pass-args ./tests+ga-no-qemu-slow \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-no-qemu-kind runs the tests from ./tests+ga-no-qemu-kind
test-no-qemu-kind:
    BUILD --pass-args ./tests+ga-no-qemu-kind \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test-qemu runs the tests from ./tests+ga-qemu
test-qemu:
    ARG GLOBAL_WAIT_END="false"
    BUILD --pass-args ./tests+ga-qemu \
        --GLOBAL_WAIT_END="$GLOBAL_WAIT_END"

# test runs both no-qemu tests and qemu tests
test:
    BUILD --pass-args +test-no-qemu
    BUILD --pass-args +test-qemu

# smoke-test is used by circleci, and aims to be a medium-weight test which covers some WITH DOCKER and multi-platform tests
smoke-test:
    BUILD ./tests/with-docker-kind+alpine-kind
    BUILD ./tests/platform+test

# test runs examples, no-qemu, qemu, and experimental tests
test-all:
    BUILD +examples
    BUILD --pass-args +test-no-qemu
    BUILD --pass-args +test-qemu
    BUILD --pass-args ./tests+experimental

examples:
    BUILD +examples-1
    BUILD +examples-2
    BUILD +examples-3

examples-1:
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

examples-2:
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

examples-3:
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

# merge-main-to-docs merges the main branch into docs-0.8
merge-main-to-docs:
    RUN git config --global user.name "littleredcorvette" && \
        git config --global user.email "littleredcorvette@users.noreply.github.com" && \
        git config --global url."git@github.com:".insteadOf "https://github.com/"

    ARG TARGETARCH
    # renovate: datasource=github-releases depName=cli/cli
    ARG gh_version=v2.49.2
    RUN curl -Lo ghlinux.tar.gz \
      https://github.com/cli/cli/releases/download/$gh_version/gh_${gh_version#v}_linux_${TARGETARCH}.tar.gz \
      && tar --strip-components=1 -xf ghlinux.tar.gz \
      && rm ghlinux.tar.gz && mv ./bin/gh /usr/local/bin/gh

    ARG git_repo="earthly/earthly"
    ARG git_url="git@github.com:$git_repo"
    ARG earthly_lib_version=3.0.1
    ARG SECRET_PATH=littleredcorvette-id_rsa
    DO --pass-args github.com/earthly/lib/utils/git:$earthly_lib_version+DEEP_CLONE \
        --GIT_URL=$git_url --SECRET_PATH=$SECRET_PATH

    ARG to_branch="docs-0.8"
    ARG from_branch="main"

    LET temp_pr_branch="soon-to-be-$to_branch"
    RUN --push --secret GH_TOKEN=littleredcorvette-github-token --mount=type=secret,id=littleredcorvette-id_rsa,mode=0400,target=/root/.ssh/id_rsa \
        # 1. checkout the docs branch and merge changes from main
         git checkout $to_branch && git pull origin $to_branch &&\
         git merge $from_branch && \
        # 2. create a new temp branch to for a PR (can't push directly to protected branch)
        git checkout -b $temp_pr_branch && git push -f origin $temp_pr_branch && \
        # 3. create a new PR and wait till checks complete (if checks fail, close the PR)
        gh pr create --title "Temp PR to merge $from_branch to $to_branch" --draft -B $to_branch \
        --body "Opened by +merge-main-to-docs" --repo $git_repo && \
        sleep 15 && \
        ( \
            timeout --signal=SIGINT 300 \
            gh run watch $(gh run list --commit $(git rev-parse HEAD) -w "Check Docs for Broken Links" --json "databaseId" --jq '.[]|.databaseId') --exit-status || \
            (gh pr close $temp_pr_branch --delete-branch && exit 1) \
        ) && \
        # 4. try to push the branch now that the PR checks have passed
        git checkout $to_branch && (git push || (gh pr close $temp_pr_branch --delete-branch && exit 1))

# check-broken-links checks for broken links in our docs website
check-broken-links:
    FROM node:20-alpine3.18
    RUN npm install broken-link-checker -g
    WORKDIR /report
    ARG ADDRESS=https://docs.earthly.dev
    ARG VERBOSE=false
    LET REPORT_FILE_NAME=report.txt
    LET BLC_COMMAND="blc $ADDRESS -rog --exclude https://twitter.com/EarthlyTech --exclude http://localhost:8080/"
    IF [ $VERBOSE = "true" ]
        RUN --no-cache $BLC_COMMAND |tee $REPORT_FILE_NAME
    ELSE
        RUN --no-cache $BLC_COMMAND &> $REPORT_FILE_NAME || true
    END
    LET RESULT=$(grep -qE '^├─BROKEN─' $REPORT_FILE_NAME; echo $?)
    LET NOCOLOR='\033[0m'
    LET RED='\033[0;31m'
    LET GREEN='\033[0;32m'
    IF [ $RESULT = "0" ]
        RUN --no-cache echo -e "${RED}Final Broken Links Report:${NOCOLOR}"
        RUN --no-cache grep --color=always -E '^(Getting links from|├─BROKEN─|Finished!|Elapsed)' $REPORT_FILE_NAME
        RUN exit 1
    ELSE
        RUN --no-cache echo -e "${GREEN}No Broken Links were found${NOCOLOR}"
    END

# open-pr-for-fork creates a new PR based on the given pr_number
open-pr-for-fork:
    RUN git config --global user.name "littleredcorvette" && \
        git config --global user.email "littleredcorvette@users.noreply.github.com" && \
        git config --global url."git@github.com:".insteadOf "https://github.com/"

    ARG TARGETARCH
    # renovate: datasource=github-releases depName=cli/cli
    ARG gh_version=v2.49.2
    RUN curl -Lo ghlinux.tar.gz \
      https://github.com/cli/cli/releases/download/$gh_version/gh_${gh_version#v}_linux_${TARGETARCH}.tar.gz \
      && tar --strip-components=1 -xf ghlinux.tar.gz \
      && rm ghlinux.tar.gz && mv ./bin/gh /usr/local/bin/gh

    ARG earthly_lib_version=3.0.1
    ARG SECRET_PATH=littleredcorvette-id_rsa
    ARG git_repo="earthly/earthly"
    LET git_url="git@github.com:$git_repo"
    DO --pass-args github.com/earthly/lib/utils/git:$earthly_lib_version+DEEP_CLONE \
        --GIT_URL=$git_url --SECRET_PATH=$SECRET_PATH

    ARG --required pr_number
    RUN --no-cache --mount=type=secret,id=$SECRET_PATH,mode=0400,target=/root/.ssh/id_rsa \
        --secret GH_TOKEN=littleredcorvette-github-token \
        gh pr checkout $pr_number --branch "test-pr-$pr_number" --repo $git_repo && \
        echo "checked out $(git rev-parse HEAD)" && \
        git merge origin/main && \
        git commit --allow-empty -m "please run the test" && \
        git push -f origin HEAD
    RUN --no-cache --secret GH_TOKEN=littleredcorvette-github-token echo $(gh pr list -H test-pr-$pr_number -B main --json number --jq '.[]|.number'|| "") > /tmp/result
    LET test_pr=$(cat /tmp/result)
    IF [[ -z $test_pr ]]
        RUN --no-cache --secret GH_TOKEN=littleredcorvette-github-token \
            gh pr create --title "Run tests for PR $pr_number" --draft \
            --body "Running tests for https://github.com/$git_repo/pull/$pr_number" --repo $git_repo
    ELSE
        RUN --no-cache echo A matching test PR for PR $pr_number already exists: https://github.com/$git_repo/pull/$test_pr
    END

check-broken-links-pr:
    FROM alpine/git
    WORKDIR /tmp
    RUN apk add github-cli
    ARG BRANCH
    ARG EARTHLY_GIT_BRANCH
    LET branch=$BRANCH
    IF [ -z $branch ]
        SET branch=$EARTHLY_GIT_BRANCH
    END
    RUN --secret GH_TOKEN=littleredcorvette-github-token gh pr checks $branch --repo earthly/earthly | grep GitBook|awk '{print $5}' > url
    ARG VERBOSE
    BUILD --pass-args +check-broken-links --ADDRESS=$(cat url)
