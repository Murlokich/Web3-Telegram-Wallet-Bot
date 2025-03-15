# syntax=docker/dockerfile:1

FROM golang:1.23-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o wallet ./cmd/app


FROM alpine:latest AS release-stage

# Required to verify Telegrams CA
RUN apk --no-cache add ca-certificates

COPY --from=build-stage /app/wallet /wallet
COPY --from=build-stage /app/migrations /migrations

ENTRYPOINT ["/wallet"]

