FROM golang:1.19-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o slack-delegate-bot ./cmd/delegatebot

FROM alpine
RUN apk update && apk --no-cache add tzdata
COPY --from=builder /app/slack-delegate-bot /usr/local/bin/slack-delegate-bot
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/slack-delegate-bot"]
