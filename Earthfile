VERSION 0.8

earthly:
    FROM alpine
    COPY medium-file a-file
    RUN sleep 1
    SAVE ARTIFACT a-file

breakit:
   FROM alpine
   COPY +earthly/a-file /a/file
   IF false
     RUN echo false
   END
