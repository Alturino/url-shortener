FROM golang:1.23.1-alpine3.20 AS builder
RUN apk update && apk add --no-cache git

LABEL authors="alturino"

WORKDIR /usr/app/url-shortener/

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY main.go  ./
COPY ./internal/ ./internal/

RUN go build main.go

FROM alpine:3.20.3 AS production
RUN apk add --no-cache dumb-init

WORKDIR /usr/app/url-shortener/

RUN addgroup --system go && adduser -S -s /bin/false -G go go

COPY --chown=go:go --from=builder /usr/app/url-shortener/main .
COPY --chown=go:go application.yaml .
COPY --chown=go:go ./migrations/ ./migrations/

RUN touch url_shortener.jsonl && chown -R go:go url_shortener.jsonl

USER go
CMD [ "dumb-init", "./main" ]
