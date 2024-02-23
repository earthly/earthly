# Debugging techniques

Traditional debugging of errors during image builds often require a developer to place various print
commands through out the build commands to help reason about the state of the system before the failure occurs.
This can be slow and cumbersome.

Earthly provides an interactive mode which gives you access to a root shell when an error occurs, which we'll
cover in this guide.

Let's consider a test example that prints out a randomly generated phrase:

```Dockerfile
# Earthfile

VERSION 0.8
FROM python:3
WORKDIR /code

test:
  RUN curl https://raw.githubusercontent.com/jsvine/markovify/master/test/texts/sherlock.txt > /sherlock.txt
  COPY generate_phrase.py .
  RUN pip3 install markovify
  RUN python3 generate_phrase.py
```

and our python code:
```Python
# generate_phrase.py

import markovify
text = open('sherlock.txt').read()
text_model = markovify.Text(text)
print(text_model.make_sentence())
```


Now we can run it with `earthly +test`, and we'll see a failure has occurred:

```
=========================== FAILURE ===========================
+test *failed* | --> RUN python3 generate_phrase.py
+test *failed* | Traceback (most recent call last):
+test *failed* |   File "generate_phrase.py", line 3, in <module>
+test *failed* |     text = open('sherlock.txt').read()
+test *failed* | FileNotFoundError: [Errno 2] No such file or directory: 'sherlock.txt'
+test *failed* | Command /bin/sh -c python3 generate_phrase.py failed with exit code 1
+test *failed* | +test *failed* | ERROR: Command exited with non-zero code: RUN python3 generate_phrase.py
Error: solve side effects: solve: failed to solve: rpc error: code = Unknown desc = executor failed running [/bin/sh -c  /usr/bin/earth_debugger /bin/sh -c 'python3 generate_phrase.py']: buildkit-runc did not terminate successfully
```

Why can't it find the sherlock.txt file? Let's re-run `earthly` with the `--interactive` (or `-i`) flag: `earthly -i +test`

This time we see a slightly different message:

```
+test | --> RUN python3 generate_phrase.py
+test | Traceback (most recent call last):
+test |   File "generate_phrase.py", line 3, in <module>
+test |     text = open('sherlock.txt').read()
+test | FileNotFoundError: [Errno 2] No such file or directory: 'sherlock.txt'
+test | Command /bin/sh -c python3 generate_phrase.py failed with exit code 1
+test | Entering interactive debugger (**Warning: only a single debugger per host is supported**)
+test | root@buildkitsandbox:/code#
```

This time rather than exiting, earthly will drop us into an interactive root shell within the container of the build environment.
This root shell will allow us to execute arbitrary commands within the container to figure out the problem:

```
root@buildkitsandbox:/code# ls
generate_phrase.py
root@buildkitsandbox:/code# find / | grep sherlock.txt
/sherlock.txt
root@buildkitsandbox:/code# ls /
bin  boot  code  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  sherlock.txt  srv	sys  tmp  usr  var
root@buildkitsandbox:/code# ls /sherlock.txt
/sherlock.txt
```

Ah ha! the corpus text file was located in the root directory rather than under `/code`. We can try moving it manually to see if that fixes the problem:

```
root@buildkitsandbox:/code# mv /sherlock.txt /code/.
root@buildkitsandbox:/code# python3 generate_phrase.py
I struck him down with the servants and with the lantern and left a fragment in the midst of my work during the last three years, although he has cruelly wronged.
```

At this point we know what needs to be done to fix the test, so we can type exit (or ctrl-D), to exit the interactive shell.

```
+test | time="2020-09-16T22:23:53Z" level=error msg="failed to read from ptmx: read /dev/ptmx: input/output error"
+test | time="2020-09-16T22:23:53Z" level=error msg="failed to read data from conn: read tcp 127.0.0.1:36672->127.0.0.1:5000: use of closed network connection"
+test | ERROR: Command exited with non-zero code: RUN python3 generate_phrase.py
```

Note that even though we fixed the problem during debugging, the image will not have been saved, so we must go back to our Earthfile and fix the problem there:

