# Dockerfile for Receiver Agent
FROM golang:1.22-alpine AS receiver-builder
LABEL name="receiver"

WORKDIR /app
RUN pwd
COPY receiver_agent.go .

RUN go mod init receiver_agent && \
    go build -o receiver_agent

FROM alpine:latest
WORKDIR /app
COPY --from=receiver-builder /app/receiver_agent .

EXPOSE 8080
CMD ["./receiver_agent"]

