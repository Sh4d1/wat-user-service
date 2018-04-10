FROM golang:1.10 as builder

WORKDIR /go/src/github.com/Sh4d1/wat-user-service

COPY . .

#RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:3.6

#RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app

COPY --from=builder /go/src/github.com/Sh4d1/wat-user-service .

CMD ["./wat-user-service"]
