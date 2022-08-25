FROM golang:1.19 as builder

WORKDIR /subscriber

COPY ./modules/subscriber/go.mod ./
COPY ./modules/subscriber/go.sum ./
RUN go mod download

COPY ./modules/subscriber ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./subscriber ./cmd/subscriber

FROM alpine:3.15
RUN apk update
WORKDIR /

COPY --from=builder /subscriber .

ENTRYPOINT [ "./subscriber" ]