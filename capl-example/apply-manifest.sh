#!/bin/bash

# Check if required environment variables are set
required_vars=("LINODE_TOKEN" "LINODE_REGION")
missing_vars=()

for var in "${required_vars[@]}"; do
    if [ -z "${!var}" ]; then
        missing_vars+=("$var")
    fi
done

if [ ${#missing_vars[@]} -ne 0 ]; then
    echo "Error: The following required environment variables are not set:"
    printf '%s\n' "${missing_vars[@]}"
    exit 1
fi

# Apply the manifest with environment variables substituted
envsubst < cluster-manifest.yaml | kubectl apply -f -
