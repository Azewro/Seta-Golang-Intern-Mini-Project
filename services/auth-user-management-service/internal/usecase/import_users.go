package usecase

import (
	"errors"
	"net/mail"
	"strings"
	"sync"

	"auth-user-management-service/internal/domain"
	"auth-user-management-service/pkg/utils"

	"gorm.io/gorm"
)

const (
	importMaxErrorsListed = 50
	defaultImportWorkers  = 5
	maxImportWorkersCap   = 64
)

// ImportUserRow is one CSV data row (header is line 1; first data row is line 2).
type ImportUserRow struct {
	RowNumber int
	Username  string
	Email     string
	Password  string
	Role      string
}

// ImportUserError describes a failed row for the capped error list.
type ImportUserError struct {
	RowNumber int    `json:"rowNumber"`
	Email     string `json:"email,omitempty"`
	Message   string `json:"message"`
}

// ImportUsersResponse is the aggregate result of a bulk import run.
type ImportUsersResponse struct {
	Success           int               `json:"success"`
	Failed            int               `json:"failed"`
	Errors            []ImportUserError `json:"errors"`
	ErrorsTruncated   bool              `json:"errorsTruncated"`
}

type importRowOutcome struct {
	ok      bool
	row     int
	email   string
	message string
}

// ImportUsers creates users concurrently using a worker pool. Imported users are email-verified by default.
func (u *authUsecaseImpl) ImportUsers(rows []ImportUserRow) *ImportUsersResponse {
	out := &ImportUsersResponse{Errors: make([]ImportUserError, 0, importMaxErrorsListed)}
	if len(rows) == 0 {
		return out
	}

	workers := readIntEnvWithDefault("IMPORT_WORKERS", defaultImportWorkers)
	if workers < 1 {
		workers = defaultImportWorkers
	}
	if workers > maxImportWorkersCap {
		workers = maxImportWorkersCap
	}
	if workers > len(rows) {
		workers = len(rows)
	}

	jobs := make(chan ImportUserRow)
	var wg sync.WaitGroup

	var batchMu sync.Mutex
	seenEmail := make(map[string]struct{})

	var aggMu sync.Mutex

	workerFn := func() {
		defer wg.Done()
		for row := range jobs {
			o := u.processImportRow(&batchMu, seenEmail, row)
			aggMu.Lock()
			if o.ok {
				out.Success++
			} else {
				out.Failed++
				if len(out.Errors) < importMaxErrorsListed {
					out.Errors = append(out.Errors, ImportUserError{
						RowNumber: o.row,
						Email:     o.email,
						Message:   o.message,
					})
				} else {
					out.ErrorsTruncated = true
				}
			}
			aggMu.Unlock()
		}
	}

	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go workerFn()
	}

	for i := range rows {
		jobs <- rows[i]
	}
	close(jobs)
	wg.Wait()

	return out
}

func (u *authUsecaseImpl) processImportRow(batchMu *sync.Mutex, seen map[string]struct{}, row ImportUserRow) importRowOutcome {
	username := strings.TrimSpace(row.Username)
	emailRaw := strings.TrimSpace(row.Email)
	password := row.Password
	role := strings.TrimSpace(row.Role)

	emailNorm := strings.ToLower(emailRaw)

	if username == "" {
		return failOutcome(row.RowNumber, emailNorm, "username is required")
	}
	if emailNorm == "" {
		return failOutcome(row.RowNumber, "", "email is required")
	}
	if !isValidEmailAddress(emailNorm) {
		return failOutcome(row.RowNumber, emailNorm, "invalid email format")
	}
	if strings.TrimSpace(password) == "" {
		return failOutcome(row.RowNumber, emailNorm, "password is required")
	}
	if len(strings.TrimSpace(password)) < 8 {
		return failOutcome(row.RowNumber, emailNorm, "password must be at least 8 characters")
	}

	if role == "" {
		role = "member"
	}
	if role != "manager" && role != "member" {
		return failOutcome(row.RowNumber, emailNorm, "role must be manager or member")
	}

	batchMu.Lock()
	if _, dup := seen[emailNorm]; dup {
		batchMu.Unlock()
		return failOutcome(row.RowNumber, emailNorm, "duplicate email in import file")
	}
	seen[emailNorm] = struct{}{}
	batchMu.Unlock()

	existing, err := u.userRepo.FindByEmail(emailNorm)
	if err == nil && existing != nil {
		return failOutcome(row.RowNumber, emailNorm, "email already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return failOutcome(row.RowNumber, emailNorm, "database error while checking email")
	}

	hashed, err := utils.HashPassword(strings.TrimSpace(password))
	if err != nil {
		return failOutcome(row.RowNumber, emailNorm, "failed to hash password")
	}

	user := &domain.User{
		Username:   username,
		Email:      emailNorm,
		Password:   hashed,
		Role:       role,
		IsVerified: true,
	}

	if err := u.userRepo.CreateUser(user); err != nil {
		if isDuplicateDBError(err) {
			return failOutcome(row.RowNumber, emailNorm, "email already exists")
		}
		return failOutcome(row.RowNumber, emailNorm, "failed to create user")
	}

	return importRowOutcome{ok: true}
}

func failOutcome(row int, email, msg string) importRowOutcome {
	return importRowOutcome{ok: false, row: row, email: email, message: msg}
}

func isValidEmailAddress(s string) bool {
	addr, err := mail.ParseAddress(s)
	if err != nil || addr.Address == "" {
		return false
	}
	return strings.EqualFold(addr.Address, s)
}

func isDuplicateDBError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "duplicate") || strings.Contains(msg, "unique")
}
