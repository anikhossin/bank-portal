# Bank Case Portal (Demo)

A demo of a one-time-link case-sharing workflow for a fictional bank ("Northfield Bank"). Three standalone apps:

- [`api/`](api) — Go + SQLite backend. Case storage, HMAC-signed one-time tokens, CORS.
- [`admin/`](admin) — Svelte + Tailwind. Create a case, get a shareable link.
- [`portal/`](portal) — Svelte + Tailwind. Customer opens the link once to view case details.

## Run all three

```sh
# terminal 1
cd api && cp .env.example .env && go run .

# terminal 2
cd admin && cp .env.example .env && npm install && npm run dev   # :5174

# terminal 3
cd portal && cp .env.example .env && npm install && npm run dev # :5173
```

Open the admin app, fill out a case, copy the generated link, open it in the portal app (or a different browser/incognito window, since the token is single-use).

See each folder's README for details. No secrets are hardcoded — see each `.env.example`.
