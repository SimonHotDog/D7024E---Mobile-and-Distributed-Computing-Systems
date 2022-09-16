FROM ubuntu:18.04
RUN apt-get update && apt-get install iputils-ping -y

FROM golang:1.19.1-alpine
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./
COPY kademlia/ kademlia/
COPY cli/ cli/
EXPOSE 8080