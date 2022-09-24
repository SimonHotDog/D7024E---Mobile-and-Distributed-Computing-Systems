FROM golang:1.19.1-alpine

WORKDIR /app

COPY go.mod go.sum *.go ./
COPY kademlia/ kademlia/
COPY cli/ cli/

RUN go build -o ./d7024e

EXPOSE 14041

CMD [ "./d7024e" ]
