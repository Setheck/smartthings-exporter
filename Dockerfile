FROM golang:alpine AS builder

ENV TZ="America/Los_Angeles"
RUN apk update && apk add build-base make git
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make build
WORKDIR /app
RUN cp /build/smartthings-exporter ./smartthings-exporter

FROM alpine

ENV TZ="America/Los_Angeles"
RUN mkdir -p /app /data

COPY --chown=65534:0 --from=builder /app /app
USER 65534
EXPOSE 8080

WORKDIR /data
ENTRYPOINT ["/app/smartthings-exporter"]
