FROM golang:1.19

WORKDIR /usr/app

COPY go.mod .
RUN go mod download && go mod verify

COPY main.go .
RUN go build -v -o /usr/local/bin/app ./...

CMD ["app"]