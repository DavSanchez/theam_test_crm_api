FROM golang:1.14

WORKDIR /go/src/
COPY . .

RUN go get -d -v ./...

EXPOSE 4000

CMD ["go run main.go"]