```Dockerfile
# Earthfile

VERSION 0.8
FROM python:3
WORKDIR /code

test:
  RUN curl https://raw.githubusercontent.com/jsvine/markovify/master/test/texts/sherlock.txt > /code/sherlock.txt
  COPY generate_phrase.py .
  RUN pip3 install markovify
  RUN python3 generate_phrase.py
```


## Debugging integration tests

Let's consider a more complicated example where we are running integration tests within an embedded docker setup:

```Dockerfile
# Earthfile

VERSION 0.8

server:
  COPY server.py .

test:
  FROM docker:19.03.12-dind
  RUN apk add curl
  WITH DOCKER --load server:latest=+server
    RUN docker run --rm -d --network=host server:latest python3 server.py && sleep 5 && curl -s localhost:8000 | grep hello
  END

```

and our server.py code:

```Python
from http.server import HTTPServer, BaseHTTPRequestHandler

class SimpleHTTPRequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(b'Hello, world!')

httpd = HTTPServer(('localhost', 8000), SimpleHTTPRequestHandler)
httpd.serve_forever()
```

Let's fire up our integration test with `earthly -P -i +test`:

```
buildkitd | Found buildkit daemon as docker container (earthly-buildkitd)
+base | --> FROM python:3
context | --> local context .
+base | resolve docker.io/library/python:3@sha256:e9b7e3b4e9569808066c5901b8a9ad315a9f14ae8d3949ece22ae339fff2cad0 100%
context | transferring .: 100%
+base | *cached* --> WORKDIR /code
+server | *cached* --> COPY server.py .
+test | --> FROM docker:19.03.12-dind
+test | resolve docker.io/library/docker:19.03.12-dind@sha256:674f1f40ff7c8ac14f5d8b6b28d8fb1f182647ff75304d018003f1e21a0d8771 100%
+test | *cached* --> RUN apk add curl
+test | --> WITH DOCKER RUN docker run --rm -d --network=host server:latest python3 server.py && sleep 5 && curl -s localhost:8000 | grep hello
+test | Loading images...
+test | Loaded image: server:latest
+test | ...done
+test | 1dc054c647cb75bde4897a2828edb095739cb9f864ed203ed2ddb54e62554aad
+test | Command /bin/sh -c docker run --rm -d --network=host server:latest python3 server.py && sleep 5 && curl -s localhost:8000 | grep hello failed with exit code 1
+test | Entering interactive debugger (**Warning: only a single debugger per host is supported**)
```



There was a failure checking that the server output contained the string `hello`; let's see what is going on:


```
/ # docker ps -a
CONTAINER ID        IMAGE               COMMAND               CREATED             STATUS              PORTS               NAMES
b8a31c54dd17        server:latest       "python3 server.py"   5 seconds ago       Up 4 seconds                            frosty_rhodes
```

The good news is our server container is running; let's see what happens when we try to connect to it:

```
/ # curl -s localhost:8000
Hello, world!/ 
```

Ah ha! The problem is our test is expecting a lowercase `h`, so we can fix our grep to look for an uppercase `H`:

```Dockerfile
# Earthfile

VERSION 0.8

server:
  COPY server.py .

test:
  FROM docker:19.03.12-dind
  RUN apk add curl
  WITH DOCKER --load server:latest=+server
    RUN docker run --rm -d --network=host server:latest python3 server.py && sleep 5 && curl -s localhost:8000 | grep Hello
  END
```

Then when we re-run our test we get:

```
+test | --> WITH DOCKER RUN docker run --rm -d --network=host server:latest python3 server.py && sleep 5 && curl -s localhost:8000 | grep Hello
+test | Loading images...
+test | Loaded image: server:latest
+test | ...done
+test | cb5299ae03cd17cfb2b528f01268ccf59761feec036cb313a3e969930d6f0815
+test | Hello, world!
+test | Target +test built successfully
=========================== SUCCESS ===========================
```

With the use of the interactive debugger; we were able to examine the state of the embedded containerized 

## Demo

[![asciicast](https://asciinema.org/a/361170.svg)](https://asciinema.org/a/361170?speed=2)

## Final tips

If you ever want to jump into an interactive debugging session at any point in your Earthfile, you can simply add a command that will fail such as:

```
  RUN false
```

and run earthly with the `--interactive` (or `-i`) flag.


Hopefully you won't run into failures, but if you do the interactive debugger may help you discover the root cause more easily. Happy coding.
