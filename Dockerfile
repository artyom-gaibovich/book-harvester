# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o /out/iprbooks-dumper ./cmd/dumper

FROM alpine:3.20
WORKDIR /app
COPY --from=build /out/iprbooks-dumper /app/iprbooks-dumper

ENTRYPOINT ["/app/iprbooks-dumper"]
