VERSION 0.8

code:
    FROM alpine
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

# earthly builds the earthly CLI and docker image.
earthly:
    FROM +code
    RUN sleep 4
    RUN mkdir build
    RUN echo data3 > build/earthly
    SAVE ARTIFACT build/earthly

dummy:
    FROM alpine
    COPY +earthly/earthly /usr/bin/earthly

multi:
   FROM +dummy
   ARG i
   IF [ "$i" -gt 2 ]
     RUN echo $i is big
   ELSE
     RUN echo $i is small
   END

breakit:
   BUILD +multi --i=0 --i=1 --i=2 --i=3 --i=4
