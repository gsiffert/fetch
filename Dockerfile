FROM golang:1.22-alpine3.19 as builder

# Because we use sqlite3 driver, we need to install gcc.
RUN apk update && apk add build-base

WORKDIR /app

COPY go.mod go.mod
COPY go.sum go.sum
COPY vendor vendor

COPY cmd cmd
COPY internal internal

RUN CGO_ENABLED=1 go build ./cmd/fetch

FROM alpine:3.19

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/fetch /fetch

ENTRYPOINT ["/fetch"]
