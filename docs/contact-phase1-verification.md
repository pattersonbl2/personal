# Contact form Phase 1 — manual verification

After deploying Pages + Cloud Run with `TURNSTILE_SECRET` set:

1. **Happy path (JS)** — open `/contact/`, wait ~3s, submit a real message. Expect inline success; email arrives via Resend.
2. **Happy path (no-JS AJAX)** — disable the contact.js script only (or force classic POST). Expect branded HTML success page from the API.
3. **Bad Turnstile token** — clear `cf-turnstile-response` in DevTools before submit. Expect inline/API error; no email.
4. **Honeypot** — fill hidden `website` field. Expect fake success; no email; log `reason=spam detail=honeypot`.
5. **Too fast** — submit immediately with `form_ts` set to now. Expect fake success; log `timing_too_fast`.
6. **Rate limit** — exceed 5 POSTs/hour from one IP. Expect 429; log `ratelimit: reject`.
7. **Resume PDF** — on `/resume/`, complete Turnstile and click Download PDF. Expect PDF download; without Turnstile expect 403.

Backend unit tests: `cd backend && go test ./...`

Spin validate (needs `TURNSTILE_SECRET` in the local shell, not pasted in chat):
`"$TURNSTILE_SPIN_SCRIPTS/validate.sh" --sitekey 0x4AAAAAAD-Wku5l0KWMgbF-`
