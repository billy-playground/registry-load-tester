#!/bin/bash

# Directory containing JSON files
JSON_DIR="images"

# Function to pick a random JSON file
pick_random_json() {
    local files=("$JSON_DIR"/*.json)
    if [ ${#files[@]} -eq 0 ]; then
        echo "Error: No JSON files found in $JSON_DIR" >&2
        exit 1
    fi
    echo "${files[RANDOM % ${#files[@]}]}"
}

# Main execution
json_file=$(pick_random_json)
echo "Processing $json_file" >&2

# Parse JSON using jq
total_size=$(jq -r '.size' "$json_file")
manifest_ref=$(jq -r '.manifest' "$json_file")
blob_refs=($(jq -r '.blob[]' "$json_file"))

# Array to store PIDs of background jobs
pids=()
# Variable to track overall success
all_success=1

# Record start time in milliseconds since epoch
start_time=$(date +%s%3N)

# Download manifest in background if present
if [ -n "$manifest_ref" ] && [ "$manifest_ref" != "null" ]; then
    echo "Downloading manifest from $manifest_ref" >&2
    oras manifest fetch --output /dev/null "$manifest_ref" --no-tty & # 2>/dev/null
    pids+=($!)
fi

# Download blobs in background
for blob_url in "${blob_refs[@]}"; do
    echo "Downloading blob from $blob_url" >&2
    oras blob fetch --output /dev/null "$blob_url" --no-tty & # 2>/dev/null
    pids+=($!)
done

# Wait for all downloads to complete and check their status
for pid in "${pids[@]}"; do
    wait "$pid"
    if [ $? -ne 0 ]; then
        all_success=0 # Set to 0 if any download fails
    fi
done

# Record end time and calculate elapsed time in milliseconds
end_time=$(date +%s%3N)
download_milliseconds=$((end_time - start_time))

if [ $all_success -eq 1 ]; then
    # output to csv line, tee to stdout
    echo "$json_file,$total_size,$download_milliseconds" | tee -a results.csv
    echo "Test completed successfully for $json_file: $total_size bytes in $download_milliseconds ms" >&2
    exit 0
else
    echo "Test failed for $json_file due to one or more download failures" >&2
    exit 1
fi
