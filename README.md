# Transaction Payment Service

A containerized microservice demonstrating safe financial transfers between user accounts using Go, TigerBeetle, PostgreSQL, and Redis.

## Project Structure

```
transaction/
├── cmd/api/              # Application entry point
├── internal/
│   ├── user/            # User domain
│   ├── account/         # Account domain (includes queries)
│   └── http/            # HTTP layer
│       ├── handler/     # HTTP handlers by domain
│       │   ├── user/    # User handler (handler, request, response, mapper)
│       │   ├── account/ # Account handler
│       │   └── health/  # Health handler
│       ├── server.go
│       ├── router.go
│       ├── config.go
│       └── middleware.go
├── pkg/                 # Shared packages
├── migrations/          # Database migrations
├── docker-compose.yml   # Docker services configuration
├── Dockerfile          # Application container
├── Makefile           # Build and test commands
└── go.mod
```

## Setup & Run

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (for local development)

### Quick Start with Docker
```bash
# Clone and navigate to the project
git clone <repository-url>
cd transaction

# Start all services
make docker-up

# Run integration tests
make integration-test

# Stop services
make docker-down
```

### Local Development
```bash
# Copy environment file
cp .env.example .env

# Start dependencies only
docker-compose up -d postgres redis tigerbeetle

# Run migrations
make migrate-up

# Start the application
make run
```

### Environment Variables
Create a `.env` file with the following variables:
```env
SERVER_PORT=8080
API_KEY=your-api-key-here

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=transaction_db
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

TIGERBEETLE_CLUSTER_ID=0
TIGERBEETLE_HOST=localhost
TIGERBEETLE_PORT=3000

LOG_LEVEL=info

MIGRATION_ENABLED=true
MIGRATION_DIRECTION=up
```

## API Endpoints

### Authentication
All endpoints (except `/health` and `POST /users`) require the `X-API-KEY` header.

### User Management
- `POST /api/v1/users` - Create a new user
- `GET /api/v1/users/:id` - Get user by ID

### Account Management
- `POST /api/v1/accounts` - Create a new account
- `GET /api/v1/accounts` - Get user's accounts
- `GET /api/v1/accounts/:id/balance` - Get account balance

### Financial Operations
- `POST /api/v1/accounts/:id/deposit` - Deposit funds to account
- `POST /api/v1/transfers` - Transfer funds between accounts
- `GET /api/v1/accounts/:id/transactions` - Get transaction history

### Health Check
- `GET /health` - Service health status

### Example Usage

```bash
# Create a user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "email": "alice@example.com"}'

# Create an account
curl -X POST http://localhost:8080/api/v1/accounts \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: test-api-key-123" \
  -d '{"user_id": "user-id", "currency": "USD"}'

# Deposit funds
curl -X POST http://localhost:8080/api/v1/accounts/account-id/deposit \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: test-api-key-123" \
  -d '{"amount": 10000, "reference": "initial-deposit"}'

# Transfer funds
curl -X POST http://localhost:8080/api/v1/transfers \
  -H "Content-Type: application/json" \
  -H "X-API-KEY: test-api-key-123" \
  -d '{"from_account_id": "from-id", "to_account_id": "to-id", "amount": 500, "reference": "transfer-1"}'

# Check balance
curl -X GET http://localhost:8080/api/v1/accounts/account-id/balance \
  -H "X-API-KEY: test-api-key-123"
```

## Testing

### Unit Tests
```bash
make test
```

### Integration Tests
```bash
make integration-test
```

### Test Coverage
```bash
go test -cover ./...
```

## Design Decisions

### Architecture
- **Clean Architecture**: Domain-driven design with clear separation of concerns
- **Hexagonal Architecture**: Infrastructure adapters for external dependencies
- **CQRS**: Separate read and write models for better scalability

### Data Consistency
- **Dual-write Pattern**: TigerBeetle ledger + PostgreSQL metadata
- **Idempotency**: Reference-based idempotency for deposits and transfers
- **Distributed Locking**: Redis-based locks to prevent race conditions

### Concurrency Control
- **Optimistic Locking**: Version-based concurrency control in PostgreSQL
- **Distributed Locks**: Redis locks for critical sections
- **Atomic Operations**: Database transactions for consistency

### Caching Strategy
- **Redis Cache**: Account balances cached with 30s TTL
- **Cache-Aside Pattern**: Read-through and write-through cache operations
- **Cache Invalidation**: Immediate invalidation on balance updates

### Error Handling
- **Domain Errors**: Structured error types for business logic
- **HTTP Status Codes**: Proper REST status codes
- **Error Propagation**: Clear error messages without internal details

### Security
- **API Key Authentication**: Simple but effective authentication
- **Input Validation**: Comprehensive request validation
- **SQL Injection Prevention**: Parameterized queries only

### Observability
- **Structured Logging**: JSON logs with correlation IDs
- **Health Checks**: Service health monitoring
- **Metrics**: Request count and transfer success/failure rates

## Limitations & Future Improvements

### Current Limitations
- Single currency support per account
- No currency conversion
- Simple API key authentication
- Limited transaction history pagination
- No audit trail for failed operations

### Future Enhancements
- Multi-currency support
- Currency conversion rates
- OAuth2/JWT authentication
- Event sourcing for audit trails
- GraphQL API
- Real-time notifications
- Circuit breakers for external services
- Distributed tracing

