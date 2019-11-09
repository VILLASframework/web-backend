FROM golang:1.12.9-buster AS builder

RUN mkdir /build
WORKDIR /build
ADD . /build

RUN go build -o villasweb-backend

FROM debian:buster

COPY --from=builder /build/villasweb-backend /usr/bin

EXPOSE 4000

CMD [ "villasweb-backend" ]

