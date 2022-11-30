FROM golang:1.19-alpine

RUN apk add --no-cache git

COPY ${PWD} /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
RUN go build api/main.go
EXPOSE 9000

ENTRYPOINT [ "/app/main" ]
