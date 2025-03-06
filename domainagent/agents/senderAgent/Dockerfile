FROM golang:1.22-alpine AS sender-builder
LABEL name="sender"

WORKDIR /app
COPY sender_agent.go .

RUN go mod init sender_agent && \
    go build -o sender_agent

FROM alpine:latest
WORKDIR /app
COPY --from=sender-builder /app/sender_agent .

CMD ["./sender_agent"]