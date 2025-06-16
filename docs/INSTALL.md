# Installation Guide

This guide covers different ways to install `aws-sso-config`.

## Quick Install (Recommended)

### Using Go Install

The fastest way to install the latest version:

```bash
go install github.com/blairham/aws-sso-config@latest
```

This will:
- Download and install the latest version from GitHub
- Install the binary to `$GOPATH/bin` or `$GOBIN`
- Include proper version information

### Verify Installation

Check that the installation was successful:

```bash
aws-sso-config --version
```

If you get "command not found", make sure your Go bin directory is in your PATH:

```bash
# Check where Go installs binaries
go env GOPATH

# Add to your shell profile (~/.bashrc, ~/.zshrc, etc.)
export PATH="$PATH:$(go env GOPATH)/bin"

# Or if you have GOBIN set
export PATH="$PATH:$(go env GOBIN)"
```

## Alternative Installation Methods

### 1. Download Pre-built Binaries

Download from the [releases page](https://github.com/blairham/aws-sso-config/releases):

```bash
# Example for Linux amd64
curl -L https://github.com/blairham/aws-sso-config/releases/latest/download/aws-sso-config_linux_amd64.tar.gz | tar xz
sudo mv aws-sso-config /usr/local/bin/
```

### 2. Build from Source

#### Using Make (Recommended)

```bash
git clone https://github.com/blairham/aws-sso-config.git
cd aws-sso-config
make install
```

#### Using Go Build

```bash
git clone https://github.com/blairham/aws-sso-config.git
cd aws-sso-config
go install .
```

### 3. Using Docker

```bash
# Pull the image
docker pull ghcr.io/blairham/aws-sso-config:latest

# Run with your config
docker run --rm -it \
  -v ~/.aws:/home/appuser/.aws:ro \
  -v ~/.awsssoconfig:/home/appuser/.awsssoconfig \
  -v $(pwd):/workspace \
  ghcr.io/blairham/aws-sso-config:latest --help
```

See [DOCKER.md](DOCKER.md) for comprehensive Docker usage guide.

## Development Installation

For developers who want to install from local changes:

```bash
# Clone and install with version info
git clone https://github.com/blairham/aws-sso-config.git
cd aws-sso-config
make go-install-dev

# Or install without version info
go install .
```

## Troubleshooting

### Command Not Found

If you get "command not found" after installation:

1. Check if the binary was installed:
   ```bash
   ls $(go env GOPATH)/bin/aws-sso-config
   ```

2. Check your PATH:
   ```bash
   echo $PATH | grep -q "$(go env GOPATH)/bin" && echo "✓ GOPATH/bin in PATH" || echo "✗ GOPATH/bin not in PATH"
   ```

3. Add Go bin to PATH (add to your shell profile):
   ```bash
   export PATH="$PATH:$(go env GOPATH)/bin"
   ```

### Permission Issues

If you get permission errors:

```bash
# Check Go environment
go env GOPATH GOBIN

# Make sure you have write permissions
ls -la $(go env GOPATH)/bin/
```

### Version Issues

To install a specific version:

```bash
# Install specific version
go install github.com/blairham/aws-sso-config@v1.0.0

# Install from specific branch
go install github.com/blairham/aws-sso-config@main

# Install from specific commit
go install github.com/blairham/aws-sso-config@abc1234
```

## Updating

To update to the latest version:

```bash
go install github.com/blairham/aws-sso-config@latest
```

## Uninstalling

To remove the installation:

```bash
rm $(go env GOPATH)/bin/aws-sso-config
```

## Platform-Specific Notes

### macOS

On macOS, you might need to allow the binary to run:
```bash
# If you get "cannot be opened because the developer cannot be verified"
xattr -d com.apple.quarantine $(which aws-sso-config)
```

### Windows

On Windows, make sure your Go bin directory is in your PATH:
```cmd
# Add to PATH (replace with your actual GOPATH)
setx PATH "%PATH%;%GOPATH%\bin"
```

### Linux

Most Linux distributions work out of the box. For systemd users, you might want to:
```bash
# Reload shell environment
systemctl --user daemon-reload
```

## Next Steps

After installation, see the [README](README.md) for usage instructions or run:

```bash
aws-sso-config --help
```
