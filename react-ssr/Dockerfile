# Start from a Node.js base image to have npm available
FROM node:14 as build-stage

# Set the current working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY client/package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the client code
COPY client/ ./

# Build the client code if necessary (e.g., if you're using a tool like webpack)
# RUN npm run build

# Start a new stage from golang base image
FROM golang:1.22

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the Go files
COPY . .

# Copy built JavaScript files from the build stage
COPY --from=build-stage /app ./client
# COPY --from=build-stage /app/dist ./client/dist  # Uncomment if you're building the client code

# Enable CGO
ENV CGO_ENABLED=1

# Build the Go app
RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 3000

# Run the executable
CMD ["./main"]