# Start from a base Golang image
FROM golang:1.21-alpine as builder

RUN apk add --no-cache ca-certificates git make musl-dev go

LABEL maintainer="Nicolaas Bijman <niekbijman@gmail.com>"
WORKDIR /cmd/cli

ENV GO111MODULE=off  
ENV GOOS=linux 
ENV GOARCH=amd64
ENV CGO_ENABLED=0

# Copy the source code into the container
COPY ./ ./

# Build the Golang CLI app
RUN go build -o bankid 

# make binary executable
RUN chmod +x ./bankid

# multi-stage docker build in order to hide the secrets used during the build
FROM alpine
RUN apk add --no-cache bash
RUN apk --no-cache add curl

WORKDIR /cli

COPY --from=builder /cmd/cli/bankid /cli/

CMD ["./cli/bankid"]