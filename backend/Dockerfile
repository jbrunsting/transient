FROM golang:latest

RUN go get github.com/codegangsta/gin
RUN mkdir /src
WORKDIR /src

EXPOSE 3001

CMD gin -p 3001 -a 3000 -b ../gin-bin
