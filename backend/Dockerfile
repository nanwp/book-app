FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/api-service ./cmd

FROM alpine:3.20
RUN apk add --no-cache tzdata ca-certificates && update-ca-certificates
WORKDIR /app
COPY --from=builder /out/api-service ./api-service
ENV TZ=Asia/Jakarta
USER nobody
ENTRYPOINT ["./api-service"]
