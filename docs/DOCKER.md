# Docker Usage Guide

This guide covers how to use `aws-sso-config` with Docker, including building and running the ultra-minimal container.

## Quick Start

### Using Pre-built Image

```bash
# Pull the latest image
docker pull ghcr.io/blairham/aws-sso-config:latest

# Run with your AWS credentials and config
docker run --rm -it \
  -v ~/.aws:/home/appuser/.aws:ro \
  -v ~/.awsssoconfig:/home/appuser/.awsssoconfig \
  -v $(pwd):/workspace \
  ghcr.io/blairham/aws-sso-config:latest --help
```

### Building Locally

The Docker images are built using **GoReleaser** for optimal size and consistency:

```bash
# Build the image using GoReleaser (recommended)
make docker-build

# Test the image
make docker-test

# Run the container
make docker-run ARGS="--version"

# Alternative: Build only Docker images (skip other artifacts)
make docker-only
```

**Why GoReleaser?**
- **Smaller images**: ~5MB vs ~18MB with traditional Docker builds
- **Consistent builds**: Same process for releases and local development
- **Optimized binaries**: Better compression and size optimization
- **Metadata**: Proper labels and versioning information

## Image Details

The Docker image is built using **GoReleaser** with a **scratch base** for minimal size:

- **Base**: `scratch` (empty base image)
- **Size**: ~5 MB (just the binary + certificates)  
- **Build**: GoReleaser-optimized binaries
- **Security**: Runs as non-root user (`appuser`)
- **Contents**: Only the `aws-sso-config` binary and essential certificates
- **Multi-arch**: Supports AMD64 and ARM64 (when published)

## Volume Mounts

### Essential Directories

| Host Path | Container Path | Purpose | Mode |
|-----------|----------------|---------|------|
| `~/.aws` | `/home/appuser/.aws` | AWS credentials and config | `ro` (read-only) |
| `~/.awsssoconfig` | `/home/appuser/.awsssoconfig` | Tool configuration | `rw` (read-write) |
| `$(pwd)` | `/workspace` | Current working directory | `rw` (read-write) |

### Mount Examples

```bash
# Read-only AWS config (recommended for security)
-v ~/.aws:/home/appuser/.aws:ro

# Read-write tool config (allows saving settings)
-v ~/.awsssoconfig:/home/appuser/.awsssoconfig

# Mount current directory for output files
-v $(pwd):/workspace

# Mount specific output directory
-v /path/to/output:/workspace

# Mount temporary directory
-v /tmp/aws-output:/workspace
```

## Usage Examples

### 1. Basic Commands

```bash
# Show version
docker run --rm ghcr.io/blairham/aws-sso-config:latest --version

# Show help
docker run --rm ghcr.io/blairham/aws-sso-config:latest --help

# List config values
docker run --rm \
  -v ~/.awsssoconfig:/home/appuser/.awsssoconfig \
  ghcr.io/blairham/aws-sso-config:latest config list
```

### 2. Configuration Management

```bash
# Get configuration value
docker run --rm \
  -v ~/.awsssoconfig:/home/appuser/.awsssoconfig \
  ghcr.io/blairham/aws-sso-config:latest config get sso.start_url

# Set configuration value
docker run --rm \
  -v ~/.awsssoconfig:/home/appuser/.awsssoconfig \
  ghcr.io/blairham/aws-sso-config:latest config set sso.start_url https://mycompany.awsapps.com/start

# Edit configuration (requires interactive terminal)
docker run --rm -it \
  -v ~/.awsssoconfig:/home/appuser/.awsssoconfig \
  -e EDITOR=vi \
  ghcr.io/blairham/aws-sso-config:latest config edit
```

### 3. Generate AWS Config

```bash
# Generate AWS config to current directory
docker run --rm \
  -v ~/.aws:/home/appuser/.aws:ro \
  -v ~/.awsssoconfig:/home/appuser/.awsssoconfig \
  -v $(pwd):/workspace \
  ghcr.io/blairham/aws-sso-config:latest generate

# Generate to specific file
docker run --rm \
  -v ~/.aws:/home/appuser/.aws:ro \
  -v ~/.awsssoconfig:/home/appuser/.awsssoconfig \
  -v /tmp:/workspace \
  ghcr.io/blairham/aws-sso-config:latest generate --output /workspace/aws-config
```

### 4. Working with Browser Authentication

Since the container can't open a browser, you'll need to handle SSO authentication differently:

```bash
# The tool will display the SSO URL - copy and paste into your browser
docker run --rm -it \
  -v ~/.aws:/home/appuser/.aws \
  -v ~/.awsssoconfig:/home/appuser/.awsssoconfig \
  ghcr.io/blairham/aws-sso-config:latest generate
```

## Makefile Shortcuts

The project includes convenient Makefile targets:

```bash
# Build Docker image
make docker-build

# Run with mounted directories
make docker-run ARGS="config list"

# Run with config mounted
make docker-run-config ARGS="generate"

# Test the image
make docker-test

# Check image size
make docker-size

# Debug (alpine shell for troubleshooting)
make docker-shell
```

## Building Custom Images

### Local Build

