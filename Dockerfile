FROM alpine:latest

# Install necessary dependencies
RUN apk add --no-cache bash jq curl grep coreutils

# Copy the ORAS binary from the official ORAS image
COPY --from=ghcr.io/oras-project/oras:v1.2.0 /bin/oras /bin/oras

# Set the working directory
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY *.sh /app/
COPY images /app/images/

# Ensure test.sh and runner.sh are executable
RUN chmod +x /app/test.sh
RUN chmod +x /app/runner.sh

# Define the entry point
ENTRYPOINT ["/app/test.sh"]