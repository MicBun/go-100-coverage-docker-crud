FROM golang:latest

WORKDIR WORKDIR /go/src/github.com/MicBun/go-microservice-kubernetes

COPY . .

RUN go get -d -v ./...

RUN go build -o go-microservice-kubernetes ./bin/web/main.go

EXPOSE 8080

ENTRYPOINT ["./go-microservice-kubernetes"]