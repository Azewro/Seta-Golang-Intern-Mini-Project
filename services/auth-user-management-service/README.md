# Auth User Management Service

Stage 1 implementation for authentication and user management.

## Features

- Register user with strict role validation (`manager` or `member`)
- Login with JWT token generation
- Logout with real token revocation via `sessions` table
- View own profile (`/users/me`)
- Manager-only user list (`/users`)

## Environment

Create `.env` in this folder:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=seta_miniproject_db
PORT=8080
JWT_SECRET=super-secret-key-for-jwt
JWT_EXPIRES_HOURS=24
```

## Run

```powershell
go mod tidy
go run ./cmd/main.go
```

## Quick API test

```powershell
# Register manager
curl -Method POST "http://localhost:8080/api/v1/auth/register" -ContentType "application/json" -Body '{"username":"admin","email":"admin@example.com","password":"password123","role":"manager"}'

# Login
$login = curl -Method POST "http://localhost:8080/api/v1/auth/login" -ContentType "application/json" -Body '{"email":"admin@example.com","password":"password123"}'
$token = ($login.Content | ConvertFrom-Json).accessToken

# Manager list users
curl -Method GET "http://localhost:8080/api/v1/users" -Headers @{ Authorization = "Bearer $token" }

# Me
curl -Method GET "http://localhost:8080/api/v1/users/me" -Headers @{ Authorization = "Bearer $token" }

# Logout
curl -Method POST "http://localhost:8080/api/v1/auth/logout" -Headers @{ Authorization = "Bearer $token" }
```

