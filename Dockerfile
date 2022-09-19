FROM golang:1.16 as build-env

ENV CONFIG_PATH $CONFIG_PATH

WORKDIR /app

COPY ./ .

RUN go mod download

RUN CGO_ENABLED=0 go build -o /opt/proxy ./cmd/app/main.go

# Release
FROM alpine:latest

WORKDIR /root/

COPY --from=build-env /opt/proxy .
COPY --from=build-env /app/config ./config

CMD ["./proxy", "run", "$CONFIG_PATH"]