# Builder Stage
FROM golang:1.18-alpine AS builder

RUN apk update
RUN apk add --no-cache curl bash nano

ARG VERSION
ARG CGO_ENABLED

WORKDIR /app/src
COPY ../.. .

RUN go mod tidy
RUN go build -ldflags "-s -w -X main.version=${VERSION}" -o /app/task-queue-asynq ./cmd/client
RUN go build -ldflags "-s -w -X main.version=${VERSION}" -o /app/asynq-server ./cmd/server

WORKDIR /app

RUN rm -rf src/
ENTRYPOINT ["/app/task-queue-asynq"]

# Production Stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/task-queue-asynq ./app
COPY --from=builder /app/asynq-server ./server
CMD ["./app", "./server"]
