# builder stage
FROM golang:alpine as builder
RUN mkdir /build
COPY ./go/* /build/
WORKDIR /build
RUN GOOS=linux go build -o app .
# worker stage
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
COPY --from=builder /build/app .
COPY ./go/hosts.txt /hosts.txt
ENTRYPOINT /app
