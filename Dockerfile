# Build
FROM golang:1.18-alpine AS builder

WORKDIR /app

COPY . ./

RUN go build -o /movie-api

# Deploy
FROM alpine:3.16
COPY --from=builder /movie-api .
ENTRYPOINT [ "./movie-api" ]
