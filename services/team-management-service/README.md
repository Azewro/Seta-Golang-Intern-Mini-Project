# Team Management Service

Stage 1 implementation for team creation and team membership management.

## Scope

- Global user role (`manager` / `member`) is owned by auth-user-management-service.
- This service manages team records and team-scoped memberships only.
- Team mutation endpoints require JWT with global `manager` role.
- Only `mainManagerUserId` can add/remove managers inside a team.

## Data Model

- `teams`
  - `id`
  - `team_name`
  - `main_manager_user_id`
- `team_memberships`
  - `team_id`
  - `user_id`
  - `membership_role` (`manager` or `member`)

## Environment

The service uses root `.env.backend` (same pattern as auth service):

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=seta_miniproject_db
JWT_SECRET=super-secret-key-for-jwt
TEAM_SERVICE_PORT=8081
```

## Run

```powershell
Set-Location "C:\Users\admin\Downloads\Seta-Golang-Intern-Mini-Project\services\team-management-service"
go mod tidy
go run ./cmd/main.go
```

## API

Base path: `/api/v1`

- `POST /teams` (manager) `{ "teamName": "Core Team" }`
- `GET /teams/my` (authenticated)
- `GET /teams/:teamId` (team member)
- `POST /teams/:teamId/members` (manager in team) `{ "userId": 2 }`
- `DELETE /teams/:teamId/members/:userId` (manager in team)
- `POST /teams/:teamId/managers` (main manager only) `{ "userId": 3 }`
- `DELETE /teams/:teamId/managers/:userId` (main manager only)

Use JWT bearer token from auth-user-management-service login response.
