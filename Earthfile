VERSION 0.8

earthly:
    FROM alpine
    COPY medium-file a-file
    RUN sleep 4
    RUN mkdir build
    RUN echo data3 > build/earthly
    SAVE ARTIFACT build/earthly

dummy:
    FROM alpine
    COPY +earthly/earthly /usr/bin/earthly

breakit:
   FROM +dummy
   IF false
     RUN echo false
   ELSE
     RUN echo true
   END
