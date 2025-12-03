#!/bin/bash
# Test script for the Go integration example

set -e

echo "Testing SOPS Go Integration Example"
echo "===================================="
echo

# Navigate to the example directory
cd "$(dirname "$0")"

# Import test keys if not already available
if ! gpg --list-secret-keys | grep -q "FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4"; then
    echo "Importing SOPS test keys..."
    if [ -f "../../pgp/sops_functional_tests_key.asc" ]; then
        gpg --import ../../pgp/sops_functional_tests_key.asc 2>/dev/null || true
        echo "✓ Test keys imported"
    else
        echo "✗ Error: Could not find test keys"
        exit 1
    fi
fi

# Check if encrypted config files exist, create them if not
if [ ! -f "config.enc.json" ] || [ ! -f "config.enc.yaml" ]; then
    echo "Creating sample encrypted configs..."
    ./create-sample-configs.sh
fi

# Run the example
echo
echo "Running the integration example..."
if go run main.go; then
    echo
    echo "✓ All tests passed!"
    exit 0
else
    echo
    echo "✗ Tests failed!"
    exit 1
fi
