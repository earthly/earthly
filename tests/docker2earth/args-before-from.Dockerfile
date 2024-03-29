ARG BASE
ARG GO_MAJOR
ARG GO_MINOR=16
ARG GO_VERSION="${GO_MAJOR}.${GO_MINOR}"
FROM "${BASE}:${GO_VERSION}"
WORKDIR /go/src/github.com/alexellis/href-counter/
RUN go get -d -v golang.org/x/net/html
COPY app.go .
RUN go mod init
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest as greet
WORKDIR /root/
RUN echo greetings > /root/hello.txt

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/alexellis/href-counter/app .
COPY --from=greet /root/hello.txt .
CMD ["./app"]
