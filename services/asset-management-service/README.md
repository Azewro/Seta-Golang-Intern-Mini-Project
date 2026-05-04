# Asset Management & Sharing Service

Stage 2 Service 1 implementation for managing folders, notes, and sharing rules.

## Scope

- Folder and note CRUD.
- Share folder/note with `read` or `write`.
- Folder-share inheritance for notes.
- Manager read-only oversight for assets owned by their team members.

## Environment

Uses root `.env.backend`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=seta_miniproject_db
ASSET_SERVICE_PORT=8082
AUTH_SERVICE_URL=http://localhost:8080
TEAM_SERVICE_URL=http://localhost:8081
```

## Run

```powershell
Set-Location "C:\Users\admin\Downloads\Seta-Golang-Intern-Mini-Project\services\asset-management-service"
go mod tidy
go run ./cmd/main.go
```

## API

Base path: `/api/v1`

- `POST /folders`
- `GET /folders`
- `GET /folders/:folderId`
- `PATCH /folders/:folderId`
- `DELETE /folders/:folderId`
- `POST /folders/:folderId/notes`
- `GET /folders/:folderId/notes`
- `GET /notes/:noteId`
- `PATCH /notes/:noteId`
- `DELETE /notes/:noteId`
- `POST /shares`
- `DELETE /shares/:shareId`
- `GET /shares/received`
- `GET /shares/granted`
