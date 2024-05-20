# +----------------------------------------------------+
# |Stage 1: Build frontend                             |
# +----------------------------------------------------+
FROM node:20-alpine3.19 AS build-fe

COPY ./html/package.json /tmp/package.json
RUN cd /tmp && yarn install
RUN mkdir /build && cp -a /tmp/node_modules /build/

# Build
WORKDIR /build
COPY ./html /build/

RUN yarn run build


# +----------------------------------------------------+
# |Stage 1: Build binary                               |
# +----------------------------------------------------+
FROM golang:1.21.10-alpine3.19 AS build

RUN apk update && apk add --no-cache git

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY --from=build-fe /build/static/dist ./html/static/dist
ENV GOCACHE=/root/.cache/go-build
RUN --mount=type=cache,target="/root/.cache/go-build" \
    go build -buildvcs -o /usr/bin/wadoh .


# +----------------------------------------------------+
# |Stage 2: Runner                                     |
# +----------------------------------------------------+
FROM ghcr.io/linuxserver/baseimage-alpine:3.17 as server

RUN apk --update add ca-certificates \
    mailcap \
    curl

ENV APP_ADDRESS=http://localhost

HEALTHCHECK --interval=5s --timeout=3s --start-period=5s --retries=2 \
    CMD curl -f ${APP_ADDRESS}/_/ping || exit 1

VOLUME /log /config
WORKDIR /log
COPY --from=build /usr/bin/wadoh /usr/bin/wadoh
EXPOSE 8080

ENTRYPOINT [ "/usr/bin/wadoh" ]
