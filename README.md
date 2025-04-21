# Go Echo API with Turso, GORM, and JWT Authentication

A RESTful API built with Go, using Echo framework, Turso as a database, GORM as the ORM, and JWT for authentication.

## Features

- RESTful API endpoints for user management
- JWT-based authentication system
- Integration with Turso (LibSQL) database
- GORM ORM for database operations
- Environment-based configuration
- Structured logging

## Prerequisites

- Go 1.21 or higher
- Turso CLI and account (to set up your database)

## Setup

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/go-echo-turso-api.git
   cd go-echo-turso-api
   ```

2. Install dependencies:
   ```
   go mod download
   ```

3. Set up your Turso database:
   ```
   turso db create my-api-db
   turso db tokens create my-api-db
   ```

4. Create a `.env` file based on `.env.example`:
   ```
   cp .env.example .env
   ```

5. Update the `.env` file with your Turso database URL, auth token, and JWT secret:
   ```
   TURSO_DATABASE_URL=libsql://my-api-db.turso.io
   TURSO_AUTH_TOKEN=your-auth-token
   JWT_SECRET_KEY=your-secure-secret-key
   ```

## Running the application

Start the server:
```
go run main.go
```

The API will be available at `http://localhost:8080`.

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Authenticate and get a JWT token

### Users (Protected Routes - Require JWT Token)

- `GET /api/users` - Get all users
- `GET /api/users/:id` - Get a specific user
- `POST /api/users` - Create a new user
- `PUT /api/users/:id` - Update a user
- `DELETE /api/users/:id` - Delete a user

### Health Check

- `GET /health` - Check if the API is up and running

## Authentication

Most API endpoints are protected and require authentication. You need to:

1. Register a user or login to get a JWT token
2. Include the token in the Authorization header of subsequent requests:
   ```
   Authorization: Bearer your-jwt-token
   ```

## Example Requests

### Register a New User

```
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login

```
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get All Users (Authenticated)

```
curl -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer your-jwt-token"
```

## Project Structure

```
.
├── config/                 # Configuration packages
│   ├── database.go         # Database connection setup
│   └── jwt.go              # JWT configuration and utilities
├── controllers/            # HTTP request handlers
│   ├── auth_controller.go  # Authentication controller
│   └── user_controller.go  # User controller
├── middleware/             # Custom middleware
│   └── jwt_middleware.go   # JWT authentication middleware
├── models/                 # Database models
│   ├── migrate.go          # Database migration
│   └── user.go             # User model and DTOs
├── .env                    # Environment variables (gitignored)
├── .env.example            # Example environment variables
├── go.mod                  # Go module file
├── go.sum                  # Go module checksums
├── main.go                 # Application entry point
└── README.md               # This file
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License.