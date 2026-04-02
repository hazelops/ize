FROM golang:1.17

ARG PROJECT_PATH=.
WORKDIR /go/src/app
COPY ${PROJECT_PATH}/main.go  ./

RUN go mod init cronjob

RUN go build -o /usr/bin/app .

CMD ["app"]
