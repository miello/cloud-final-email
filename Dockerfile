# Base Image
FROM golang:1.20.3-alpine3.17 as base

# Working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download && go mod verify

# Copy the source code
COPY . .

# Build the application
RUN go build -o server ./main.go

# Create master image
FROM alpine AS master

# Working directory
WORKDIR /app

# Copy execute file
COPY --from=base /app/server ./

# Set ENV to production
ENV GO_ENV production

# Run the application
CMD ["./server"]