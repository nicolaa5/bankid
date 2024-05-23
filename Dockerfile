# Start from a base Golang image
FROM golang:1.21-alpine

LABEL maintainer="Nicolaas Bijman <niekbijman@gmail.com>"
WORKDIR /cli

# Copy the source code into the container
COPY go.mod /cli
COPY go.sum /cli
COPY . /cli

# Build the Golang CLI app
RUN go build -o bankid

ENTRYPOINT ["./bankid"]