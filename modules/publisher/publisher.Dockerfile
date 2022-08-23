FROM golang:1.19 as builder

WORKDIR /publisher

# Copy and download dependency using go mod
COPY ./modules/publisher/go.mod ./
COPY ./modules/publisher/go.sum ./
RUN go mod download

# Copy the code into the container
COPY ./modules/publisher ./

ENTRYPOINT ["go", "run", "./cmd/front-end_publisher/main.go"]