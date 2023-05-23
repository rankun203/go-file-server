# Use the official Golang image to create a build artifact.
# This is based on Debian and sets the GOPATH to /go.
FROM golang:1.20-alpine as builder

# Create and change to the app directory.
WORKDIR /app

# # Retrieve application dependencies using go modules.
# # Allows container builds to reuse downloaded dependencies.
# COPY go.* ./
# RUN go mod download

# Copy local code to the container image.
COPY ./main.go ./

# Build the binary.
RUN go build -ldflags="-w -s" -o server ./main.go

# Use the official lightweight Alpine image for a lean production container.
# https://hub.docker.com/_/alpine
FROM alpine:3.18

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /server

# Expose the service on port 3000
EXPOSE 3000

# Run the web service on container startup.
ENTRYPOINT ["/server"]
