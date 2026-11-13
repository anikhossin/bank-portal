package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type createCaseRequest struct {
	CaseID           string `json:"case_id"`
	IBAN             string `json:"iban"`
	AESReference     string `json:"aes_reference"`
	Verwendungszweck string `json:"verwendungszweck"`
	Betrag           string `json:"betrag"`
}

type caseRecord struct {
	CaseID           string `json:"case_id"`
	IBAN             string `json:"iban"`
	AESReference     string `json:"aes_reference"`
	Verwendungszweck string `json:"verwendungszweck"`
	Betrag           string `json:"betrag"`
	Status           string `json:"status"`
	CreatedAt        string `json:"created_at"`
}

type caseSummary struct {
	CaseID           string `json:"case_id"`
	IBAN             string `json:"iban"`
	AESReference     string `json:"aes_reference"`
	Verwendungszweck string `json:"verwendungszweck"`
	Betrag           string `json:"betrag"`
	Status           string `json:"status"`
	CreatedAt        string `json:"created_at"`
	PortalLink       string `json:"portal_link"`
	Expired          bool   `json:"expired"`
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
	caseID := strings.TrimSpace(req.CaseID)
	if caseID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Fallnummer ist erforderlich."})
		return
	}
	if req.IBAN == "" || req.AESReference == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "iban and aes_reference are required"})
		return
	}

	var existing int
	if err := s.db.QueryRow(`SELECT COUNT(*) FROM cases WHERE case_id = ?`, caseID).Scan(&existing); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not check case id"})
		return
	}
	if existing > 0 {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "Diese Fallnummer existiert bereits. Bitte eine andere Fallnummer verwenden."})
		return
	}

	token, err := generateToken(s.tokenSecret)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not generate token"})
		return
	}

	_, err = s.db.Exec(
		`INSERT INTO cases (case_id, iban, aes_reference, issue_type, issue_description, verwendungszweck, betrag, status, token, token_used)
		 VALUES (?, ?, ?, 'failed_transfer', '', ?, ?, 'open', ?, 0)`,
		caseID, req.IBAN, req.AESReference, req.Verwendungszweck, req.Betrag, token,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "Diese Fallnummer existiert bereits. Bitte eine andere Fallnummer verwenden."})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "could not create case"})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"case_id":     caseID,
		"portal_link": s.portalLink(caseID),
	})
}

func (s *server) listCases(w http.ResponseWriter, r *http.Request) {
	rows, err := s.db.Query(
		`SELECT case_id, iban, aes_reference, verwendungszweck, betrag, status, created_at, token, token_used FROM cases ORDER BY created_at DESC`,
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
		if err := rows.Scan(&c.CaseID, &c.IBAN, &c.AESReference, &c.Verwendungszweck, &c.Betrag, &c.Status, &c.CreatedAt, &token, &tokenUsed); err != nil {
			continue
		}
		c.PortalLink = s.portalLink(c.CaseID)
		c.Expired = tokenUsed == 1
		list = append(list, c)
	}
	if list == nil {
		list = []caseSummary{}
	}
	writeJSON(w, http.StatusOK, list)
}

func (s *server) deleteCase(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CaseID string `json:"case_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.CaseID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "case_id required"})
		return
	}
	res, err := s.db.Exec(`DELETE FROM cases WHERE case_id = ?`, req.CaseID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "delete failed"})
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "case not found"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
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

func (s *server) portalLink(caseID string) string {
	return s.portalOrigin + "/case?c=" + url.QueryEscape(caseID)
}

func (s *server) getCase(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.URL.Query().Get("id"))
	if id == "" {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}

	var rec caseRecord
	var tokenUsed int
	err := s.db.QueryRow(
		`SELECT case_id, iban, aes_reference, verwendungszweck, betrag, status, created_at, token_used FROM cases WHERE case_id = ?`, id,
	).Scan(&rec.CaseID, &rec.IBAN, &rec.AESReference, &rec.Verwendungszweck, &rec.Betrag, &rec.Status, &rec.CreatedAt, &tokenUsed)

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
