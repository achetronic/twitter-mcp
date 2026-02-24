# Contributing to Twitter MCP

Thanks for your interest in contributing! Here's how you can help.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/twitter-mcp.git`
3. Create a branch: `git checkout -b feature/my-feature`
4. Make your changes
5. Run tests: `go test -v ./...`
6. Commit: `git commit -m "Add my feature"`
7. Push: `git push origin feature/my-feature`
8. Open a Pull Request

## Development Setup

```bash
# Install dependencies
go mod tidy

# Build
make build

# Run (requires config.yaml with Twitter credentials)
make run
```

## Code Style

- Run `go fmt` before committing
- Follow standard Go conventions
- Add comments for exported functions
- Keep functions focused and small

## Adding New Tools

See [AGENTS.md](AGENTS.md) for detailed instructions on adding new MCP tools.

## Testing

- Write tests for new functionality
- Ensure existing tests pass
- Test with real Twitter API when possible (use a test account)

## Pull Request Guidelines

- Keep PRs focused on a single change
- Update documentation if needed
- Add tests for new features
- Ensure CI passes

## Reporting Issues

When reporting bugs, please include:
- Go version (`go version`)
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.
