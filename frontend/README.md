# Stanza Bonanza — frontend

React + Vite + TypeScript client for [Stanza Bonanza](../README.md). React Query,
Zustand, Tailwind, react-markdown; WebAuthn passkeys via `@simplewebauthn/browser`.

## Commands

```bash
pnpm install
pnpm dev        # dev server on :5173
pnpm build      # type-check + production build
pnpm lint
pnpm test       # vitest (watch); `pnpm exec vitest run` for a single run
```

The dev server proxies the API to the Go backend on `:8080`. See the root
[README](../README.md) for the full stack and env setup.
