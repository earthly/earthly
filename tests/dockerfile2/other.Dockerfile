FROM alpine:3.16

COPY a.txt .
RUN cat a.txt
ENTRYPOINT ["cat", "a.txt"]
