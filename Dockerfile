
#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/image-process
COPY . .
RUN apk add --no-cache git
RUN go get -d -v ./...
RUN go install -v ./...

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/image-process /app
ENTRYPOINT ./app
LABEL Name=golang-workspace Version=0.0.1
EXPOSE 8000
