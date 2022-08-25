FROM golang:1.19 as builder

WORKDIR /publisher

# Copy and download dependency using go mod
COPY ./modules/publisher/go.mod ./
COPY ./modules/publisher/go.sum ./
RUN go mod download

# Copy the code into the container
COPY ./modules/publisher ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./publisher ./cmd/front-end_publisher

FROM alpine:3.15
RUN apk update
WORKDIR /


COPY --from=builder /publisher .

ENTRYPOINT [ "./publisher" ]

#ENTRYPOINT ["go", "run", "./cmd/front-end_publisher/main.go"]

#FROM golang:1.17-alpine as builder
#
#WORKDIR /build
#
#ADD go.mod go.mod
#ADD go.sum go.sum
#ADD cmd cmd
#ADD pkg pkg
#
#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o /tmp/cal-service ./cmd/cal-service
#
#FROM alpine:3.15
#RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
#WORKDIR /
#COPY --from=builder /tmp/cal-service .
#
#ENTRYPOINT [ "./cal-service" ]