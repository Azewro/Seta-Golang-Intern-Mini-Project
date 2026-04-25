# Frontend (React)

Simple React frontend for Stage 1 Auth & User Management.

## Features

- Register (`/register`) with email verification required
- Login (`/login`)
- Verify email page (`/verify-email?token=...`)
- Resend verification email from login if account is not verified
- Logout with backend revoke session
- View own profile (`/profile`)
- Manager-only users list (`/users`)

## Prerequisites

- Backend auth service running on `http://localhost:8080`
- Node.js 18+
- Root env file `C:\Users\admin\Downloads\Seta-Golang-Intern-Mini-Project\.env.frontend` configured

## Run

```powershell
Set-Location "C:\Users\admin\Downloads\Seta-Golang-Intern-Mini-Project\frontend"
npm install
npm run dev
```

Open: `http://localhost:5173`

## Notes

- Frontend config is loaded from root `.env.frontend`.
- Dev server uses Vite proxy from `/api` to backend URL from `VITE_API_BASE_URL`.
- JWT token is stored in `localStorage` key: `seta_access_token`.
