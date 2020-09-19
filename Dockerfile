FROM golang:1.14-alpine

WORKDIR /app
COPY . .

RUN go build .

CMD ["./be_nms"]