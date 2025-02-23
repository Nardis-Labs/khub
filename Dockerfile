
# Use a multi-stage build to compile both the Go binary and React client app
FROM golang:1.23.0-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the Go source code into the container at /app
COPY . .

COPY client/build /app/client/build/

# Install the necessary dependencies
RUN set -eux \
    & apk update && apk add \
        --no-cache \
        git gcc g++ libc-dev musl-dev ca-certificates bash;

# Set the working directory to /app
WORKDIR /app

# Compile the Go binary with optimizations enabled

COPY tmp/khub .

# Expose the port on which the Go binary will listen
EXPOSE 8080