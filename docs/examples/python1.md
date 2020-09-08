# Python example

A complete Python example is available on [the Basics page](../guides/basics.md).

```Dockerfile
FROM python:3
WORKDIR /code

deps:
    RUN pip install wheel
    COPY requirements.txt ./
    RUN pip wheel -r requirements.txt --wheel-dir=wheels
    SAVE IMAGE

build:
    FROM +deps
    COPY src src
    SAVE ARTIFACT src /src
    SAVE ARTIFACT wheels /wheels

docker:
    COPY +build/src src
    COPY +build/wheels wheels
    COPY requirements.txt ./
    RUN pip install --no-index --find-links=wheels -r requirements.txt
    ENTRYPOINT ["python3", "./src/hello.py"]
    SAVE IMAGE python-example:latest
```

For the complete code see the [examples/python GitHub directory](https://github.com/earthly/earthly/tree/master/examples/python).
