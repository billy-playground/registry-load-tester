# Use the ORAS image as a parent image
FROM ghcr.io/oras-project/oras:v1.2.0

# Install bash
RUN apk add --no-cache bash jq

# Set the working directory
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Ensure test.sh and runner.sh are executable
RUN chmod +x /app/test.sh
RUN chmod +x /app/runner.sh

# Set the environment variable for the ORAS binary
ENV PATH="/bin/oras:${PATH}"

# Define the entry point
ENTRYPOINT ["/app/test.sh"]