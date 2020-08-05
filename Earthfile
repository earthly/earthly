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
	COPY --dir buildcontext builder cleanup cmd config conslogging debugger dockertar \
		domain llbutil logging ./
	COPY --dir buildkitd/buildkitd.go buildkitd/settings.go buildkitd/
	COPY --dir earthfile2llb/antlrhandler earthfile2llb/dedup earthfile2llb/image \
		earthfile2llb/imr earthfile2llb/variables earthfile2llb/*.go earthfile2llb/
	COPY ./earthfile2llb/parser+parser/*.go ./earthfile2llb/parser/
	SAVE IMAGE

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

unittest:
	FROM +code
	RUN go test ./...

buildkitd:
	BUILD ./buildkitd+buildkitd

debugger:
	FROM +code
	ARG GOCACHE=/go-cache
	ARG EARTHLY_TARGET_TAG
	ARG VERSION=$EARTHLY_TARGET_TAG
	ARG EARTHLY_GIT_HASH
	RUN --mount=type=cache,target=$GOCACHE \
		go build \
			-ldflags "-d -X main.Version=$VERSION -X main.GitSha=$EARTHLY_GIT_HASH $GO_EXTRA_LDFLAGS" \
			-tags netgo -installsuffix netgo \
			-o build/earth_debugger \
			cmd/debugger/*.go
	SAVE ARTIFACT build/earth_debugger

debugger-docker:
	# cant be FROM scratch because the args require sh to exist
	FROM busybox:1.32.0
	COPY +debugger/earth_debugger /earth_debugger
	ARG EARTHLY_TARGET_TAG
	ARG TAG=$EARTHLY_TARGET_TAG
	SAVE IMAGE --push earthly/debugger:$TAG

earth:
	FROM +code
	ARG GOOS=linux
	ARG GOARCH=amd64
	ARG GO_EXTRA_LDFLAGS="-linkmode external -extldflags -static"
	RUN test -n "$GOOS" && test -n "$GOARCH"
	ARG EARTHLY_TARGET_TAG
	ARG VERSION=$EARTHLY_TARGET_TAG
	ARG EARTHLY_GIT_HASH
	ARG DEFAULT_BUILDKITD_IMAGE=earthly/buildkitd:$VERSION
	ARG DEFAULT_DEBUGGER_IMAGE=earthly/debugger:debugger-logging@sha256:226b55f8d0a5f8fce684b4c7f114d4781e560089962686b441665e92c64d1c5a
	ARG GOCACHE=/go-cache
	RUN --mount=type=cache,target=$GOCACHE \
		go build \
			-ldflags "-X main.DefaultBuildkitdImage=$DEFAULT_BUILDKITD_IMAGE -X main.DefaultDebuggerImage=$DEFAULT_DEBUGGER_IMAGE -X main.Version=$VERSION -X main.GitSha=$EARTHLY_GIT_HASH $GO_EXTRA_LDFLAGS" \
			-o build/earth \
			cmd/earth/*.go
	SAVE ARTIFACT build/earth AS LOCAL "build/$GOOS/$GOARCH/earth"

earth-darwin:
	BUILD \
		--build-arg GOOS=darwin \
		--build-arg GOARCH=amd64 \
		--build-arg GO_EXTRA_LDFLAGS= \
		+earth

earth-all:
	BUILD +earth
	BUILD +earth-darwin

earth-docker:
	FROM ./buildkitd+buildkitd
	RUN apk add --update --no-cache docker-cli
	ENV ENABLE_LOOP_DEVICE=false
	ENV FORCE_LOOP_DEVICE=false
	COPY earth-buildkitd-wrapper.sh /usr/bin/earth-buildkitd-wrapper.sh
	ENTRYPOINT ["/usr/bin/earth-buildkitd-wrapper.sh"]
	ARG EARTHLY_TARGET_TAG
	ARG TAG=$EARTHLY_TARGET_TAG
	COPY --build-arg VERSION=$TAG +earth/earth /usr/bin/earth
	SAVE IMAGE --push earthly/earth:$TAG

for-linux:
	BUILD +buildkitd
	BUILD +earth

for-darwin:
	BUILD +buildkitd
	BUILD +earth-darwin

all:
	BUILD +buildkitd
	BUILD +debugger-docker
	BUILD +earth-all
	BUILD +earth-docker

test:
	BUILD +lint
	BUILD +unittest
	BUILD ./examples/tests+all
	BUILD +examples

examples:
	BUILD ./examples/go+docker
	BUILD ./examples/java+docker
	BUILD ./examples/js+docker
	BUILD ./examples/cpp+docker
	BUILD ./examples/scala+docker
	BUILD ./examples/dotnet+docker
	BUILD ./examples/monorepo+all
	BUILD ./examples/multirepo+docker
	BUILD ./examples/readme/go1+all
	BUILD ./examples/readme/go2+all
	BUILD ./examples/readme/go3+build
	BUILD ./examples/readme/proto+docker
	BUILD github.com/earthly/hello-world+hello

test-experimental:
	BUILD ./examples/tests+experimental
