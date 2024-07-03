#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Check if the output directory is provided as a command line argument
if [[ -z "$1" ]]; then
    echo "Usage: $0 <output_directory>"
    exit 1
fi

OUTPUT_DIR="$1"

# Create the output directory if it does not exist
if [[ ! -d "$OUTPUT_DIR" ]]; then
    mkdir -p "$OUTPUT_DIR"
fi

# Loop through all directories matching the pattern ./plugin_*
for dir in ./plugin_*; do
    if [[ -d $dir ]]; then
        # Extract the plugin name from the directory name
        plugin_name=$(basename "$dir")
        
        # Construct the main Go file name
        go_file="${dir}/${plugin_name}.go"
        
        # Check if the main Go file exists
        if [[ -f $go_file ]]; then
            # Build the plugin
            output_file="${OUTPUT_DIR}/${plugin_name}.so"
            echo "Building ${output_file} from ${go_file}"
            go build -buildmode=plugin -o "$output_file" "$go_file"
        else
            echo "Main Go file ${go_file} not found, skipping ${plugin_name}"
        fi
    else
        echo "${dir} is not a directory, skipping"
    fi
done

echo "All plugins built successfully."

