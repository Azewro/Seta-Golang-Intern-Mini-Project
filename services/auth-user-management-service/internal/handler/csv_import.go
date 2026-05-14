package handler

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode/utf8"

	"auth-user-management-service/internal/usecase"
)

const (
	maxImportUploadBytes = 3 * 1024 * 1024
	formImportFileField  = "file"
)

var errImportEmptyFile = errors.New("csv file is empty")

// parseImportCSV reads CSV with comma delimiter; header row must name columns username, email, password; role is optional.
func parseImportCSV(r io.Reader) ([]usecase.ImportUserRow, error) {
	cr := csv.NewReader(r)
	cr.TrimLeadingSpace = true
	cr.ReuseRecord = false

	header, err := cr.Read()
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil, errImportEmptyFile
		}
		return nil, fmt.Errorf("invalid csv header: %w", err)
	}
	if len(header) == 0 {
		return nil, errImportEmptyFile
	}

	col := make(map[string]int)
	for i, raw := range header {
		name := normalizeCSVHeader(raw)
		if name == "" {
			continue
		}
		if _, dup := col[name]; dup {
			return nil, fmt.Errorf("duplicate header column %q", name)
		}
		col[name] = i
	}

	for _, req := range []string{"username", "email", "password"} {
		if _, ok := col[req]; !ok {
			return nil, fmt.Errorf("missing required column %q", req)
		}
	}

	for name := range col {
		if name != "username" && name != "email" && name != "password" && name != "role" {
			return nil, fmt.Errorf("unknown column %q (allowed: username, email, password, role)", name)
		}
	}

	roleIdx, hasRoleCol := col["role"]

	var rows []usecase.ImportUserRow
	lineNo := 1

	for {
		rec, err := cr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", lineNo+1, err)
		}
		lineNo++

		if rowIsAllBlank(rec) {
			continue
		}

		u := csvField(rec, col["username"])
		e := csvField(rec, col["email"])
		p := csvField(rec, col["password"])
		role := ""
		if hasRoleCol {
			role = csvField(rec, roleIdx)
		}

		rows = append(rows, usecase.ImportUserRow{
			RowNumber: lineNo,
			Username:  u,
			Email:     e,
			Password:  p,
			Role:      role,
		})
	}

	return rows, nil
}

func normalizeCSVHeader(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	// UTF-8 BOM on first header cell
	if len(s) >= 3 && s[0] == 0xEF && s[1] == 0xBB && s[2] == 0xBF {
		s = s[3:]
	}
	if r, sz := utf8.DecodeRuneInString(s); r == '\ufeff' {
		s = s[sz:]
	}
	return strings.ToLower(strings.TrimSpace(s))
}

func csvField(rec []string, idx int) string {
	if idx < 0 || idx >= len(rec) {
		return ""
	}
	return rec[idx]
}

func rowIsAllBlank(rec []string) bool {
	for _, cell := range rec {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}
