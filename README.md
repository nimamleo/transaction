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
└── go.mod
```

## Setup & Run

Coming soon...

## API Endpoints

Coming soon...

## Design Decisions

Coming soon...

