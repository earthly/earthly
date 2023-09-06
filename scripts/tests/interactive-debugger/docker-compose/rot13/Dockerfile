FROM alpine:3.18
RUN apk add --no-cache nmap-ncat

# rot13 echo server
CMD ncat -l $LISTEN_PORT -c 'xargs -n1 -I {} sh -c "echo {} | tr \"A-Za-z\" \"N-ZA-Mn-za-m\""'
