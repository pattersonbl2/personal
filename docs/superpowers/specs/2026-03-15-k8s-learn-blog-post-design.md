# K8s Learn Blog Post — Design Spec

## Overview

A concise showcase blog post for ark31.info reviewing the self-hosted K8s Learn training platform. Hybrid structure: architecture layers as skeleton, problem-solution pairs as the content within each layer.

## Target Audience

- Homelabbers/DevOps engineers who might want to build something similar
- Hiring managers evaluating engineering skills (portfolio piece)

## Tone

Concise, punchy. Short paragraphs. Highlight the most impressive engineering bits without going tutorial-deep.

## Post Structure

### Title

"K8s Learn: A Self-Hosted Kubernetes Training Platform"

### Opening (2-3 sentences)

What it is, who it's for, link to `learn.bp31app.com`. One-liner on the stack: K3s, ArgoCD, n8n, Ollama, Cloudflare Tunnel.

### Section 1: The Frontend

- **Problem:** Needed a guided learning experience, not just a shell.
- **Solution:** Static HTML/CSS served by nginx in K8s. 5 progressive tasks (Pod → Deployment → Service → Ingress → Scale). No framework — ConfigMap-mounted HTML via kustomize configMapGenerator. Lightweight: nginx:alpine containers at 10m CPU/16Mi RAM.

### Section 2: User Provisioning

- **Problem:** Each user needs an isolated namespace with RBAC, created on-demand from a signup form.
- **Solution:** n8n workflow triggered by signup webhook. Creates namespace, service account, scoped RBAC — all via K8s API using a dedicated ServiceAccount with ClusterRole. Rate-limited at nginx layer (2r/m per IP).

### Section 3: Per-User Routing

- **Problem:** Users deploy apps in their tasks and need real public URLs (`learn-[name]-app.bp31app.com`).
- **Solution:** Cloudflare Tunnel wildcard rule for `*.learn.bp31app.com` + Traefik IngressRoute. User runs `kubectl create ingress` with their subdomain and it just works. This required an explicit tunnel rule since Cloudflare only supports one level of subdomain on free plans.

### Section 4: AI Help Widget

- **Problem:** Users get stuck, but this is fully self-hosted — no SaaS LLM dependency.
- **Solution:** Floating "?" button → n8n webhook → Ollama running qwen2.5:14b on a dedicated GPU node (RTX 3060 via PCI passthrough). Scoped system prompt keeps answers relevant to the 5 tasks only. Rate-limited at 6r/m per IP with 45s response timeout.

### Section 5: Cleanup

- **Problem:** Can't let abandoned user namespaces accumulate indefinitely.
- **Solution:** 24-hour auto-cleanup lifecycle for provisioned environments.

### Closing (2-3 sentences)

What you'd change or add next (e.g., more tasks, better observability per user, Helm chart for the whole platform). Link to gitops repo if public.

## Frontmatter

```yaml
---
title: "K8s Learn: A Self-Hosted Kubernetes Training Platform"
date: 2026-03-15T12:00:00-05:00
draft: false
---
```

## Conventions

- Follow existing post style: narrative-driven, markdown headers (##/###), code blocks with language hints where useful
- Minimal frontmatter (title, date, draft) per existing posts
- No tags/categories (not used in existing posts)
- File name: `k8s-learn-review.md` in `content/posts/`
