FROM golang:1.23.1-alpine3.20 as builder
RUN apk update && apk add --no-cache git

LABEL authors="alturino"

WORKDIR /usr/app/url-shortener

COPY ["go.mod", "go.sum", "./pkg/", "./"]
RUN go mod tidy

RUN go build main.go

FROM alpine:3.20.3 as production
RUN apk add --no-cache dumb-init

WORKDIR /usr/src/url-shortener

RUN addgroup --system go && adduser -S -s /bin/false -G go go

COPY --chown=go:go --from=builder /usr/app/url-shortener/main .
COPY --chown=go:go ./pkg/ .

USER go
CMD [ "dumb-init", "./main" ]
