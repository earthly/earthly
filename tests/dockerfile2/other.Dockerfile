FROM alpine:3.18

COPY a.txt .
RUN cat a.txt
ENTRYPOINT ["cat", "a.txt"]
