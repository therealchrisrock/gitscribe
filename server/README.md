# Teammate API Server

A Golang REST API built with Gin framework, PostgreSQL with GORM, and Firebase Authentication.

## Features

- **Gin Framework**: High-performance HTTP server
- **PostgreSQL with GORM**: Robust database layer with ORM
- **Firebase Authentication**: User authentication with Firebase Auth
- **Migration System**: Database versioning with golang-migrate
- **Repository Pattern**: Clean separation of concerns
- **Environment-based Configuration**: Using .env files
- **CORS Support**: Ready for cross-origin requests
- **Middleware**: Logging, error handling, authentication
- **Clean Project Structure**: Well-organized codebase

## Project Structure

```
server/
├── database/         # Database connection and migrations
├── firebase/         # Firebase authentication integration
├── handlers/         # HTTP handlers
├── middleware/       # Middleware components
├── migrations/       # SQL migration files
├── models/           # Data models
├── repositories/     # Data access layer
├── .env              # Environment configuration
├── go.mod            # Go modules file
├── go.sum            # Go modules checksums
├── main.go           # Application entrypoint
└── README.md         # This file
```

## Prerequisites

- Go 1.21.x (preferably Go 1.21.13)
- PostgreSQL 13 or later
- Firebase project (for authentication)

> **Note:** The project is configured to work with Go 1.21.x. If you're using a newer Go version and encounter compatibility issues, you might need to update the `go` directive in `go.mod` file.

## Setup

1. **Clone the repository**

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure the environment**
   - Copy `.env.example` to `.env`
   - Update the database credentials
   - Configure Firebase authentication

4. **Create PostgreSQL database**
   ```bash
   createdb teammate_db
   ```

5. **Run the server**
   ```bash
   go run main.go
   ```

   The server will start at `http://localhost:8080`

## API Endpoints

### Public Endpoints

- `GET /health` - Health check
- `POST /api/register` - Register a new user
- `GET /api/users` - List all users (would typically be restricted)
- `GET /api/users/:id` - Get user by ID (would typically be restricted)

### Protected Endpoints (require Firebase Authentication)

- `GET /api/me` - Get the current authenticated user
- `PUT /api/users/:id` - Update a user
- `DELETE /api/users/:id` - Delete a user

## Authentication

This API uses Firebase Authentication. Include the Firebase ID token in the `Authorization` header:

```
Authorization: Bearer your-firebase-id-token
```

## Development

### Adding a New Migration

1. Create migration files:
   ```bash
   migrate create -ext sql -dir migrations -seq name_of_migration
   ```

2. Edit the `.up.sql` and `.down.sql` files in the migrations directory

### Running Migrations Manually

Migrations run automatically when the server starts, but you can also run them manually:

```bash
migrate -path migrations -database "postgres://user:password@localhost:5432/teammate_db?sslmode=disable" up
```

## Docker Deployment

This project includes Docker and Docker Compose configurations for easy deployment:

### Using Docker Compose

1. Build and start all services:
   ```bash
   cd server
   docker-compose up -d
   ```

   This will:
   - Start a PostgreSQL database container
   - Build and start the API container
   - Create a network for the containers to communicate
   - Mount volumes for data persistence

2. View logs:
   ```bash
   docker-compose logs -f
   ```

3. Stop all services:
   ```bash
   docker-compose down
   ```

4. To remove volumes when stopping:
   ```bash
   docker-compose down -v
   ```

### Using Docker without Compose

1. Create a network:
   ```bash
   docker network create teammate-network
   ```

2. Start PostgreSQL:
   ```bash
   docker run -d \
     --name teammate-postgres \
     --network teammate-network \
     -p 5432:5432 \
     -e POSTGRES_USER=postgres \
     -e POSTGRES_PASSWORD=your-super-secret-and-long-postgres-password \
     -e POSTGRES_DB=teammate_db \
     -v postgres-data:/var/lib/postgresql/data \
     postgres:15-alpine
   ```

3. Build and start the API:
   ```bash
   docker build -t teammate-api -f dockerfile .
   
   docker run -d \
     --name teammate-api \
     --network teammate-network \
     -p 8080:8080 \
     -e DB_HOST=postgres \
     -e DB_PORT=5432 \
     -e DB_USER=postgres \
     -e DB_PASSWORD=your-super-secret-and-long-postgres-password \
     -e DB_NAME=teammate_db \
     -e DB_SSLMODE=disable \
     -e GIN_MODE=debug \
     -e APP_ENV=development \
     teammate-api
   ```

## Production Deployment

For production:

1. Set `GIN_MODE=release` in the .env file or Docker environment
2. Use a proper Firebase service account
3. Configure proper SSL/TLS for the database connection
4. Use Docker Compose with production overrides
5. Implement proper secrets management instead of hardcoded passwords
