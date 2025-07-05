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

### Available Endpoints

#### System
- `GET /health` - Health check with database status
- `GET /` - API information

#### Authentication (Coming in Milestone 4)
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login

#### Clients (Coming in Milestone 5)
- `GET /api/v1/clients` - List all clients
- `POST /api/v1/clients` - Create new client
- `GET /api/v1/clients/:id` - Get client details
- `PUT /api/v1/clients/:id` - Update client
- `DELETE /api/v1/clients/:id` - Delete client

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

2. In Bruno:
   - Navigate to the request you want to test
   - Click "Send" to execute the request
   - Check the "Tests" tab to see automated test results
   - All tests should show green checkmarks when passing

### Available Test Collections

Currently implemented:
- **System**
  - Health Check - Verifies server and database health
  - API Info - Returns available endpoints

Coming soon:
- **Auth** - User registration and login
- **Clients** - Client CRUD operations
- **Projects** - Project management
- **Allocations** - Weekly hour planning
- **Time Entries** - Daily time tracking
- **Reports** - Analytics and summaries

### Understanding Test Results

Each Bruno request includes automated tests that verify:
- HTTP status codes
- Response structure
- Required fields
- Data validity

Failed tests will show in red with details about what went wrong.

## Database

The application uses SQLite for data persistence with automatic migrations on startup.

### Database Location
- Development: `./data/timetracker.db`
- Docker: Mounted at `/app/data/timetracker.db`

### Schema

The database includes the following tables:

#### users
- `id` (UUID) - Primary key
- `username` (string) - Unique username for login
- `password` (string) - Bcrypt hashed password
- `email` (string) - Unique email address
- `full_name` (string) - User's full name
- `is_active` (boolean) - Account status
- `created_at`, `updated_at`, `deleted_at` - Timestamps

#### clients
- `id` (UUID) - Primary key
- `name` (string) - Client company name
- `code` (string) - Unique client code
- `email` (string) - Contact email
- `phone` (string) - Contact phone
- `address` (string) - Client address
- `is_active` (boolean) - Client status
- `user_id` (UUID) - Owner user reference

#### projects
- `id` (UUID) - Primary key
- `name` (string) - Project name
- `code` (string) - Unique project code
- `description` (string) - Project details
- `status` (enum) - active, on_hold, completed, cancelled
- `billable_rate` (float) - Hourly rate
- `currency` (string) - Currency code (default: USD)
- `start_date` (date) - Project start
- `end_date` (date) - Project end (nullable)
- `client_id` (UUID) - Client reference
- `user_id` (UUID) - Owner user reference

#### allocations
- `id` (UUID) - Primary key
- `project_id` (UUID) - Project reference
- `user_id` (UUID) - User reference
- `week_starting` (date) - Monday of the week
- `hours` (float) - Allocated hours for the week
- `notes` (string) - Additional notes

#### time_entries
- `id` (UUID) - Primary key
- `project_id` (UUID) - Project reference
- `user_id` (UUID) - User reference
- `date` (date) - Entry date
- `hours` (float) - Hours worked
- `description` (string) - Work description
- `is_billable` (boolean) - Billable flag

### Viewing the Database

To inspect the database directly:

```bash
# Install SQLite if needed
brew install sqlite

# View all tables
sqlite3 data/timetracker.db ".tables"

# View schema
sqlite3 data/timetracker.db ".schema"

# Query data (example)
sqlite3 data/timetracker.db "SELECT * FROM users;"
```

## Project Structure

```
consultant-time-tracker/
├── cmd/
│   └── server/         # Application entrypoint
│       └── main.go
├── internal/           # Private application code
│   ├── models/         # Database models
│   │   ├── base.go    # Base model with UUID
│   │   ├── user.go    # User model
│   │   ├── client.go  # Client model
│   │   ├── project.go # Project model
│   │   ├── allocation.go # Allocation model
│   │   └── time_entry.go # Time entry model
│   ├── database/       # Database configuration
│   │   ├── connection.go # Database connection
│   │   └── migrations.go # Migration logic
│   ├── handlers/       # HTTP handlers (Coming soon)
│   ├── middleware/     # HTTP middleware (Coming soon)
│   ├── graphql/        # GraphQL schema (Coming soon)
│   └── services/       # Business logic (Coming soon)
├── docker/             # Docker configurations
│   ├── Dockerfile      # Production image
│   ├── Dockerfile.dev  # Development image
│   └── docker-compose.yml
├── bruno-collections/  # API test collections
│   └── TimeTracker-REST/
│       ├── environments/
│       ├── system/
│       └── bruno.json
├── data/              # Database storage (git-ignored)
├── configs/           # Configuration files
├── scripts/           # Utility scripts
├── go.mod            # Go module definition
├── go.sum            # Go module checksums
├── .env.example      # Environment template
├── .gitignore        # Git ignore rules
├── Makefile          # Build commands
└── README.md         # This file
```

## Environment Variables

Create a `.env` file from `.env.example`:

```bash
APP_ENV=development
APP_PORT=8080
APP_HOST=0.0.0.0
DB_PATH=./data/timetracker.db
GRAPHQL_PLAYGROUND=true
```

## Development Workflow

1. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes and test:
   ```bash
   make docker-dev
   # Make changes, test with Bruno
   ```

3. Commit and push:
   ```bash
   git add .
   git commit -m "Description of changes"
   git push -u origin feature/your-feature-name
   ```

4. Create a Pull Request on GitHub

## Troubleshooting

### Port Already in Use
If port 8080 is already in use:
```bash
# Find process using port 8080
lsof -i :8080
# Kill the process or change APP_PORT in .env
```

### Database Issues
If you encounter database errors:
```bash
# Stop containers
make docker-down
# Remove database
rm -f data/timetracker.db*
# Restart
make docker-dev
```

### Docker Issues
If Docker containers won't start:
```bash
# Clean up containers
docker system prune
# Rebuild without cache
docker-compose -f docker/docker-compose.yml build --no-cache
```
