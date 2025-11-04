# Guidelines for AI Coding Agents

## Project Overview

LuckySix is a comprehensive Go application designed to systematically generate, process, and monitor a vast number of blockchain wallets derived from BIP39 mnemonic phrases. The project generates every possible two-word combination, builds upon them to create 12-word seeds, derives wallets for multiple blockchains (EVM and Cosmos), checks their balances, and sends Telegram notifications when a wallet with a non-zero balance is found.

Key features include:
- Combinatorial generation of BIP39 word combinations (LuckyTwo, LuckyFive, LuckySix).
- Multi-chain wallet derivation for Ethereum (EVM) and Cosmos chains.
- Automated balance checking via public RPC endpoints.
- Telegram notifications for wallets with positive balances.
- Service-oriented design with Clean Architecture principles.

The application uses PostgreSQL for data storage and integrates with external APIs like Telegram Bot API and blockchain RPCs.

## Build and Test Commands

### Prerequisites
- Go (version 1.21 or higher)
- PostgreSQL
- Docker and Docker Compose (recommended for local setup)

### Build Commands
- **Build the application**: `go build ./cmd/app`
- **Run the application locally**: `go run ./cmd/app`
- **Docker setup**: `docker-compose up -d` (use `--force-recreate` to rebuild containers)

### Test Commands
- **Run all tests**: `go test ./...`
- **Run tests with coverage**: `go test -cover ./...`
- **Run specific package tests**: `go test ./internal/service`
- **Run integration tests**: `go test ./integration-test`

### Linting and Formatting
- **Format code**: `go fmt ./...`
- **Import organization**: `goimports -w .`
- **Linting**: `golangci-lint run`

### Dependency Management
- **Install/update dependencies**: `go get -u ./... && go mod tidy`
- **Tidy modules**: `go mod tidy`

# Guidelines for AI Coding Agents

## Restricted files
Files in the list contain sensitive data, they MUST NOT be read
- .env

## Code style guidelines

### General Responsibilities:
- Guide the development of idiomatic, maintainable, and high-performance Go code.
- Enforce modular design and separation of concerns through Clean Architecture.
- Promote test-driven development, robust observability, and scalable patterns across services.

### Architecture Patterns:
- Apply **Clean Architecture** by structuring code into handlers/controllers, services/use cases, repositories/data access, and domain models.
- Use **domain-driven design** principles where applicable.
- Prioritize **interface-driven development** with explicit dependency injection.
- Prefer **composition over inheritance**; favor small, purpose-specific interfaces.
- Ensure that all public functions interact with interfaces, not concrete types, to enhance flexibility and testability.

### Project Structure Guidelines:
- Use a consistent project layout:
  - cmd/: application entrypoints
  - internal/: core application logic (not exposed externally)
  - pkg/: shared utilities and packages
  - configs/: configuration schemas and loading
  - test/: test utilities, mocks, and integration tests
- Group code by feature when it improves clarity and cohesion.
- Keep logic decoupled from framework-specific code.

### Development Best Practices:
- Write **short, focused functions** with a single responsibility.
- Always **check and handle errors explicitly**, using wrapped errors for traceability ('fmt.Errorf("context: %w", err)').
- Avoid **global state**; use constructor functions to inject dependencies.
- Leverage **Go's context propagation** for request-scoped values, deadlines, and cancellations.
- Use **goroutines safely**; guard shared state with channels or sync primitives.
- **Defer closing resources** and handle them carefully to avoid leaks.

### Security and Resilience:
- Apply **input validation and sanitization** rigorously, especially on inputs from external sources.
- Isolate sensitive operations with clear **permission boundaries**.
- Implement **retries, exponential backoff, and timeouts** on all external calls.
- Use **circuit breakers and rate limiting** for service protection.
- Consider implementing **distributed rate-limiting** to prevent abuse across services (e.g., using Redis).

### Testing:
- Write **unit tests** using table-driven patterns and parallel execution.
- **Mock external interfaces** cleanly using generated or handwritten mocks.
- Separate **fast unit tests** from slower integration and E2E tests.
- Ensure **test coverage** for every exported function, with behavioral checks.
- Use tools like 'go test -cover' to ensure adequate test coverage.

### Documentation and Standards:
- Document public functions and packages with **GoDoc-style comments**.
- Provide concise **READMEs** for services and libraries.
- Maintain a 'CONTRIBUTING.md' and 'ARCHITECTURE.md' to guide team practices.
- Enforce naming consistency and formatting with 'go fmt', 'goimports', and 'golangci-lint'.

### Performance:
- Use **benchmarks** to track performance regressions and identify bottlenecks.
- Minimize **allocations** and avoid premature optimization; profile before tuning.
- Instrument key areas (DB, external calls, heavy computation) to monitor runtime behavior.

### Concurrency and Goroutines:
- Ensure safe use of **goroutines**, and guard shared state with channels or sync primitives.
- Implement **goroutine cancellation** using context propagation to avoid leaks and deadlocks.

### Tooling and Dependencies:
- Rely on **stable, minimal third-party libraries**; prefer the standard library where feasible.
- Use **Go modules** for dependency management and reproducibility.
- Version-lock dependencies for deterministic builds.
- Integrate **linting, testing, and security checks** in CI pipelines.

### Key Conventions:
1. Prioritize **readability, simplicity, and maintainability**.
2. Design for **change**: isolate business logic and minimize framework lock-in.
3. Emphasize clear **boundaries** and **dependency inversion**.
4. Ensure all behavior is **observable, testable, and documented**.
5. **Automate workflows** for testing, building, and deployment.