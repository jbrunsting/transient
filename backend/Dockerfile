FROM golang:latest 

RUN mkdir /src
ADD . /src
WORKDIR /src

RUN go build -o main .

CMD ["/src/main"]
