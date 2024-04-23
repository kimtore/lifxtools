# vi: se ft=dockerfile:

FROM golang:1.17-alpine as builder
RUN apk add --no-cache git make
ENV GOOS=linux
ENV CGO_ENABLED=0
ENV GO111MODULE=on
COPY . /src
WORKDIR /src
RUN go mod download
RUN go build -o bin/lifxtool cmd/lifxtool/*.go

FROM alpine:3.15
EXPOSE 8080
WORKDIR /app
COPY --from=builder /src/bin/lifxtool /app/lifxtool
ENV BIND_ADDRESS=0.0.0.0:8080
ENV CONFIG_FILE=/config/config.yaml
CMD /app/lifxtool
