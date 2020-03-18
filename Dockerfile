FROM alpine
RUN apk update && apk --no-cache add tzdata
ADD slack-delegate-bot /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/slack-delegate-bot"]
EXPOSE 8080
