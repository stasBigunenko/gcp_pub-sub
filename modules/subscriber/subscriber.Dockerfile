FROM golang:1.19 as builder

WORKDIR /subscriber

# Copy and download dependency using go mod
COPY ./modules/subscriber/go.mod ./
COPY ./modules/subscriber/go.sum ./
RUN go mod download

# Copy the code into the container
COPY ./modules/subscriber ./

ENTRYPOINT ["go", "run", "./cmd/subscriber/main.go"]