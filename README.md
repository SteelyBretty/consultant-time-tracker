# Consultant Time Tracker

A comprehensive time tracking and billing system for consultants with REST and GraphQL APIs.

## Features

- Multi-client project management
- Time allocation and tracking
- Billable rate configuration
- Weekly/monthly reporting
- REST and GraphQL APIs
- SQLite database with automatic migrations
- Docker support for development and production
- Basic authentication system

## Tech Stack

- **Language**: Go 1.24
- **Web Framework**: Gin
- **Database**: SQLite with GORM
- **Authentication**: Basic Auth with bcrypt
- **API Testing**: Bruno
- **Containerization**: Docker & Docker Compose

## Getting Started for Developers

### Prerequisites

- Docker Desktop installed and running
- Git
- Bruno (for API testing) - https://www.usebruno.com/

### Quick Start (Docker)

1. Clone the repository:
   ```bash
   git clone https://github.com/SteelyBretty/consultant-time-tracker.git
   cd consultant-time-tracker
   ```

2. Create environment file:
   ```bash
   cp .env.example .env
   ```

3. Start the development server:
   ```bash
   make docker-dev
   ```

4. Verify it's running:
   ```bash
   curl http://localhost:8080/health
   ```
   
   You should see:
   ```json
   {
     "database": "healthy",
     "service": "consultant-time-tracker",
     "status": "healthy",
     "version": "0.1.0"
   }
   ```

The server will automatically reload when you make changes to the code.

### Stopping the Server

To stop the Docker development server:
```bash
make docker-down
```

### Local Development (Alternative)

If you prefer running locally without Docker:

1. Install Go 1.24+
2. Install Air for hot-reload:
   ```bash
   make install-air
   ```
3. Install dependencies:
   ```bash
   make deps
   ```
4. Run locally:
   ```bash
   make dev
   ```

### Available Make Commands

- `make docker-dev` - Run development server in Docker with hot-reload
- `make docker-down` - Stop Docker containers
- `make docker-logs` - View Docker container logs
- `make dev` - Run development server locally with hot-reload
- `make build` - Build production binary
- `make test` - Run tests
- `make clean` - Clean build artifacts and database
- `make deps` - Download Go dependencies

## API Usage Guide

### Base URL
- Local: `http://localhost:8080`

### Authentication
All API endpoints (except health, info, register, and login) require Basic Authentication.
Include the Authorization header: `Basic base64(username:password)`

### Available Endpoints

#### System
- `GET /health` - Health check with database status
- `GET /` - API information

#### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login user (verify credentials)
- `GET /api/v1/auth/me` - Get current user info (requires auth)

#### Clients
- `POST /api/v1/clients` - Create new client
- `GET /api/v1/clients` - List all clients (paginated)
- `GET /api/v1/clients/:id` - Get client details
- `PUT /api/v1/clients/:id` - Update client
- `DELETE /api/v1/clients/:id` - Delete client (soft delete)

#### Projects (Coming in Milestone 6)
- `POST /api/v1/projects` - Create new project
- `GET /api/v1/projects` - List all projects
- `GET /api/v1/projects/:id` - Get project details
- `PUT /api/v1/projects/:id` - Update project
- `DELETE /api/v1/projects/:id` - Delete project

## Testing with Bruno

### Setup Bruno Collection

1. Download and install Bruno from https://www.usebruno.com/

2. Open the collection:
   - Launch Bruno
   - Click "Open Collection"
   - Navigate to `consultant-time-tracker/bruno-collections/TimeTracker-REST`
   - Click "Open"

3. Select environment:
   - In the top-right dropdown, select "Local"
   - This sets the base URL to `http://localhost:8080`

### Running Tests

1. Ensure the server is running:
   ```bash
   make docker-dev
   ```

2. Follow this test flow:
   - **Register** a new user (or use existing)
   - **Login** to verify credentials
   - **Create Client** to add test data
   - **List Clients** to see all clients
   - **Update/Delete** as needed

3. Each request includes automated tests that verify the response

### API Testing Workflow

#### 1. First Time Setup
```
1. Register User → 2. Login (verify) → 3. Create Clients
```

#### 2. Client Management Flow
```
1. Create Client → 2. List Clients → 3. Get Client → 4. Update Client → 5. Delete Client
```

### Test Data

The Bruno collection includes sample data for:
- **Users**: johndoe/password123, janedoe/password456
- **Clients**: Acme Corporation (ACME), TechCorp Solutions (TECH)

## API Documentation

### Authentication Endpoints

#### Register User
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123",
  "full_name": "John Doe"
}
```

#### Login
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "johndoe",
  "password": "password123"
}
```

### Client Endpoints

