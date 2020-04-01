# JS example

A complete JavaScript example is available on [the Basics page](../guides/basics.md).

```Dockerfile
# build.earth

FROM node:13.10.1-alpine3.11
WORKDIR /js-example

deps:
    COPY package.json package-lock.json ./
    RUN npm install
    SAVE ARTIFACT package-lock.json AS LOCAL ./package-lock.json
    SAVE IMAGE

build:
    FROM +deps
    COPY src src
    COPY dist dist
    RUN npx webpack
    SAVE ARTIFACT dist /dist AS LOCAL dist

docker:
    FROM +deps
    COPY +build/dist ./dist
    EXPOSE 8080
    ENTRYPOINT ["/js-example/node_modules/http-server/bin/http-server", "./dist"]
    SAVE IMAGE js-example:latest
```

For the complete code see the [examples/js GitHub directory](https://github.com/vladaionescu/earthly/tree/master/examples/js).
