# Case Portal — Admin

Svelte + Tailwind admin UI for creating case links. No auth, single page.

## Setup

```sh
cp .env.example .env
npm install
npm run dev
```

Runs on **http://localhost:5174**. Requires the [`api`](../api) server running (default `http://localhost:8080`; override with `VITE_API_URL` in `.env`).

All validation and business logic live in the API — this app only collects form input and displays whatever the server returns.
