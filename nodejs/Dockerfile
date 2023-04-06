# Stage 1: Install Node.js v16 and Puppeteer
FROM node:16 AS nodejs

# Install necessary dependencies for Puppeteer and Chromium
RUN apt-get update && \
    apt-get install -y \
        libx11-xcb1 \
        libxcomposite1 \
        libxdamage1 \
        libxext6 \
        libxfixes3 \
        libnss3 \
        libnspr4 \
        libgconf-2-4 \
        libasound2 \
        libatk1.0-0 \
        libatk-bridge2.0-0 \
        libgtk-3-0

# Install Puppeteer
WORKDIR /nodejs
COPY package.json .
RUN npm install --unsafe-perm

# Stage 2: Install Golang 1.20 and copy Node.js files
FROM golang:1.20

# Copy Node.js and Puppeteer files from the previous stage
COPY --from=nodejs /nodejs /nodejs
COPY --from=nodejs /usr/local/bin/node /usr/local/bin/

# Set environment variables for the Node.js and Puppeteer paths
ENV NODE_PATH="/nodejs/node_modules"
ENV PATH="/nodejs/node_modules/.bin:/usr/local/bin:$PATH"

# Set the working directory
WORKDIR /app

# Copy your source code into the container
COPY . .
COPY pa11y.js .

# Download and install any required dependencies
RUN go mod download

# Build the application
RUN go build -o app

# Expose the port your application uses
EXPOSE 8080

# Run the application
CMD ["./app"]
