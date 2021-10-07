# Import example

The `IMPORT` command can be used to alias a project reference, allowing it to be reused without duplicating the path.

For example:

```Dockerfile
init:
    DO ./some/local/path+PRINT --string="Foo Bar"
    COPY ./some/local/path+PRINT/log.txt ./
    RUN cat log.txt
    BUILD github.com/earthly/hello-world:main+hello
```

can be refactored as:

```Dockerfile
IMPORT ./some/local/path AS lib
IMPORT github.com/earthly/hello-world:main

init:
    DO lib+PRINT --string="Foo Bar"
    COPY lib+PRINT/log.txt ./
    RUN cat log.txt
    BUILD hello-world+hello
```

Note that the `IMPORT` command only supports project references, not target references.

To run this example, execute:

```bash
earthly +init
```
