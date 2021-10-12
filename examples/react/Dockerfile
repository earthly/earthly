FROM node:16.10.0-alpine as build

WORKDIR /app

COPY package.json ./

COPY yarn.lock ./

RUN yarn install --frozen-lockfile

COPY . ./

RUN yarn build

FROM nginx:stable-alpine 

COPY --from=build /app/build /usr/share/nginx/html
