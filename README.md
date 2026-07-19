# Stanza Bonanza

Unlike other poetry sites, Stanza Bonanza encourages open collaboration between poets. Users can create a new poem, write a stanza or two, and then wait for other users to extend that poem, adding their own personal spin and enriching the piece in the process. Before long, a diverse tapestry comes to life!

Poems on Stanza Bonanza are more than the sum of their parts!

If poets desire more control over their work, they can specify certain constraints.

---

## Stack

- **Backend**: Go 1.25, chi router, Postgres (raw SQL migrations), WebAuthn passkeys
  and magic-link auth. In `backend/`.
- **Frontend**: React + Vite + TypeScript, React Query, Zustand, Tailwind,
  react-markdown. In `frontend/` (pnpm).

## Develop

```bash
make dev        # Postgres (docker-compose) + backend + frontend
make migrate    # apply DB migrations (needs DATABASE_URL)
make test       # backend + frontend tests
make lint
```

Or per side: `make dev-backend`, `make dev-frontend`. Frontend on `:5173`, backend on `:8080`.

## Environment

Backend config lives in `backend/internal/config/config.go`, set via env or a `.env`:

| Var | Required | Default |
|---|---|---|
| `DATABASE_URL` | yes | — |
| `SESSION_SECRET` | yes | — |
| `PORT` | no | `8080` |
| `ALLOWED_ORIGINS` | no | `http://localhost:5173` |
| `WEBAUTHN_RP_ID` / `WEBAUTHN_RP_NAME` / `WEBAUTHN_RP_ORIGINS` | no | localhost defaults |
| `RESEND_API_KEY` | for magic-link email | — (no key: link is logged, not sent) |
| `EMAIL_FROM` | no | `Stanza Bonanza <onboarding@resend.dev>` |
| `MAGIC_LINK_BASE_URL` | no | `http://localhost:5173/auth/verify` |

## Deploy

The backend deploys to Fly.io from `backend/fly.toml` (`fly deploy`); the frontend
builds to static assets (`pnpm build`).

## Layout

```
backend/
  cmd/server      entrypoint
  internal/       handler, service, repository, middleware, domain, config
  migrations/     SQL migrations
frontend/
  src/            components, pages, hooks, store
```
