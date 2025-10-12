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
│       ├── handler/     # HTTP handlers
│       ├── request/     # Request DTOs
│       ├── response/    # Response DTOs
│       ├── server.go
│       ├── router.go
│       └── middleware.go
├── pkg/                 # Shared packages
├── migrations/          # Database migrations
└── go.mod
```

## Setup & Run

Coming soon...

## API Endpoints

Coming soon...

## Design Decisions

Coming soon...

