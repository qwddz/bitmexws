FROM golang:alpine as go

RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /app
COPY . /app

RUN go build -v cmd/app/app.go

EXPOSE 8090

CMD /app/app