VERSION 0.8

code:
    FROM alpine
    IF false
        RUN false
    END
    IF false
        RUN false
    END
    COPY --dir analytics autocomplete billing buildcontext builder logbus cleanup cloud cmd config conslogging debugger \
        dockertar docker2earthly domain features internal outmon slog states util variables regproxy ./
    COPY --dir buildkitd/buildkitd.go buildkitd/settings.go buildkitd/certificates.go buildkitd/
    COPY --dir earthfile2llb/*.go earthfile2llb/
    COPY --dir ast/antlrhandler ast/spec ast/hint ast/command ast/commandflag ast/*.go ast/
    COPY --dir inputgraph/*.go inputgraph/testdata inputgraph/

# earthly builds the earthly CLI and docker image.
earthly:
    FROM alpine
    COPY --dir analytics autocomplete billing buildcontext builder logbus cleanup cloud cmd config conslogging debugger \
        dockertar docker2earthly domain features internal outmon slog states util variables regproxy ./
    COPY --dir buildkitd/buildkitd.go buildkitd/settings.go buildkitd/certificates.go buildkitd/
    COPY --dir earthfile2llb/*.go earthfile2llb/
    COPY --dir ast/antlrhandler ast/spec ast/hint ast/command ast/commandflag ast/*.go ast/
    COPY --dir inputgraph/*.go inputgraph/testdata inputgraph/
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
