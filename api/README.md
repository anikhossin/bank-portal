# Bank Portal API

Go backend for the case portal demo. SQLite storage, HMAC-signed one-time case tokens.

## Setup

```sh
cp .env.example .env   # edit TOKEN_SECRET before any real use
go mod tidy
go run .
```

Runs on `PORT` from `.env` (default `8080`).

## Endpoints

- `POST /api/cases` — create a case. Body:
  ```json
  { "case_id": "CASE-2024-001", "iban": "...", "aes_reference": "...", "issue_type": "failed_transfer", "issue_description": "..." }
  ```
  `case_id` is optional (auto-generated if omitted). Returns `{ case_id, portal_link }`.
- `GET /api/case?token=...` — fetch a case by its one-time token.
  - Valid + unused → 200 with case data, marks token used.
  - Already used → 410.
  - Unknown token → 404.

## Config (`.env`)

| Var | Purpose |
|---|---|
| `PORT` | HTTP port |
| `DB_PATH` | SQLite file path |
| `TOKEN_SECRET` | HMAC signing secret — required, no default |
