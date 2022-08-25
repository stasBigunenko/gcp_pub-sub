FROM golang:1.19 as builder

WORKDIR /publisher

COPY ./modules/publisher/go.mod ./
COPY ./modules/publisher/go.sum ./
RUN go mod download

COPY ./modules/publisher ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./publisher ./cmd/front-end_publisher

FROM alpine:3.15
RUN apk update
WORKDIR /

COPY --from=builder /publisher .

ENTRYPOINT [ "./publisher" ]