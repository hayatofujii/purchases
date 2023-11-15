# https://hub.docker.com/_/golang
FROM golang:1.19.1

# Create and change to the app directory.
WORKDIR /app

# Copy application dependency manifests to the container image.
COPY go.mod go.sum ./

# Install dependencies
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Run the web service on container startup.
CMD [ "go", "run", "main.go" ]