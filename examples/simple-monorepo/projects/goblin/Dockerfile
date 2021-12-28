FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN go mod init goblin

RUN go build -o /usr/bin/app .
EXPOSE 8000
CMD ["app"]