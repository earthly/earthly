# Import example

The `IMPORT` command can be used to alias an Earthfile reference, allowing it to be reused without duplicating the path.

For example:

```Dockerfile
build:
    DO ./some/local/path+PRINT --string="IMPORT example"
    COPY ./some/local/path+get-file/file.txt ./
    RUN cat file.txt
    BUILD github.com/earthly/hello-world:main+hello
```

can be refactored as:

```Dockerfile
IMPORT ./some/local/path AS lib
IMPORT github.com/earthly/hello-world:main

build:
    DO lib+PRINT --string="IMPORT example"
    COPY lib+get-file/file.txt ./
    RUN cat file.txt
    BUILD hello-world+hello
```

Note that the `IMPORT` command only supports Earthfile references, not target references.

To run this example, execute:

```bash
earthly +build
```
