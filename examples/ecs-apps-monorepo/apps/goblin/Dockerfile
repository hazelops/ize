FROM golang:1.17

ENV \
    EXAMPLE_API_KEY=Default \
    EXAMPLE_SECRET=Default

ARG PROJECT_PATH=.
WORKDIR /go/src/app
COPY ${PROJECT_PATH}/main.go  ./

RUN go mod init goblin

RUN go build -o /usr/bin/app .

CMD ["app"]

EXPOSE 3000
