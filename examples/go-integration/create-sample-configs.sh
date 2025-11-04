#!/bin/bash
# This script creates sample encrypted configuration files for the Go integration example

set -e

# SOPS test key fingerprint
PGP_KEY_FINGERPRINT="FBC7B9E2A4F9289AC0C1D4843D16CEE4A27381B4"

echo "Creating sample encrypted configuration files..."
echo

# Check if sops is available
if ! command -v sops &> /dev/null; then
    echo "Error: sops is not installed or not in PATH"
    echo "Please install sops first: go install github.com/getsops/sops/v3/cmd/sops@latest"
    exit 1
fi

# Check if GPG keys are imported
if ! gpg --list-secret-keys | grep -q "$PGP_KEY_FINGERPRINT"; then
    echo "Warning: SOPS test keys not found in GPG keyring"
    echo "Importing test keys from repository..."
    if [ -f "../../pgp/sops_functional_tests_key.asc" ]; then
        gpg --import ../../pgp/sops_functional_tests_key.asc 2>/dev/null || true
        echo "Test keys imported successfully"
    else
        echo "Error: Could not find test keys at ../../pgp/sops_functional_tests_key.asc"
        echo "Please run this script from the examples/go-integration directory"
        exit 1
    fi
fi

# Create sample JSON config
echo "Creating sample JSON configuration..."
cat > config.json << 'EOF'
{
  "application": {
    "name": "MyApp",
    "environment": "production",
    "port": 8080
  },
  "database": {
    "host": "db.example.com",
    "port": 5432,
    "username": "app_user",
    "password": "super_secret_db_password_123",
    "database": "myapp_db"
  },
  "api_keys": {
    "stripe_key": "sk_live_1234567890abcdef",
    "sendgrid_key": "SG.1234567890abcdef",
    "aws_access_key": "AKIAIOSFODNN7EXAMPLE"
  }
}
EOF

# Encrypt JSON config
echo "Encrypting JSON configuration with SOPS..."
sops -e --pgp "$PGP_KEY_FINGERPRINT" config.json > config.enc.json
rm config.json
echo "  ✓ Created config.enc.json"

# Create sample YAML config
echo "Creating sample YAML configuration..."
cat > config.yaml << 'EOF'
application:
  name: MyApp
  environment: production
  port: 8080

database:
  host: db.example.com
  port: 5432
  username: app_user
  password: super_secret_db_password_123
  database: myapp_db

api_keys:
  stripe_key: sk_live_1234567890abcdef
  sendgrid_key: SG.1234567890abcdef
  aws_access_key: AKIAIOSFODNN7EXAMPLE
EOF

# Encrypt YAML config
echo "Encrypting YAML configuration with SOPS..."
sops -e --pgp "$PGP_KEY_FINGERPRINT" config.yaml > config.enc.yaml
rm config.yaml
echo "  ✓ Created config.enc.yaml"

echo
echo "Sample encrypted configuration files created successfully!"
echo
echo "You can now run the example with:"
echo "  go run main.go"
echo
echo "To edit the encrypted files, use:"
echo "  sops config.enc.json"
echo "  sops config.enc.yaml"
