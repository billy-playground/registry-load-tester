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
all_success=0

# Record start time in seconds since epoch
start_time=$(date +%s)

# Download manifest in background if present
if [ -n "$manifest_ref" ] && [ "$manifest_ref" != "null" ]; then
    oras manifest fetch --output /dev/null "$manifest_ref"  --no-tty & # 2>/dev/null
    pids+=($!)
fi

# Download blobs in background
for blob_url in "${blob_refs[@]}"; do
    oras blob fetch --output /dev/null "$blob_url"  --no-tty & # 2>/dev/null
    pids+=($!)
done

# Wait for all downloads to complete and check their status
for pid in "${pids[@]}"; do
    wait "$pid"
    if [ $? -ne 0 ]; then
        all_success=1  # Set to 1 if any download fails
    fi
done       

# Record end time and calculate elapsed time
end_time=$(date +%s)
download_time=$((end_time - start_time))

if [ $all_success -eq 0 ]; then
    echo "Test completed successfully for $json_file with total layer size: $total_size bytes in $download_time seconds" >&2
    exit 0
else
    echo "Test failed for $json_file due to one or more download failures" >&2
    exit 1
fi