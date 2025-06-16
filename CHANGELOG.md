# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- **Configuration file support with Viper**: Full support for YAML, JSON, and TOML configuration files
- **`init` command**: Generate example configuration files in multiple formats
- **`-config` flag**: Specify custom configuration files for the `generate` command
- **Automatic config file discovery**: Searches in current directory, home directory, XDG config directory (`~/.config/aws-config/`), and system config locations
- **Environment variable prefix support**: Use `AWS_CONFIG_` prefix for all configuration options
- Configuration file validation and error handling
- Comprehensive README.md with installation and usage instructions
- MIT License
- Contributing guidelines
- GitHub Actions CI/CD pipeline
- GoReleaser configuration for automated releases
- Configuration management system with environment variable support
- Comprehensive unit tests for core functionality
- Pre-commit hooks for code quality
- Version information embedded in binary builds
- Modern SSO authentication flow with browser integration
- Automatic AWS profile detection based on repository name
- Custom test runner script (`run-tests.sh`) to handle problematic tests

### Changed
- **BREAKING: Renamed commands**: `initconfig` → `init` (now `initialize` in code to avoid Go keyword conflicts)
- **BREAKING: Removed `config` command**: The `config write` and `config read` commands have been removed
- Improved code coverage significantly across all packages
- Updated CLI help to show only the three main commands: `generate`, `init`, and `run`
- Fixed Makefile binary name from `aws-config` to `aws-sso-config` for consistency
- Updated test suite to prevent browser windows from opening during testing
- Removed all TripAdvisor-specific logic and hardcoded values
- Updated SSO URL to be configurable via environment variables
- Improved error handling throughout the codebase
- Enhanced Makefile with comprehensive build targets
- Migrated to GoReleaser for professional build and release process
- Updated .gitignore to follow gitignore.io standards

### Fixed
- Unused import statements
- Linting issues throughout the codebase
- Version handling in main.go
- Build artifact management

### Removed
- Company-specific account name mappings
- Hardcoded SSO URLs and account filters
- Deprecated configuration options

## [0.1.0] - 2025-06-14

### Added
- Initial release
- AWS config file generation from SSO accounts
- Basic CLI structure with subcommands
- SSO authentication flow
- Profile management based on repository detection

---

## Release Process

Releases are automatically created when a new tag is pushed:

```bash
git tag v1.0.0
git push origin v1.0.0
```

This will trigger the GitHub Actions release workflow that builds binaries for multiple platforms and creates a GitHub release.
