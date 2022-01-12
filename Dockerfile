FROM golang:1.17

RUN mkdir /etc/incus

COPY . /go/src/github.com/Salesflare/incus
WORKDIR /go/src/github.com/Salesflare/incus

RUN ./scripts/build.sh

CMD ["/go/bin/incus", "-conf='/etc/incus/'"]
