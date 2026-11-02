# syntax=docker/dockerfile:1

# ---- build -------------------------------------------------------------
FROM golang:1.24-alpine AS build
WORKDIR /src

# Cache module downloads in their own layer.
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/abrayly-mebel-backend .

# ---- runtime -----------------------------------------------------------
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app -u 10001 app \
    && mkdir -p /app/uploads \
    && chown -R app:app /app

WORKDIR /app
COPY --from=build /out/abrayly-mebel-backend /app/abrayly-mebel-backend

USER app

ENV GO_ENV=production \
    PORT=3030

EXPOSE 3030

# Uploaded category, product and banner images.
VOLUME ["/app/uploads"]

HEALTHCHECK --interval=30s --timeout=3s --start-period=15s --retries=3 \
    CMD wget -qO- "http://127.0.0.1:${PORT}/health" > /dev/null || exit 1

ENTRYPOINT ["/app/abrayly-mebel-backend"]
