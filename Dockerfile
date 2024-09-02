FROM golang:1.21-bullseye AS builder

RUN mkdir /build
WORKDIR /build

# Make use of layer caching
ADD go.* ./
RUN go mod download
RUN go install github.com/swaggo/swag/cmd/swag@v1.8.3

ADD . .

RUN swag init --propertyStrategy pascalcase \
              --generalInfo "start.go" \
              --output "./doc/api/" \
              --parseDependency \
              --parseInternal \
              --parseVendor \
              --parseDepth 2

RUN go build -o villasweb-backend

FROM debian:bullseye-slim

RUN apt-get update && \
    apt-get install -y \
        ca-certificates && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/villasweb-backend /usr/bin

EXPOSE 4000

CMD [ "villasweb-backend" ]

