# syntax=docker/dockerfile:1

# ---------- build stage ----------
FROM golang:1.26-alpine AS build
WORKDIR /src

# Dependencies first, so this layer is cached until go.mod/go.sum change.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

# ---------- runtime stage ----------
FROM alpine:3.22
RUN apk add --no-cache ca-certificates tzdata \
 && addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=build /out/api /app/api
COPY config /app/config

USER app
EXPOSE 8080

ENTRYPOINT ["/app/api"]
CMD ["-config-file=/app/config/docker.yaml"]