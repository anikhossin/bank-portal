# Case Portal — Customer

Svelte + Tailwind customer-facing portal. Reads a one-time `token` from the URL and shows the case once.

## Setup

```sh
cp .env.example .env
npm install
npm run dev
```

Runs on **http://localhost:5173**. Requires the [`api`](../api) server running (default `http://localhost:8080`; override with `VITE_API_URL` in `.env`).

Visit `http://localhost:5173/case?token=<token>` (get a token from the admin app). States are driven entirely by the API's response:

- 404 → "This link is invalid"
- 410 → "This link has already been opened"
- 200 → two-step reveal (case ID → full details), each field has its own copy button
