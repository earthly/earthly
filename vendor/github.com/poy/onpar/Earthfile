VERSION 0.7

# Make sure these are up to date
ARG goversion=1.20
ARG distro=alpine3.17

FROM golang:${goversion}-${distro}
WORKDIR /go-workdir

deps:
    # Copying only go.mod and go.sum at first allows earthly to cache the result of
    # `go mod download` unless go.mod or go.sum have changed.
    COPY go.mod go.sum .
    RUN go mod download

    # `go mod download` may have changed these files, so we save them locally when
    # `earthly +deps` is called directly.
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

test-base:
    FROM +deps
    # gcc and g++ are required for `go test -race`
    RUN apk add --update gcc g++
    COPY . .

project-base:
    FROM +deps
    COPY . .

# tidy runs 'go mod tidy' and updates go.mod and go.sum.
tidy:
    FROM +project-base
    RUN go mod tidy
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum

# vet runs 'go vet'.
vet:
    FROM +project-base

    # pkg is the package to run 'go vet' against.
    ARG pkg = ./...

    RUN go vet $path

# fmt formats source files.
fmt:
    FROM +project-base

    RUN go install golang.org/x/tools/cmd/goimports@latest
    RUN goimports -w .
    FOR --sep "\n" gofile IN $(find . -name \'*.go\')
      SAVE ARTIFACT $gofile AS LOCAL $gofile
    END

# generate re-generates code and saves the outputs locally.
generate:
    FROM +project-base

    RUN go install git.sr.ht/~nelsam/hel@latest
    RUN go install golang.org/x/tools/cmd/goimports@latest

    RUN go generate ./...

    # 'go install' can update our go.mod and go.sum, so save them too.
    SAVE ARTIFACT go.mod AS LOCAL go.mod
    SAVE ARTIFACT go.sum AS LOCAL go.sum
    FOR gofile IN $(find . -name helheim_test.go)
      SAVE ARTIFACT $gofile AS LOCAL $gofile
    END

# test runs tests.
test:
    FROM +test-base

    # pkg is the package to run tests against.
    ARG pkg = ./...

    RUN go test -race $pkg
