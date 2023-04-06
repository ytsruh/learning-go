# Stage 1: Install Node.js v16
FROM node:16 AS nodejs

# Stage 2: Install Golang 1.20
FROM golang:1.20

# Set the working directory
WORKDIR /app

# Copy your source code into the container
COPY . .

# Download and install any required dependencies
RUN go mod download
RUN npm install

# Build the application
RUN go build -o app

# Expose the port your application uses
EXPOSE 8080

# Run the application
CMD ["./app"]
