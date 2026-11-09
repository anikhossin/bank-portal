package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type createCaseRequest struct {
	CaseID           string `json:"case_id"`
	IBAN             string `json:"iban"`
	AESReference     string `json:"aes_reference"`
	IssueType        string `json:"issue_type"`
	IssueDescription string `json:"issue_description"`
}

type caseRecord struct {
	CaseID           string `json:"case_id"`
	IBAN             string `json:"iban"`
	AESReference     string `json:"aes_reference"`
	IssueType        string `json:"issue_type"`
	IssueDescription string `json:"issue_description"`
	Status           string `json:"status"`
	CreatedAt        string `json:"created_at"`
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
	req.IssueType = strings.TrimSpace(req.IssueType)
	if req.IBAN == "" || req.AESReference == "" || req.IssueDescription == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "iban, aes_reference and issue_description are required"})
		return
	}
	if req.IssueType != "failed_transfer" && req.IssueType != "dispute" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "issue_type must be failed_transfer or dispute"})
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
		 VALUES (?, ?, ?, ?, ?, 'open', ?, 0)`,
		caseID, req.IBAN, req.AESReference, req.IssueType, req.IssueDescription, token,
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

func (s *server) getCase(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	var rec caseRecord
	err := s.db.QueryRow(
		`SELECT case_id, iban, aes_reference, issue_type, issue_description, status, created_at
		 FROM cases WHERE token = ?`, token,
	).Scan(&rec.CaseID, &rec.IBAN, &rec.AESReference, &rec.IssueType, &rec.IssueDescription, &rec.Status, &rec.CreatedAt)

	if err == sql.ErrNoRows {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "lookup failed"})
		return
	}

	writeJSON(w, http.StatusOK, rec)
}
