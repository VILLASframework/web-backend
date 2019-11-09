FROM golang:1.12.9-buster AS builder

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

FROM debian:buster

COPY --from=builder /go/bin/villasweb-backend-go /usr/bin

EXPOSE 4000

CMD [ "villasweb-backend-go" ]
