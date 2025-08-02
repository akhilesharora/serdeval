# Contributing to SerdeVal

First off, thank you for considering contributing to SerdeVal! We want to make contributing to this project as easy and transparent as possible.

## Our Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

## Code of Conduct

This project adheres to a Code of Conduct. By participating, you are expected to uphold this code:
- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on what is best for the community
- Show empathy towards other community members

## Pull Requests

We actively welcome your pull requests:

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. If you've changed APIs, update the documentation.
4. Ensure the test suite passes.
5. Make sure your code follows the existing style.
6. Issue that pull request!

## Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/serdeval.git
cd serdeval

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build -o serdeval

# Run locally
./serdeval web --port 8080
```

## Any contributions you make will be under the MIT Software License

When you submit code changes, your submissions are understood to be under the same [MIT License](LICENSE) that covers the project.

## Report bugs using GitHub's [issues](https://github.com/akhilesharora/serdeval/issues)

We use GitHub issues to track public bugs. Report a bug by [opening a new issue](https://github.com/akhilesharora/serdeval/issues/new).

**Great Bug Reports** tend to have:
- A quick summary and/or background
- Steps to reproduce
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening)

## Privacy First

Remember, SerdeVal is a privacy-focused tool. Any contributions should:
- Not add tracking or analytics
- Not make network requests (except for essential functionality)
- Not store or log user data
- Keep all validation local/client-side where possible

## License

By contributing, you agree that your contributions will be licensed under its MIT License.