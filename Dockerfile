FROM golang:1.18-alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Add required packages
RUN apk add --update git curl bash

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

# Copy app files
COPY main.go handlers ./

# Install Reflex for development
RUN go install github.com/cespare/reflex@latest

# Expose port
EXPOSE 9090

# Start app
CMD reflex -r '\.go' -s -- sh -c "go run main.go"
