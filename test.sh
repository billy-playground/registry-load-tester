#!/bin/bash

# Number of instances to run, default is 50
NUM_INSTANCES=${1:-50}

# Registry
registry="mcr.azure.cn"

# Array to store PIDs
pids=()

# Parses the challenge from the authentication header
function parse_challenge() {
    sed -n "s/.*$2=\"\([^\"]*\).*/\1/p" <<< $1
}

# Function to wait for all processes and check their exit status
wait_for_all() {
    local all_success=0
    for pid in "${pids[@]}"; do
        wait "$pid"
        status=$?
        if [ $status -ne 0 ]; then
            all_success=1
        fi
    done
    return $all_success
}

# Get anonymous token and login
auth=$(curl -LIs "https://$registry/v2/" | grep -i "Www-Authenticate:")
realm=$(parse_challenge "$auth" realm)
service=$(parse_challenge "$auth" service)
curl -s -X GET "$realm?service=$service&scope=repository:*:pull" | jq -r '.access_token' | oras login $registry --identity-token-stdin > /dev/null

echo "json_file,total_size,download_milliseconds"
# Run instances in parallel
for ((i=1; i<=NUM_INSTANCES; i++)); do
    ./runner.sh &
    pids+=($!)
done

# Wait for all instances to complete
wait_for_all
EXIT_CODE=$?

# if [ $EXIT_CODE -eq 0 ]; then
#     echo "All $NUM_INSTANCES instances completed successfully"
# else
#     echo "One or more instances failed"
# fi

exit $EXIT_CODE