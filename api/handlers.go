package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type createCaseRequest struct {
	CaseID       string `json:"case_id"`
	IBAN         string `json:"iban"`
	AESReference string `json:"aes_reference"`
}

type caseRecord struct {
	CaseID       string `json:"case_id"`
	IBAN         string `json:"iban"`
	AESReference string `json:"aes_reference"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}

type caseSummary struct {
	CaseID       string `json:"case_id"`
	IBAN         string `json:"iban"`
	AESReference string `json:"aes_reference"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	PortalLink   string `json:"portal_link"`
	Expired      bool   `json:"expired"`
}

type server struct {
	db           *sql.DB
	tokenSecret  string
	portalOrigin string
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func (s *server) createCase(w http.ResponseWriter, r *http.Request) {
	var req createCaseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if req.IBAN == "" || req.AESReference == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "iban and aes_reference are required"})
		return
	}

	caseID := strings.TrimSpace(req.CaseID)
	if caseID == "" {
		var n int
		s.db.QueryRow("SELECT COUNT(*) FROM cases").Scan(&n)
		caseID = fmt.Sprintf("CASE-2026-%03d", n+1)
	}

	token, err := generateToken(s.tokenSecret)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
		return
	}

	_, err = s.db.Exec(
		`INSERT INTO cases (case_id, iban, aes_reference, issue_type, issue_description, status, token, token_used)
		 VALUES (?, ?, ?, '', '', 'open', ?, 0)`,
		caseID, req.IBAN, req.AESReference, token,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "case_id already exists"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create case"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"case_id":     caseID,
		"portal_link": fmt.Sprintf("%s/case?token=%s", s.portalOrigin, token),
	})
}

func (s *server) listCases(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(
		`SELECT case_id, iban, aes_reference, status, created_at, token, token_used FROM cases ORDER BY created_at DESC`,
	)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "query failed"})
		return
	}
	defer rows.Close()

	var list []caseSummary
	for rows.Next() {
		var c caseSummary
		var token string
		var tokenUsed int
		if err := rows.Scan(&c.CaseID, &c.IBAN, &c.AESReference, &c.Status, &c.CreatedAt, &token, &tokenUsed); err != nil {
			continue
		}
		c.PortalLink = fmt.Sprintf("%s/case?token=%s", s.portalOrigin, token)
		c.Expired = tokenUsed == 1
		list = append(list, c)
	}
	if list == nil {
		list = []caseSummary{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (s *server) expireCase(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CaseID string `json:"case_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CaseID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "case_id required"})
		return
	}
	res, err := s.db.Exec(`UPDATE cases SET token_used = 1 WHERE case_id = ?`, req.CaseID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "update failed"})
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "case not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (s *server) getCase(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	var rec caseRecord
	var tokenUsed int
	err := s.db.QueryRow(
		`SELECT case_id, iban, aes_reference, status, created_at, token_used FROM cases WHERE token = ?`, token,
	).Scan(&rec.CaseID, &rec.IBAN, &rec.AESReference, &rec.Status, &rec.CreatedAt, &tokenUsed)

	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "lookup failed"})
		return
	}
	if tokenUsed == 1 {
		writeJSON(w, http.StatusGone, map[string]string{"error": "expired"})
		return
	}

	writeJSON(w, http.StatusOK, rec)
}
