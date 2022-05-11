FROM golang:1.18-alpine

RUN apk update && apk add git

WORKDIR /sample-app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o main .
EXPOSE 8080

COPY run.sh ./
CMD [ "sh", "./run.sh" ]
