FROM golang:1.18-alpine

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go get -d -v ./...

EXPOSE 9090

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

CMD ["./main"]