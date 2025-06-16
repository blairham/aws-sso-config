# Contributing to aws-sso-config

We love your input! We want to make contributing to aws-sso-config as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

## Pull Requests

Pull requests are the best way to propose changes to the codebase. We actively welcome your pull requests:

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Development Setup

### Prerequisites

- Go 1.19 or later
- Make
- Git
- golangci-lint
- goreleaser (for releases)

### Setup

1. Clone your fork:
   ```bash
   git clone https://github.com/yourusername/aws-sso-config.git
   cd aws-sso-config
   ```

2. Install dependencies:
   ```bash
   make deps
   ```

3. Install pre-commit hooks:
   ```bash
   pre-commit install
   ```

### Development Workflow

1. Create a feature branch:
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. Make your changes and write tests

3. Run the full test suite:
   ```bash
   make check
   ```

4. Run tests with coverage:
   ```bash
   make test-coverage
   ```

5. Build the project:
   ```bash
   make build
   ```

6. Commit your changes (this will run pre-commit hooks):
   ```bash
   git add .
   git commit -m "Add my new feature"
   ```

7. Push to your fork and create a pull request

### Code Style

We use the standard Go formatting tools:

- `gofmt` for formatting
- `golangci-lint` for linting
- Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Testing

- Write unit tests for new functionality
- Ensure all tests pass: `make test`
- Check test coverage: `make test-coverage`
- Integration tests should be added for major features

### Commit Messages

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `style:` for formatting changes
- `refactor:` for code refactoring
- `test:` for adding tests
- `chore:` for maintenance tasks

Examples:
```
feat: add support for multiple SSO regions
fix: handle empty profile names gracefully
docs: update installation instructions
```

## Any contributions you make will be under the MIT Software License

In short, when you submit code changes, your submissions are understood to be under the same [MIT License](http://choosealicense.com/licenses/mit/) that covers the project. Feel free to contact the maintainers if that's a concern.

## Report bugs using GitHub's [issue tracker](https://github.com/blairham/aws-sso-config/issues)

We use GitHub issues to track public bugs. Report a bug by [opening a new issue](https://github.com/blairham/aws-sso-config/issues/new).

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

## Feature Requests

We welcome feature requests! Please open an issue with:

- A clear description of the feature
- Why you need it
- How it should work
- Any relevant examples or mockups

## Questions

If you have questions about the project, please:

1. Check the [README](README.md) first
2. Search existing [issues](https://github.com/blairham/aws-sso-config/issues)
3. Open a new issue with the "question" label

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
