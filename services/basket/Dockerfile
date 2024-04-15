# Use the official Golang image to create a build artifact.
FROM golang:1.21.6 AS builder

# Set the Current Working Directory inside the container.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed.
RUN go mod download && go mod tidy

# Copy the source code into the container.
COPY . .

# Build the Go app.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o myapp .

# Start from a scratch (empty) container to keep the final image clean and small.
FROM scratch

# Set the working directory in the scratch image.
WORKDIR /app

# Copy the compiled binary from the build container.
COPY --from=builder /app/myapp /app/myapp

# Copy the configuration file from the build context to the same directory as the binary.
COPY ./basket.prod.toml /app/basket.prod.toml

# Command to run the executable.
CMD ["./myapp"]
