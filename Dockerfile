FROM golang:1.4.0

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN go get github.com/tools/godep github.com/codegangsta/gin 

RUN mkdir -p /etc/service/doorbot-api && mkdir -p /go/src/bitbucket.org/msamson/doorbot-api

WORKDIR /go/src/bitbucket.org/msamson/doorbot-api

COPY ./Godeps /go/src/bitbucket.org/msamson/doorbot-api/Godeps
RUN godep restore

COPY . /go/src/bitbucket.org/msamson/doorbot-api

RUN go-wrapper download && \
    go-wrapper install


ENTRYPOINT [ "gin" ]
EXPOSE 3000
