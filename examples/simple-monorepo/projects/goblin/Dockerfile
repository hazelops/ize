FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN go build -o /usr/bin/app src/main.go
EXPOSE 8000
CMD ["app"]