FROM golang:1.18-buster AS builder

RUN mkdir /build
WORKDIR /build

# Make use of layer caching
ADD go.* ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag

ADD . .

RUN swag init --propertyStrategy pascalcase \
              --generalInfo "start.go" \
              --output "./doc/api/" \
              --parseDependency \
              --parseInternal \
              --parseVendor \
              --parseDepth 2

RUN go build -o villasweb-backend

FROM debian:buster

RUN apt-get update && \
    apt-get install -y \
        ca-certificates

COPY --from=builder /build/villasweb-backend /usr/bin

EXPOSE 4000

CMD [ "villasweb-backend" ]

