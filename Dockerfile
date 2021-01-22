FROM golang:1.13-buster AS builder

RUN mkdir /build
WORKDIR /build
ADD . /build

RUN go install github.com/swaggo/swag/cmd/swag
RUN swag init -p pascalcase -g "start.go" -o "./doc/api/"
RUN go build -o villasweb-backend

FROM debian:buster

RUN apt-get update && \
    apt-get install -y \
        ca-certificates

COPY --from=builder /build/villasweb-backend /usr/bin

EXPOSE 4000

CMD [ "villasweb-backend" ]

