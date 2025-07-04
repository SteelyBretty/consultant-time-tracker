# Consultant Time Tracker

A time tracking and billing system for consultants with REST and GraphQL APIs.

## Features

- Multi-client project management
- Time allocation and tracking
- Billing rate configuration
- Weekly/monthly reporting
- REST and GraphQL APIs

## Development

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Bruno API client

### Getting Started

1. Clone the repository
2. Copy `.env.example` to `.env`
3. Run `make deps` to download dependencies
4. Run `make run` to start the server

### API Documentation

- REST API: http://localhost:8080/api
- GraphQL Playground: http://localhost:8080/graphql

## Testing

Use Bruno collections in the `bruno-collections` directory.