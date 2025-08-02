# Contributing to SerdeVal

Thank you for your interest in contributing to SerdeVal! We welcome contributions from the community.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone git@github.com:YOUR_USERNAME/serdeval.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`
4. Make your changes
5. Run tests: `make test`
6. Run linters: `make lint`
7. Commit your changes
8. Push to your fork
9. Create a pull request

## Development Setup

```bash
# Install dependencies
go mod download

# Run tests
make test

# Run linters
make lint

# Build the project
make build

# Run the web interface locally
make run-web
```

## Code Style

- Follow standard Go conventions
- Run `gofmt -w .` before committing
- Ensure all tests pass
- Add tests for new functionality
- Keep commits focused and atomic

## Testing

- Write unit tests for new functionality
- Ensure all existing tests pass
- Include integration tests where appropriate
- Test on multiple platforms if possible

## Pull Request Process

1. Update documentation if needed
2. Add tests for your changes
3. Ensure all tests pass
4. Update the README.md if needed
5. Follow the pull request template

## Reporting Issues

- Use the issue templates provided
- Include as much detail as possible
- Provide sample data that reproduces the issue
- Include your environment details

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on constructive criticism
- Assume good intentions

## Questions?

Feel free to open an issue with the question template if you need help!