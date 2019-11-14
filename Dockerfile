FROM golang:1.12.9-buster AS builder

RUN mkdir /build
WORKDIR /build
ADD . /build

RUN go install github.com/swaggo/swag/cmd/swag
RUN swag init -p pascalcase -g "start.go" -o "./doc/api/"
RUN go build -o villasweb-backend

FROM debian:buster

COPY --from=builder /build/villasweb-backend /usr/bin

EXPOSE 4000

CMD [ "villasweb-backend" ]

