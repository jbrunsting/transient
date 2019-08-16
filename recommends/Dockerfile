FROM golang:latest

RUN go get github.com/codegangsta/gin
RUN mkdir /src
WORKDIR /src

EXPOSE 4001

CMD gin -p 4001 -a 4000 -b ../gin-bin
