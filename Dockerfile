FROM golang:bullseye
RUN apt-get update -y && apt-get install libvips-dev -y

WORKDIR /go/src/app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
RUN go build -o app
CMD ["./app"]
