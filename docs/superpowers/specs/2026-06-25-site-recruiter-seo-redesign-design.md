# Site Redesign: Recruiter-First, Terminal-Craftsman, SEO/AEO/GEO

**Date:** 2026-06-25
**Status:** Approved (direction), implementing
**Site:** ark31.info — Hugo + custom theme, Go backend on Cloud Run, deployed via Cloudflare Pages

## Goal

Make the site easy for recruiters to land on, understand, and act on ("hire me"), while
keeping consulting prominent. Replace the generic GitHub-dark clone look with an owned
"terminal-craftsman" identity that feels human-crafted (AI-assisted, not AI-generated).
Add a robust, on-brand CI/CD pipeline hero animation. Maximize discoverability across
classic search (SEO), answer engines (AEO), and generative engines (GEO).

## Approved Direction

- **Audience:** Balanced, recruiter-first. Hero leads with hire-me signals; consulting stays a strong second track.
- **Visual:** Terminal / craftsman — mono display type, terminal motifs, ONE owned signature accent, done with restraint.
- **Animation:** CI/CD pipeline (commit → build → test → scan → deploy → healthy), matching the SRE/GitOps story.

## Design System

Single owned accent replaces the GitHub-blue/leftover-purple inconsistency.

- **Signature accent (amber/gold):** dark `#e3b341`, hover `#f2cc60`; light `#9a6700`, hover `#7d5300`.
  Used for logo, primary CTAs, links, section markers, hover states.
- **Success green (semantic only):** dark `#3fb950`, light `#1a7f37` — pipeline "healthy/passed" states.
- **Surfaces:** keep dark `#0d1117` base; warm the border slightly (`#2a2f37`), lift muted text contrast (`#9aa4b2`).
- **Type:** body stays humanist sans for prose readability. Display type (logo, hero name, nav,
  section titles, page/post titles) becomes mono with tight tracking. Section titles get a mono
  `//` marker prefix. Blog post body headings stay sans (long-form readability).
- All accent usages flow through CSS custom properties — no hard-coded one-offs.

## Homepage (recruiter-first, balanced)

Order: Hero → Impact strip → Primary CTAs → Pipeline terminal → Selected work → Consulting → Recent posts.

- **Hero:** name + role + a one-line value proposition + location/availability line
  ("Charlotte, NC · Remote · open to senior Platform/SRE roles & select consulting").
- **Impact strip:** scannable quantified proof pulled from the resume —
  `50M+ users supported`, `40% cloud cost cut`, `6+ yrs SRE/DevOps`, `19 GitOps apps self-hosted`.
- **Primary CTAs:** "View résumé" (→ /resume/), "Get in touch" (→ /contact/), LinkedIn.
- Keep existing project + consulting + recent-posts sections; tighten copy voice so it reads human.

## CI/CD Pipeline Animation

Rewrite `layouts/partials/terminal.html`. A terminal window animating a realistic deploy:
`git push` → commit line → build → test (N passed) → trivy scan → argocd sync → rollout pods healthy
→ probes/p99 → "deployed in …". Spinners resolve into colored check/cross marks with fake timings;
loops with a pause. Honors `prefers-reduced-motion` (renders final static state, no animation).
Pure vanilla JS, no deps, no layout shift.

## SEO / AEO / GEO

- **robots.txt** (`layouts/robots.txt`): allow all, explicitly welcome AI crawlers
  (GPTBot, OAI-SearchBot, ChatGPT-User, ClaudeBot, anthropic-ai, PerplexityBot, Google-Extended,
  CCBot, Applebot-Extended, Bytespider), and emit `Sitemap:` absolute URL.
- **JSON-LD** (in `head.html`):
  - `WebSite` (site-wide).
  - Home `Person` enriched: `knowsAbout` (Kubernetes, SRE, Terraform, GitOps, observability…),
    `address` (Charlotte, NC), `email`, `alumniOf` (VCU), `description`.
  - Posts: `BlogPosting` — headline, datePublished/modified, author, image, keywords (tags),
    wordCount, `mainEntityOfPage`.
  - `BreadcrumbList` on posts and section/leaf pages.
  - Consulting: `FAQPage` from a new FAQ section.
- **Meta:** `twitter:card` → `summary_large_image`; add `og:image:alt`, `og:locale`,
  `og:image:width/height`; keep canonical.
- **llms.txt** (`static/llms.txt`): concise factual profile + key links for generative engines.
- **Content for AEO:** add a short factual "About" block and a FAQ section on the consulting page
  (extractable Q→A pairs); ensure copy states facts plainly.

## Files Touched

- `assets/css/main.css` — tokens, mono display type, impact strip, CTA buttons, pipeline styles.
- `layouts/index.html` — hero value prop, availability, impact strip, CTAs, voice.
- `layouts/partials/terminal.html` — CI/CD pipeline animation.
- `layouts/partials/head.html` — meta + expanded JSON-LD (Person/WebSite/BlogPosting/Breadcrumb).
- `layouts/robots.txt` (new) — AI crawlers + sitemap.
- `static/llms.txt` (new).
- `content/consulting.md` — FAQ section (+ FAQPage schema).
- `config.yaml` — params for location, availability, knowsAbout, education, email.

## Out of Scope (YAGNI)

- No JS framework / build pipeline changes. No new backend endpoints. No blog content rewrites
  beyond light voice edits. No light/dark redesign beyond token cleanup.

## Verification

- `hugo --minify` builds clean; spot-check generated `robots.txt`, `sitemap.xml`, JSON-LD in `public/`.
- Validate JSON-LD parses (jsonify). Manual visual check of home in dark + light.
- Reduced-motion: animation renders static final state.
