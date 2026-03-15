# K8s Learn Blog Post Implementation Plan

> **For agentic workers:** REQUIRED: Use superpowers:subagent-driven-development (if subagents available) or superpowers:executing-plans to implement this plan. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Write a concise, punchy blog post reviewing the self-hosted K8s Learn training platform for ark31.info.

**Architecture:** Single Hugo markdown file following existing PaperMod conventions. Problem-solution pairs within architecture layers. No code changes — content only.

**Tech Stack:** Hugo, Markdown, PaperMod theme

**Spec:** `docs/superpowers/specs/2026-03-15-k8s-learn-blog-post-design.md`

---

## Chunk 1: Research & Write

### Task 1: Confirm cleanup mechanism

Before writing, verify how the 24-hour namespace cleanup works.

**Files:**
- Read: `/Users/bp31/Documents/gitops/infrastructure/k8s-learn/` (all files)
- Read: n8n workflows if accessible, or check for CronJob manifests in gitops repo

- [ ] **Step 1: Search gitops repo for cleanup mechanism**

Search for CronJob, TTL, or scheduled cleanup related to k8s-learn namespaces:
```bash
grep -r "cleanup\|cronjob\|ttl\|schedule" /Users/bp31/Documents/gitops/infrastructure/k8s-learn/
grep -ri "k8s-learn.*clean\|learn.*namespace.*delet" /Users/bp31/Documents/gitops/
```

- [ ] **Step 2: Read key k8s-learn manifests for technical accuracy**

Read these files to confirm details before writing:
```
/Users/bp31/Documents/gitops/infrastructure/k8s-learn/kustomization.yaml
/Users/bp31/Documents/gitops/infrastructure/k8s-learn/nginx-signup.conf
/Users/bp31/Documents/gitops/infrastructure/k8s-learn/nginx-tasks.conf
/Users/bp31/Documents/gitops/infrastructure/k8s-learn/n8n-sa.yaml
/Users/bp31/Documents/gitops/infrastructure/k8s-learn/tasks-page.yaml
/Users/bp31/Documents/gitops/infrastructure/k8s-learn/signup-page.yaml
```

- [ ] **Step 3: Note findings for use in writing**

Record cleanup mechanism (CronJob, n8n workflow, or manual). If unclear, write it as "scheduled cleanup" without specifying mechanism. Note any technical details (exact resource limits, rate limit values, etc.) for accuracy.

---

### Task 2: Write the blog post

**Files:**
- Create: `/Users/bp31/Documents/personal/content/posts/k8s-learn-review.md`
- Reference: `/Users/bp31/Documents/personal/docs/superpowers/specs/2026-03-15-k8s-learn-blog-post-design.md`
- Reference: `/Users/bp31/Documents/personal/content/posts/homelab-end-to-end-2026.md` (style guide)
- Reference: `/Users/bp31/Documents/gitops/infrastructure/k8s-learn/` (technical details)

- [ ] **Step 1: Write the frontmatter and opening**

```markdown
---
title: "K8s Learn: A Self-Hosted Kubernetes Training Platform"
date: 2026-03-15T12:00:00-05:00
draft: false
---

[Opening: 2-3 sentences. Lead with live link, frame as project you built, one-liner on stack.]
```

Match the voice of existing posts — first person, direct, no fluff.

- [ ] **Step 2: Write Section 1 — The Frontend (`##` heading)**

Problem-solution format. Mention: static HTML/CSS, nginx in K8s, 5 progressive tasks, ConfigMap-mounted via kustomize, lightweight resource requests. Keep to 1-2 short paragraphs. Use `##` for section heading (major sections), `###` only for subsections if needed.

- [ ] **Step 3: Write Section 2 — User Provisioning**

Problem-solution format. Mention: n8n running outside cluster on VM 110 with SA token auth, webhook-triggered namespace + RBAC creation, rate limiting, progress bar UX. Keep to 1-2 short paragraphs.

- [ ] **Step 4: Write Section 3 — Per-User Routing**

Problem-solution format. Mention: `*.bp31app.com` catch-all to Traefik, single-level subdomain pattern for Cloudflare free plan, user just runs `kubectl create ingress`. Keep to 1-2 short paragraphs.

- [ ] **Step 5: Write Section 4 — AI Help Widget**

Problem-solution format. Mention: floating button, n8n webhook to Ollama, qwen2.5:14b on RTX 3060 via PCI passthrough, scoped system prompt, rate limiting. Keep to 1-2 short paragraphs.

- [ ] **Step 6: Write Section 5 — Cleanup**

Problem-solution format. Mention: 24-hour lifecycle, mechanism found in Task 1. Shortest section — 2-3 sentences max.

- [ ] **Step 7: Write the closing**

2-3 sentences on what you'd change/add next. End with italicized repo link line (single repo since K8s Learn lives entirely in gitops):
```markdown
*Check out the repo: [gitops](https://github.com/pattersonbl2/gitops)*
```

- [ ] **Step 8: Review the full post**

Read the complete file. Check:
- Tone is concise and punchy (not tutorial-verbose)
- Each section is problem → solution
- No section exceeds ~150 words
- Horizontal rules (`---`) between major sections (matching existing posts)
- Code blocks have language hints where used
- Total post length: ~100-150 lines target

- [ ] **Step 9: Create branch and commit**

All git commands run from `/Users/bp31/Documents/personal`:
```bash
cd /Users/bp31/Documents/personal && git checkout -b blog/k8s-learn-review
cd /Users/bp31/Documents/personal && git add content/posts/k8s-learn-review.md && git commit -m "Add K8s Learn platform review blog post"
```

---

### Task 3: Local preview verification

- [ ] **Step 1: Build the Hugo site locally**

```bash
cd /Users/bp31/Documents/personal && hugo server -D
```

- [ ] **Step 2: Verify the post renders correctly**

Check `http://localhost:1313/posts/k8s-learn-review/` — confirm formatting, headers, code blocks, links all render as expected.

- [ ] **Step 3: Fix any rendering issues and re-commit if needed**

---

### Task 4: Push and open PR

- [ ] **Step 1: Push branch and create PR**

```bash
cd /Users/bp31/Documents/personal && git push -u origin blog/k8s-learn-review
cd /Users/bp31/Documents/personal && gh pr create --title "Add K8s Learn platform review blog post" --body "New blog post reviewing the self-hosted K8s training platform (learn.bp31app.com). Concise tech review covering architecture, provisioning, routing, AI help widget, and cleanup."
```
