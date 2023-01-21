FROM golang:alpine as go

RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /app
COPY . /app

RUN go build -v cmd/api/api.go

EXPOSE 80

CMD /app/api