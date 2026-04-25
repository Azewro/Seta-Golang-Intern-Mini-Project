# Auth User Management Service

Stage 1 implementation for authentication and user management.

## Features

- Register user (public sign-up creates `member` role)
- Email verification via SMTP link (token expires in 5 minutes)
- Resend verification email endpoint
- Login with JWT token generation (only for verified accounts)
- Logout with real token revocation via `sessions` table
- View own profile (`/users/me`)
- Manager-only user list (`/users`)

## Environment

Use the root backend env file:

```powershell
Set-Location "C:\Users\admin\Downloads\Seta-Golang-Intern-Mini-Project"
# Edit .env.backend with your values
```

Required keys:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=seta_miniproject_db
PORT=8080
APP_BASE_URL=http://localhost:8080
JWT_SECRET=super-secret-key-for-jwt
JWT_EXPIRES_HOURS=24
EMAIL_VERIFY_TOKEN_TTL_MINUTES=5
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email@example.com
SMTP_PASSWORD=your-app-password
SMTP_FROM_EMAIL=your-email@example.com
SMTP_FROM_NAME=Seta App
```

## Run

```powershell
go mod tidy
go run ./cmd/main.go
```

## Quick API test

```powershell
# Register (role is assigned as member on server)
curl -Method POST "http://localhost:8080/api/v1/auth/register" -ContentType "application/json" -Body '{"username":"user1","email":"user1@example.com","password":"password123"}'

# Verify from email link (example)
curl -Method GET "http://localhost:8080/api/v1/auth/verify-email?token=<token-from-email>"

# Resend verification email
curl -Method POST "http://localhost:8080/api/v1/auth/resend-verification" -ContentType "application/json" -Body '{"email":"user1@example.com"}'

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
