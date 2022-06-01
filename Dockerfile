# Build image
FROM golang:1.17-alpine3.15 AS build

WORKDIR /go/src/github.com/shopsmart/ssm2ssm

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build executable
COPY . ./
RUN go build -o /go/bin/ssm2ssm ./cmd/ssm2ssm

# Main image
FROM alpine:3.15

COPY --from=build /go/bin/ssm2ssm /usr/local/bin/ssm2ssm

CMD [ "ssm2ssm" ]
