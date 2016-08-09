FROM golang:1.6
MAINTAINER sklirg

ENV APP_DIR=/go/src/github.com/sklirg/discord-bot
RUN mkdir -p $APP_DIR

WORKDIR $APP_DIR

COPY . $APP_DIR

RUN curl -s https://glide.sh/get | sh
RUN glide install
RUN go build

CMD ["./discord-bot"]
