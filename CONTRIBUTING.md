# Contributing to protoc-gen-utcp

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Development Setup

1. **Fork and clone** the repository
```bash
git clone https://github.com/YOUR_USERNAME/protoc-gen-utcp.git
cd protoc-gen-utcp
```

2. **Install dependencies**
```bash
go mod download
```

3. **Build the plugin**
```bash
make build
```

4. **Run tests**
```bash
make test
```

## Making Changes

### Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Run `go vet` to catch common issues
- Add comments for exported types and functions

### Testing

- Add tests for new functionality
- Ensure all tests pass with `make test`
- Maintain or improve test coverage

### Commits

- Write clear, descriptive commit messages
- Reference related issues in commits
- Keep commits focused and atomic

### Pull Requests

1. Create a feature branch from `main`
2. Make your changes
3. Add tests
4. Update documentation if needed
5. Ensure `make test` passes
6. Submit a pull request

## Project Structure

```
.
├── cmd/protoc-gen-utcp/    # Plugin entrypoint
├── internal/
│   ├── generator/          # UTCP generation logic
│   └── utcp/              # UTCP type definitions
├── examples/              # Example proto files
├── Makefile              # Build automation
└── README.md             # Documentation
```

## Testing Your Changes

### Unit Tests

```bash
make test
```

### Integration Tests

```bash
make examples
# Check generated files in examples/
```

### Manual Testing

```bash
# Install your local version
make install

# Test with your own proto file
protoc --utcp_out=. \
  --utcp_opt=base_url=https://api.example.com \
  your_service.proto
```

## Reporting Issues

When reporting issues, please include:

- Go version (`go version`)
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Sample proto file if applicable

## Questions?

- Open an issue for questions
- Check existing issues first
- Be respectful and constructive

Thank you for contributing!
