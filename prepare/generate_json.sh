#!/bin/bash

# Base registry URL
REGISTRY="mcr.azure.cn"

# Output directory for JSON files
OUTPUT_DIR="json_files"

# Input file containing repositories
INPUT_FILE="repositories.json"

# Ensure output directory exists
mkdir -p "$OUTPUT_DIR"

# Check if input file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: Input file $INPUT_FILE not found"
    exit 1
fi

# Check if oras is installed
if ! command -v oras &> /dev/null; then
    echo "Error: oras command not found. Please install ORAS."
    exit 1
fi

# Function to fetch tags for a repository (using oras tags)
get_tags() {
    local repo="$1"
    oras repo tags "$REGISTRY/$repo" 2>/dev/null | head -n 5
}

# Function to fetch manifest content and extract layers
get_manifest_and_layers() {
    local repo="$1"
    local tag="$2"
    local ref="$REGISTRY/$repo:$tag"
    
    # Fetch manifest for linux/amd64
    local manifest_output
    echo "oras manifest fetch --platform linux/amd64 --format json $ref"
    manifest_output=$(oras manifest fetch --platform linux/amd64 --format json "$ref")

    if [ $? -ne 0 ]; then
        echo "Error fetching manifest for $ref" >&2
        return 1
    fi
    manifest=$(echo "$manifest_output" | jq -r '.content')

    # Get manifest digest
    local manifest_digest
    manifest_digest="$(echo "$manifest_output" | jq -r '.digest')"
    if [ -z "$manifest_digest" ]; then
        echo "Failed to extract digest for $ref" >&2
        return 1
    fi

    # Extract size and blobs from manifest
    local total_size=0
    local blobs=()
    
    # Parse config (if present)
    local config_digest
    local config_size
    config_digest=$(echo "$manifest" | jq -r '.config.digest // empty')
    config_size=$(echo "$manifest" | jq -r '.config.size // 0')
    
    if [ -n "$config_digest" ] && [ "$config_size" -gt 0 ]; then
        blobs+=("$REGISTRY/$repo@$config_digest")
        total_size=$((total_size + config_size))
    fi

    # Parse layers
    local layer_digests
    local layer_sizes
    mapfile -t layer_digests < <(echo "$manifest" | jq -r '.layers[].digest')
    mapfile -t layer_sizes < <(echo "$manifest" | jq -r '.layers[].size')
    
    for ((i=0; i<${#layer_digests[@]}; i++)); do
        local digest="${layer_digests[$i]}"
        local size="${layer_sizes[$i]}"
        if [ -n "$digest" ] && [ "$size" -gt 0 ]; then
            blobs+=("$REGISTRY/$repo@$digest")
            total_size=$((total_size + size))
        fi
    done
    # Call generate_json_file directly with the results
    generate_json_file "$repo" "$tag" "$manifest_digest" "$total_size" "${blobs[@]}"
}

# Function to generate JSON file
generate_json_file() {
    local repo="$1"
    local tag="$2"
    local digest="$3"
    local size="$4"
    shift 4
    local blobs=("$@")
    
    # Replace / with _ in filename
    local filename="${repo//\//_}:${tag}.json"
    local filepath="$OUTPUT_DIR/$filename"
    
    # Create JSON content manually
    local blob_array=""
    for blob in "${blobs[@]}"; do
        blob_array="$blob_array\"$blob\","
    done
    blob_array="${blob_array%,}"  # Remove trailing comma
    
    cat << EOF > "$filepath"
{
    "size": $size,
    "manifest": "$REGISTRY/$repo@$digest",
    "blob": [$blob_array]
}
EOF
    
    echo "Generated $filepath"
}

# Main execution
echo "Reading repositories from $INPUT_FILE..."
repos=$(grep '"repositories"' -A 100 "$INPUT_FILE" | grep -o '"[^"]\+/[^"]\+' | tr -d '"')

for repo in $repos; do
    echo "Processing repository: $repo"
    tags=$(get_tags "$repo")
    
    if [ -z "$tags" ]; then
        echo "No tags found or error fetching tags for $repo" >&2
        continue
    fi
    
    for tag in $tags; do
        echo "Fetching manifest for $repo:$tag"
        get_manifest_and_layers "$repo" "$tag"
        # No need to check return code here since generate_json_file is called directly
    done
done

echo "Done!"