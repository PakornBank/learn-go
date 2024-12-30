# Go Authentication API

An authentication system built with Go, featuring JWT authentication, PostgreSQL database, and clean architecture.

## Features
- User registration and login
- JWT-based authentication
- Protected routes
- PostgreSQL database with GORM
- Clean architecture (Handler -> Service -> Repository layers)
- Docker support for PostgreSQL

## Tech Stack
- Go 1.21+
- Gin (Web Framework)
- GORM (ORM)
- PostgreSQL (Database)
- JWT (Authentication)
- Docker & Docker Compose

## Project Structure
```
.
├── cmd
│   └── api
│       └── main.go
├── internal
│   ├── config
│   │   └── config.go
│   ├── database
│   │   └── database.go
│   ├── handler
│   │   └── auth_handler.go
│   ├── middleware
│   │   └── auth_middleware.go
│   ├── model
│   │   └── user.go
│   ├── repository
│   │   └── user_repository.go
│   └── service
│       └── auth_service.go
├── docker-compose.yml
├── .env
├── .env.example
├── .gitignore
├── go.mod
└── README.md
```

## Prerequisites
- Go 1.21 or higher
- Docker and Docker Compose
- Make (optional)

## Getting Started

1. Clone the repository
```bash
git clone https://github.com/PakornBank/learn-go.git
cd learn-go
```

2. Copy environment file and update values
```bash
cp .env.example .env
```

3. Start PostgreSQL using Docker
```bash
docker-compose up -d
```

4. Install dependencies
```bash
go mod tidy
```

5. Run the application
```bash
go run cmd/api/main.go
```

The server will start at `http://localhost:8080`

## Environment Variables
Create a `.env` file in the root directory:

```env
# .env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_auth_db
DB_PORT=5432
SERVER_PORT=8080
JWT_SECRET=your-super-secret-key-here
```

## API Endpoints

### Public Routes
- `POST /api/register` - Register a new user
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "full_name": "John Doe"
  }'
```

- `POST /api/login` - Login and get JWT token
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Protected Routes (Requires JWT Token)
- `GET /api/profile` - Get user profile
```bash
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Testing
Run all tests:
```bash
go test ./...
```

## Development

### Database Management
Start PostgreSQL:
```bash
docker-compose up -d
```

Stop PostgreSQL:
```bash
docker-compose down
```

Reset Database:
```bash
docker-compose down -v
docker-compose up -d
```

### Common Issues

1. Database Connection
If you can't connect to the database, ensure:
- PostgreSQL container is running (`docker ps`)
- Database credentials in `.env` match `docker-compose.yml`
- Database port is not in use

2. JWT Token
If authentication fails:
- Check token expiration
- Verify JWT_SECRET in `.env`
- Ensure token format: `Bearer YOUR_TOKEN`

## Contributing
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License
This project is licensed under the MIT License - see the LICENSE file for details

## Acknowledgments
- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io)
- [JWT Go](https://github.com/golang-jwt/jwt)
