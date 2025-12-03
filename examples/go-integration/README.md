# Go Integration Example

This example demonstrates how to programmatically integrate SOPS-encrypted data into your Go applications using the SOPS decrypt package.

## Overview

This example shows how to:
- Load encrypted configuration files using SOPS
- Decrypt data programmatically in your Go application
- Parse decrypted data into Go structs
- Handle multiple configuration formats (JSON, YAML)

## Prerequisites

1. Install SOPS:
   ```bash
   go install github.com/getsops/sops/v3/cmd/sops@latest
   ```

2. Set up PGP keys (for this example):
   ```bash
   # From the sops root directory
   gpg --import pgp/sops_functional_tests_key.asc
   ```

## Files

- `main.go` - Example application that loads and uses encrypted configuration
- `config.enc.json` - Example encrypted JSON configuration
- `config.enc.yaml` - Example encrypted YAML configuration
- `go.mod` - Go module definition

## Running the Example

1. From the `examples/go-integration` directory:

   ```bash
   # Run the example
   go run main.go
   ```

2. To encrypt your own configuration:

   ```bash
   # Create a new config file
   cat > myconfig.json << EOF
   {
     "database": {
       "host": "localhost",
       "password": "secret123"
     }
   }
   EOF

   # Encrypt it with SOPS
   sops -e myconfig.json > myconfig.enc.json
   ```

3. To edit encrypted files:

   ```bash
   sops config.enc.json
   ```

## Using in Your Own Application

Add SOPS as a dependency to your project:

```bash
go get github.com/getsops/sops/v3/decrypt
```

Then use it in your code:

```go
package main

import (
    "encoding/json"
    "log"
    
    "github.com/getsops/sops/v3/decrypt"
)

type Config struct {
    Database struct {
        Host     string `json:"host"`
        Password string `json:"password"`
    } `json:"database"`
}

func main() {
    // Decrypt the file
    cleartext, err := decrypt.File("config.enc.json", "json")
    if err != nil {
        log.Fatalf("Failed to decrypt config: %v", err)
    }
    
    // Parse into struct
    var cfg Config
    if err := json.Unmarshal(cleartext, &cfg); err != nil {
        log.Fatalf("Failed to parse config: %v", err)
    }
    
    // Use the configuration
    log.Printf("Database host: %s", cfg.Database.Host)
}
```

## Integration Patterns

### Pattern 1: Environment-based Configuration

Load different encrypted config files based on environment:

```go
env := os.Getenv("APP_ENV")
if env == "" {
    env = "development"
}
configFile := fmt.Sprintf("config.%s.enc.json", env)
cleartext, err := decrypt.File(configFile, "json")
```

### Pattern 2: Configuration Override

Load base config and override with encrypted secrets:

```go
// Load base config (unencrypted)
baseConfig := loadBaseConfig()

// Load and decrypt secrets
secretsData, _ := decrypt.File("secrets.enc.json", "json")
json.Unmarshal(secretsData, &baseConfig.Secrets)
```

### Pattern 3: Multiple Format Support

```go
import "gopkg.in/yaml.v3"

// YAML config
yamlData, _ := decrypt.File("config.enc.yaml", "yaml")
yaml.Unmarshal(yamlData, &cfg)

// JSON config
jsonData, _ := decrypt.File("config.enc.json", "json")
json.Unmarshal(jsonData, &cfg)
```

## Security Best Practices

1. **Never commit unencrypted secrets** - Always use `.gitignore` to exclude decrypted files
2. **Use appropriate key management** - Store PGP keys or cloud KMS credentials securely
3. **Rotate keys regularly** - Use `sops rotate` to update encryption keys
4. **Limit access** - Only give decryption access to services that need it
5. **Audit access** - Use SOPS audit features to track who accesses secrets

## Additional Resources

- [SOPS Documentation](https://github.com/getsops/sops)
- [SOPS Decrypt Package](https://pkg.go.dev/github.com/getsops/sops/v3/decrypt)
