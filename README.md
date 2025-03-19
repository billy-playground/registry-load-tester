# MC-MCR Load Test

This project contains scripts to perform load testing on the Mooncake Microsoft Container Registry (MCR). The main script, `test.sh`, runs multiple instances of `runner.sh` in parallel to simulate concurrent downloads.

## Prerequisites

Before running the scripts, ensure you have the following installed:
- Bash
- `jq` (for parsing JSON)
- `curl` (for making HTTP requests)
- `oras` (OCI Registry As Storage CLI)

If you are using the provided Dockerfile, these dependencies are already included in the container.

## Usage

### Running the Test Script

The `test.sh` script is the main entry point for the load test. It accepts one optional argument to specify the number of parallel instances to run.

```bash
./test.sh [NUM_INSTANCES]
```

- `NUM_INSTANCES`: The number of parallel instances to run. Defaults to `50` if not provided.

### Example

To run the script with 1 instance:

```bash
./test.sh 1
```

To run the script with 10 instances:

```bash
./test.sh 10
```

### Output

The script generates a CSV file named `results.csv` in the current directory. The file contains the following columns:
- `json_file`: The name of the JSON file being processed.
- `total_size`: The total size of the blobs downloaded.
- `download_milliseconds`: The time taken to download the blobs, in milliseconds.

### Using the Docker Container

You can also run the scripts inside a Docker container using the provided `Dockerfile`.

1. Build the Docker image:

   ```bash
   docker build -t mc-mcr-load-test .
   ```

2. Run the container:

   ```bash
   docker run --rm -v $(pwd):/app mc-mcr-load-test [NUM_INSTANCES]
   ```

   Replace `[NUM_INSTANCES]` with the desired number of parallel instances.

### Notes

- Ensure that the `images` directory contains valid JSON files with the required structure before running the scripts.
- The `results.csv` file will be overwritten each time `test.sh` is executed.

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.