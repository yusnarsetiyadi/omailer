FROM golang:1.25.0-alpine AS build

RUN mkdir /app

WORKDIR /app

COPY ./ /app

RUN go mod tidy

RUN go build -o omailer

EXPOSE 80

CMD [ "./omailer" ]