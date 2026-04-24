# Frontend (React)

Simple React frontend for Stage 1 Auth & User Management.

## Features

- Register (`/register`)
- Login (`/login`)
- Logout with backend revoke session
- View own profile (`/profile`)
- Manager-only users list (`/users`)

## Prerequisites

- Backend auth service running on `http://localhost:8080`
- Node.js 18+

## Run

```powershell
Set-Location "C:\Users\admin\Downloads\Seta-Golang-Intern-Mini-Project\frontend"
npm install
npm run dev
```

Open: `http://localhost:5173`

## Notes

- Dev server uses Vite proxy from `/api` to backend `http://localhost:8080`.
- JWT token is stored in `localStorage` key: `seta_access_token`.

