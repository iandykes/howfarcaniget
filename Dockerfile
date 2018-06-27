FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go get -d
RUN go build -o main .
FROM alpine
COPY --from=builder /build/main /app/
WORKDIR /app
CMD ["./main"]
EXPOSE 80