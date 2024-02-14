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
