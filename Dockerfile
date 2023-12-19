# Use an official Golang runtime as a base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the local package files to the container's workspace
COPY . .

# Build the Go server binary
RUN go build -o app .

# Expose port 8000 to the outside world
EXPOSE 8000

# Command to run the executable
CMD ["./app"]
