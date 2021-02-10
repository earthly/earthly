#!/usr/bin/env bash
set -eu

earthly=${earthly:=earthly}
earthly=$(realpath $earthly)
echo "running tests with $earthly"

# use host IP, otherwise earthly-buildkit won't be able to connect to it
ip=$(ifconfig eth0 | grep -w 'inet' | awk '{print $2}')
test -n "$ip"

# start up a new ssh-agent
eval "$(ssh-agent)"

# create a new key
ssh-keygen -b 3072 -t rsa -f /tmp/sshkey -q -N "" -C "testkey"
pubkey=$(cat /tmp/sshkey.pub)
ssh-add /tmp/sshkey

sudo /bin/sh -c "echo 127.0.0.1 ip4-localhost >> /etc/hosts"
sshhost="ip4-localhost"

# add test ssh server to known hosts
mkdir -p ~/.ssh
{
	ssh-keyscan -p "$SSH_PORT" -H $sshhost
	ssh-keyscan -p "$SSH_PORT" -H 127.0.0.1
	ssh-keyscan -p "$SSH_PORT" -H "$ip"
} > ~/.ssh/known_hosts

cat ~/.ssh/known_hosts

# setup passwordless login
sshpass -p "root" ssh root@$sshhost -p "$SSH_PORT" "/bin/sh -c \"echo $pubkey > /root/.ssh/authorized_keys\""

# setup a non-standard self-hosted git repo under the root user
ssh root@$sshhost -p "$SSH_PORT" "/bin/sh -c \"apt-get update && apt-get install -y git\""
ssh root@$sshhost -p "$SSH_PORT" "/bin/sh -c \"mkdir -p /root/my/really/weird/path/project.git\""
ssh root@$sshhost -p "$SSH_PORT" "/bin/sh -c \"cd /root/my/really/weird/path/project.git; git init --bare \""

# setup git
git config --global user.email "inigo@montoya.com"
git config --global user.name "my name is Inigo Montoya"

# create an Earthfile for our new private git repo
mkdir -p ~/odd-project
cd ~/odd-project
git init
cat <<EOF >> Earthfile

FROM alpine:latest

build:
  RUN echo -e "#!/bin/sh\necho hello weird world" > say-hi
  RUN chmod +x say-hi
  SAVE ARTIFACT say-hi

docker:
  COPY +build/say-hi /bin/say-hi
  RUN chmod +x /bin/say-hi
  ENTRYPOINT ["/bin/say-hi"]
  SAVE IMAGE weirdrepo:latest
EOF

git add Earthfile
git commit -m 'This is my weird commit'
git branch -M trunk
git remote add origin ssh://root@$sshhost:2222/root/my/really/weird/path/project.git
git push -u origin trunk

# delete the repo now that we've pushed it
cd ~
rm -rf odd-project

# test that earthly has access to it
$earthly config git "{myserver: {pattern: 'myserver/([^/]+)', substitute: 'ssh://root@$ip:2222/root/my/really/weird/path/\$1.git', auth: ssh}}"

if ! $earthly -V myserver/project:trunk+docker; then
    docker ps -a
    docker logs earthly-buildkitd
    exit 1
fi

# finally test that the container was built and runs
docker run --rm weirdrepo:latest | grep "hello weird world"
