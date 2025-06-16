# aws-sso-config

[![Go Report Card](https://goreportcard.com/badge/github.com/blairham/aws-sso-config)](https://goreportcard.com/report/github.com/blairham/aws-sso-config)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A command-line tool for managing AWS configuration files and SSO authentication.

## Features

- 🔧 **Automatic AWS Config Generation**: Generate AWS config files from SSO accounts
- 🚀 **Modern SSO Flow**: Browser-based authentication with automatic polling
- 🔄 **Profile Management**: Automatically detect and configure AWS profiles
- 🛡️ **Secure**: Uses AWS SDK v2 and follows security best practices
- 📦 **Easy Distribution**: Available as binary releases for multiple platforms

## Installation

### Quick Install (Recommended)

```bash
go install github.com/blairham/aws-sso-config@latest
```

**Note**: The binary will be installed to your `$GOPATH/bin` or `$GOBIN` directory. Make sure this directory is in your `$PATH` to use the command globally.

If you get "command not found", add your Go bin directory to PATH:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Other Installation Methods

- **Pre-built binaries**: Download from the [releases page](https://github.com/blairham/aws-sso-config/releases)
- **Build from source**: See [INSTALL.md](docs/INSTALL.md) for detailed instructions

For complete installation instructions, troubleshooting, and platform-specific notes, see [INSTALL.md](docs/INSTALL.md).

### Verify Installation

```bash
aws-sso-config --version
```

## Usage

### Configuration Management

The configuration file (`~/.awsssoconfig`) is automatically created when first needed. You can manage configuration values using git-like commands:

```bash
# Read configuration values
aws-sso-config config get sso_start_url
aws-sso-config config get default_region

# Write configuration values
aws-sso-config config set sso_start_url "https://mycompany.awsapps.com/start"
aws-sso-config config set default_region "us-west-2"
aws-sso-config config set dry_run true

# List all configuration values
aws-sso-config config list
```

### Generate AWS Config

Generate an AWS config file with all accounts you have access to:

```bash
# Generate using default configuration
aws-sso-config generate

# Generate using a custom config file
aws-sso-config generate -config=my-config.toml
```

Show differences before applying changes:

```bash
aws-sso-config generate --diff
```

## Configuration

aws-sso-config supports multiple configuration methods with the following precedence order (highest to lowest):

1. **Command-line flags** (e.g., `-config=my-config.toml`)
2. **Configuration file** (`~/.awsssoconfig` in TOML format)
3. **Environment variables** (with `AWS_CONFIG_` prefix)
4. **Default values**

### Configuration File

aws-sso-config automatically creates and manages a configuration file at `~/.awsssoconfig` in TOML format. The file is created automatically when first needed.

Manage configuration using git-like commands:

```bash
# Get configuration values
aws-sso-config config get sso_start_url
aws-sso-config config get default_region

# Set configuration values
aws-sso-config config set sso_start_url "https://mycompany.awsapps.com/start"
aws-sso-config config set default_region "us-west-2"

# List all configuration values
aws-sso-config config list
```

The configuration file (`~/.awsssoconfig`) contains:

```toml
# AWS Config Tool Configuration
# This file is in TOML format

# SSO Configuration
sso_start_url = "https://your-sso-portal.awsapps.com/start"
sso_region = "us-east-1"
sso_role = "AdministratorAccess"

# AWS Configuration
default_region = "us-east-1"
config_file = "~/.aws/config"

# Behavior Settings
backup_configs = true
dry_run = false

# Environment variables can also be used with AWS_CONFIG_ prefix:
# AWS_CONFIG_SSO_START_URL=https://your-sso-portal.awsapps.com/start
# AWS_CONFIG_SSO_REGION=us-east-1
# AWS_CONFIG_SSO_ROLE=AdministratorAccess
# AWS_CONFIG_DEFAULT_REGION=us-east-1
# AWS_CONFIG_CONFIG_FILE=~/.aws/config
# AWS_CONFIG_BACKUP_CONFIGS=true
# AWS_CONFIG_DRY_RUN=false
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `sso_start_url` | Your AWS SSO start URL | `"https://your-sso-portal.awsapps.com/start"` |
| `sso_region` | AWS region for SSO | `"us-east-1"` |
| `sso_role` | SSO role name | `"AdministratorAccess"` |
| `default_region` | Default AWS region for profiles | `"us-east-1"` |
| `config_file` | Path to AWS config file | `"~/.aws/config"` |
| `backup_configs` | Backup existing config files | `true` |
| `dry_run` | Show changes without applying | `false` |

### Using Custom Configuration Files

You can specify a custom configuration file (must be in TOML format):

```bash
# Generate AWS config using a custom file
aws-sso-config generate -config=my-config.toml

# Show differences before applying
aws-sso-config generate -config=my-config.toml -diff
```

### Environment Variables

You can override any configuration setting using environment variables with the `AWS_CONFIG_` prefix:

- `AWS_CONFIG_SSO_START_URL`: Your AWS SSO start URL
- `AWS_CONFIG_SSO_REGION`: AWS region for SSO (default: us-east-1)
- `AWS_CONFIG_SSO_ROLE`: SSO role name (default: AdministratorAccess)
- `AWS_CONFIG_DEFAULT_REGION`: Default AWS region (default: us-east-1)
- `AWS_CONFIG_CONFIG_FILE`: Path to AWS config file (default: ~/.aws/config)
- `AWS_CONFIG_BACKUP_CONFIGS`: Backup existing configs (default: true)
- `AWS_CONFIG_DRY_RUN`: Show changes without applying (default: false)

Example:

```bash
export AWS_CONFIG_SSO_START_URL="https://mycompany.awsapps.com/start"
export AWS_CONFIG_SSO_REGION="us-west-2"
export AWS_CONFIG_BACKUP_CONFIGS="false"

aws-sso-config generate
```

### Legacy Environment Variables

For compatibility, these environment variables are still supported:

- `AWS_PROFILE`: Override automatic profile detection

## Development

### Prerequisites

- Go 1.23.10+
- Make
- golangci-lint (for linting)
- goreleaser (for releases)
- Docker (for running GitHub Actions locally with `act`)
- [act](https://github.com/nektos/act) (optional - for local CI testing)

### Building

```bash
# Install dependencies
make deps

# Run all checks
make check

# Build for development
make build-dev

# Build for all platforms
make build-all

# Run tests
make test

# Run tests with coverage
make test-coverage
```

### Testing

```bash
# Run all tests
make test

# Run tests with race detection
make test-race

# Generate coverage report (may fail on some tests)
make test-coverage

# Run local CI checks (without Docker)
make ci-local

# Run full GitHub Actions CI locally (requires Docker + act)
make ci

# Run specific CI job locally
make ci-job JOB=security
make ci-job JOB=test
make ci-job JOB=lint
```

### Release Process

This project uses [GoReleaser](https://goreleaser.com) for automated releases:

1. **Create a new tag**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **GitHub Actions automatically**:
   - Builds binaries for multiple platforms (Linux, macOS, Windows)
   - Creates checksums and archives
   - Publishes the release on GitHub
   - Updates the changelog

3. **Manual release** (if needed):
   ```bash
   # Check GoReleaser configuration
   make goreleaser-check

   # Create a snapshot build for testing
   make snapshot

   # Create a full release (requires clean git state and proper tag)
   make release
   ```

### Pre-commit Hooks

Install pre-commit hooks for code quality:

```bash
pre-commit install
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run linting (`make lint`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

Please ensure your code follows the existing style and includes appropriate tests.

## Documentation

Comprehensive documentation is available in the [`docs/`](docs/) directory:

- [Installation Guide](docs/INSTALL.md) - Detailed installation instructions and troubleshooting
- [Docker Guide](docs/DOCKER.md) - Docker usage and deployment
- [Contributing Guidelines](docs/CONTRIBUTING.md) - How to contribute to the project
- [Changelog](docs/CHANGELOG.md) - Version history and changes
- [Security Policy](docs/SECURITY.md) - Security vulnerability reporting
- [Paging Guide](docs/PAGING.md) - Interactive pager controls and features

## Security

If you discover a security vulnerability, please send an email to [your-email]. All security vulnerabilities will be promptly addressed.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- AWS SDK for Go v2
- Mitchell Hashimoto's CLI library
- The Go community

---

**⚠️ Note**: This tool modifies your AWS configuration files. Please ensure you have backups before running.
