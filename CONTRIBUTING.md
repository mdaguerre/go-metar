# Contributing to go-metar

Thanks for your interest in contributing to go-metar! This document outlines how to get started.

## Development Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/mdaguerre/go-metar.git
   cd go-metar
   ```

2. **Install Go** (1.21 or later)
   ```bash
   # macOS
   brew install go
   ```

3. **Install dependencies**
   ```bash
   go mod tidy
   ```

4. **Build and run**
   ```bash
   go build -o go-metar
   ./go-metar KJFK
   ```

## Running Tests

```bash
# Run unit tests only (fast, no network)
go test ./... -short

# Run all tests including integration tests
go test ./...
```

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions focused and small

## Project Structure

```
go-metar/
├── main.go           # CLI entry point (Cobra commands)
├── metar/
│   ├── client.go     # API client (Fetch, FetchTAF, etc.)
│   ├── client_test.go
│   └── formatter.go  # Output formatting (Decode, DecodeTAF)
├── docs/
│   └── index.html    # Landing page
└── README.md
```

## Making Changes

1. **Create a branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Write tests for new functionality
   - Ensure all tests pass
   - Update documentation if needed

3. **Commit with clear messages**
   ```bash
   git commit -m "Add feature X"
   ```

4. **Push and create a PR**
   ```bash
   git push origin feature/your-feature-name
   ```

## Pull Request Guidelines

- Keep PRs focused on a single feature or fix
- Include tests for new functionality
- Update README.md if adding new flags or features
- Ensure CI passes before requesting review

## Adding New Features

When adding new data types (like TAF support):

1. Add structs to `metar/client.go`
2. Add fetch functions (`FetchX`, `FetchMultipleX`)
3. Add formatting in `metar/formatter.go`
4. Add CLI flag in `main.go`
5. Add tests in `metar/client_test.go`
6. Update README.md and docs

## Reporting Issues

- Check existing issues first
- Include steps to reproduce
- Include go-metar version (`go-metar --version`)
- Include OS and Go version

## Data Source

This project uses the [Aviation Weather Center API](https://aviationweather.gov/). When adding features, check their API documentation for available endpoints.

## Questions?

Open an issue with the "question" label.
