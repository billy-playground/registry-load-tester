#!/bin/bash

# Number of instances to run, default is 50
NUM_INSTANCES=${1:-50}

# Array to store PIDs
pids=()

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


# cleanup csv with header
echo "json_file,total_size,download_milliseconds" > results.csv

# Run instances in parallel
echo "Starting $NUM_INSTANCES test instances..."
for ((i=1; i<=NUM_INSTANCES; i++)); do
    ./runner.sh &
    pids+=($!)
done

# Wait for all instances to complete
wait_for_all
EXIT_CODE=$?

if [ $EXIT_CODE -eq 0 ]; then
    echo "All $NUM_INSTANCES instances completed successfully"
else
    echo "One or more instances failed"
fi

cat results.csv

exit $EXIT_CODE