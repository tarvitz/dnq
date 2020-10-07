FROM golang:1.15-alpine3.12 as builder
RUN set -x \
    && apk add --no-cache git make
COPY . /build/app/
WORKDIR /build/app
ARG APP_VERSION
ARG APP_BIN="dnq"
RUN set -x \
    && ./hack/up-sources.sh > sources.txt \
    && GOOS=linux GOARCH=amd64 \
        go build -ldflags "-X main.appVersion=${APP_VERSION} -w -s" \
                 -o /build/${APP_BIN} $(cat sources.txt) \
    && chmod +x /build/${APP_BIN}

FROM alpine:3.12
COPY --from=builder /build/${APP_BIN} /bin
