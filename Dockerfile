FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct

WORKDIR $GOPATH/src/github.com/webbleen/go-gin

COPY . $GOPATH/src/github.com/webbleen/go-gin

RUN go build .

EXPOSE 8000
ENTRYPOINT ["./go-gin"]