```bash
# Build with current version
make docker-build

# Build with custom tag
docker build -t my-aws-sso-config:latest .

# Build with version info
docker build \
  --build-arg VERSION=v1.0.0 \
  --build-arg COMMIT=$(git rev-parse HEAD) \
  --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -t my-aws-sso-config:v1.0.0 \
  .
```

### Multi-platform Build

```bash
# Build for multiple architectures
docker buildx build \
  --platform linux/amd64,linux/arm64 \
  --build-arg VERSION=v1.0.0 \
  --build-arg COMMIT=$(git rev-parse HEAD) \
  --build-arg BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) \
  -t my-aws-sso-config:v1.0.0 \
  .
```

## Docker Compose

Create a `docker-compose.yml` for easier management:

```yaml
version: '3.8'

services:
  aws-sso-config:
    image: ghcr.io/blairham/aws-sso-config:latest
    # Or build locally:
    # build: .
    volumes:
      - ~/.aws:/home/appuser/.aws:ro
      - ~/.awsssoconfig:/home/appuser/.awsssoconfig
      - ./output:/workspace
    working_dir: /workspace
    # Override entrypoint for interactive use
    entrypoint: ""
    command: ["sleep", "infinity"]
    stdin_open: true
    tty: true

  # One-shot service for commands
  aws-sso-config-cmd:
    image: ghcr.io/blairham/aws-sso-config:latest
    volumes:
      - ~/.aws:/home/appuser/.aws:ro
      - ~/.awsssoconfig:/home/appuser/.awsssoconfig
      - ./output:/workspace
    working_dir: /workspace
    # Use: docker-compose run --rm aws-sso-config-cmd config list
```

Usage with Docker Compose:

```bash
# Run commands
docker-compose run --rm aws-sso-config-cmd config list
docker-compose run --rm aws-sso-config-cmd generate

# Interactive shell
docker-compose run --rm aws-sso-config /bin/sh

# Keep container running for multiple commands
docker-compose up -d aws-sso-config
docker-compose exec aws-sso-config config list
docker-compose exec aws-sso-config generate
docker-compose down
```

## Security Considerations

### 1. Read-only AWS Directory

Mount AWS credentials as read-only to prevent accidental modification:

```bash
-v ~/.aws:/home/appuser/.aws:ro
```

### 2. Non-root User

The container runs as `appuser` (non-root) for security:

```dockerfile
USER appuser
```

### 3. Minimal Attack Surface

- Uses `scratch` base image (no OS, no shell, no utilities)
- Only contains the single binary and essential certificates
- No package manager or additional software

### 4. Temporary Directories

Use temporary directories for output to avoid polluting host filesystem:

```bash
# Use temporary directory
-v /tmp/aws-output:/workspace

# Or create temporary directory
mkdir -p /tmp/aws-sso-output
docker run --rm \
  -v ~/.aws:/home/appuser/.aws:ro \
  -v ~/.awsssoconfig:/home/appuser/.awsssoconfig \
  -v /tmp/aws-sso-output:/workspace \
  ghcr.io/blairham/aws-sso-config:latest generate
```

## Troubleshooting

### 1. Permission Issues

If you get permission errors:

```bash
# Check file ownership
ls -la ~/.awsssoconfig

# Fix ownership (if needed)
sudo chown -R $(id -u):$(id -g) ~/.awsssoconfig

# Or use different mount point
-v ~/.awsssoconfig:/tmp/.awsssoconfig
```

### 2. Browser Authentication

Since containers can't open browsers:

1. The tool will display an SSO URL
2. Copy the URL manually
3. Open in your host browser
4. The container will poll for completion

### 3. Debug Container

Use the debug shell for troubleshooting:

```bash
make docker-shell

# Or manually
docker run --rm -it \
  -v ~/.aws:/root/.aws:ro \
  -v $(pwd):/workspace \
  --entrypoint="" \
  alpine:latest sh
```

### 4. Check Image Size

Verify the image is minimal:

```bash
make docker-size
# Should show ~8-12 MB for the scratch-based image

docker images ghcr.io/blairham/aws-sso-config
```

### 5. Verbose Output

Enable verbose mode for debugging:

```bash
docker run --rm -it \
  -v ~/.aws:/home/appuser/.aws:ro \
  -v ~/.awsssoconfig:/home/appuser/.awsssoconfig \
  -v $(pwd):/workspace \
  ghcr.io/blairham/aws-sso-config:latest --verbose generate
```

## Performance Notes

- **Image Size**: ~8-12 MB (vs ~100+ MB for typical Go images)
- **Startup Time**: Near-instantaneous (no OS to boot)
- **Memory Usage**: Minimal (just the Go binary)
- **Security**: Minimal attack surface (scratch base)

## Alternative: Alpine-based Image

If you need shell access or debugging tools, you can create an Alpine-based variant:

```dockerfile
# Alternative Dockerfile.alpine
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY aws-sso-config /usr/local/bin/
ENTRYPOINT ["aws-sso-config"]
```

Build with: `docker build -f Dockerfile.alpine -t aws-sso-config:alpine .`

This provides a shell for debugging while still being small (~15-20 MB).