#### Create Client
```
POST /api/v1/clients
Authorization: Basic base64(username:password)
Content-Type: application/json

{
  "name": "Acme Corporation",
  "code": "ACME",
  "email": "contact@acme.com",
  "phone": "+1-555-1234",
  "address": "123 Business St, New York, NY"
}
```

#### List Clients
```
GET /api/v1/clients?limit=20&offset=0&search=acme&is_active=true
Authorization: Basic base64(username:password)
```

Query Parameters:
- `limit` - Number of results (max 100, default 20)
- `offset` - Skip N results for pagination
- `search` - Search in name, code, or email
- `is_active` - Filter by active status (true/false)

#### Get Client
```
GET /api/v1/clients/{id}
Authorization: Basic base64(username:password)
```

#### Update Client
```
PUT /api/v1/clients/{id}
Authorization: Basic base64(username:password)
Content-Type: application/json

{
  "name": "Updated Name",
  "is_active": false
}
```

#### Delete Client
```
DELETE /api/v1/clients/{id}
Authorization: Basic base64(username:password)
```

## Database Schema

### users
- `id` (UUID) - Primary key
- `username` (string) - Unique, for login
- `password` (string) - Bcrypt hashed
- `email` (string) - Unique
- `full_name` (string) - Display name
- `is_active` (boolean) - Account status
- `created_at`, `updated_at`, `deleted_at` - Timestamps

### clients
- `id` (UUID) - Primary key
- `name` (string) - Company name
- `code` (string) - Unique per user, uppercase
- `email` (string) - Contact email
- `phone` (string) - Contact phone
- `address` (string) - Full address
- `is_active` (boolean) - Client status
- `user_id` (UUID) - Owner reference
- `created_at`, `updated_at`, `deleted_at` - Timestamps

### projects (Coming Soon)
- `id` (UUID) - Primary key
- `name` (string) - Project name
- `code` (string) - Unique project code
- `client_id` (UUID) - Client reference
- `billable_rate` (float) - Hourly rate
- `status` (enum) - active/on_hold/completed/cancelled

## Error Handling

The API returns consistent error responses:

```json
{
  "error": "Error message here",
  "details": "Additional context (optional)"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `204` - No Content (successful delete)
- `400` - Bad Request (validation error)
- `401` - Unauthorized (missing/invalid auth)
- `404` - Not Found
- `409` - Conflict (duplicate code/username)
- `500` - Internal Server Error

## Project Structure

```
consultant-time-tracker/
├── cmd/
│   └── server/         # Application entrypoint
├── internal/           # Private application code
│   ├── api/           # Route definitions
│   ├── database/      # Database configuration
│   ├── handlers/      # HTTP request handlers
│   │   ├── auth.go   # Authentication endpoints
│   │   └── clients.go # Client CRUD endpoints
│   ├── middleware/    # HTTP middleware
│   │   └── auth.go   # Basic Auth middleware
│   ├── models/        # Database models
│   ├── schemas/       # Request/Response schemas
│   └── services/      # Business logic
│       ├── auth_service.go    # Auth logic
│       └── client_service.go  # Client logic
├── docker/            # Docker configurations
├── bruno-collections/ # API test collections
│   └── TimeTracker-REST/
│       ├── auth/     # Auth test requests
│       ├── clients/  # Client test requests
│       └── system/   # Health checks
├── data/             # SQLite database file
└── README.md         # This file
```

## Development Workflow

### Adding New Features

1. Create feature branch:
   ```bash
   git checkout -b feature/your-feature
   ```

2. Implement feature with TDD approach:
   - Write Bruno tests first
   - Implement the feature
   - Ensure all tests pass

3. Submit PR with:
   - Clear description
   - Test instructions
   - Bruno collection updates

### Code Organization

- **Models**: Database entities in `internal/models/`
- **Services**: Business logic in `internal/services/`
- **Handlers**: HTTP handlers in `internal/handlers/`
- **Schemas**: Request/Response types in `internal/schemas/`
- **Middleware**: Auth, CORS, etc. in `internal/middleware/`

## Troubleshooting

### Cannot Connect to Server
```bash
# Check if server is running
docker ps

# View logs
make docker-logs

# Restart server
make docker-down
make docker-dev
```

### Authentication Issues
- Ensure you're using Basic Auth format
- Username and password are case-sensitive
- Check user is active in database

### Database Issues
```bash
# Reset database
rm -f data/timetracker.db*
make docker-dev
```

### Port Conflicts
```bash
# Check what's using port 8080
lsof -i :8080

# Change port in .env file
APP_PORT=8081
```

## Next Steps

- **Milestone 6**: Project Management
- **Milestone 7**: Time Allocations
- **Milestone 8**: Time Entry Tracking
- **Milestone 9**: GraphQL API
- **Milestone 10**: Reporting & Analytics

## License

[Your License Here]
