---
title: "The Cheapest Production Stack I Could Build: $3.33 a Month, Plus $70 of Electricity"
date: 2026-11-15T09:00:00-05:00
draft: false
summary: "I run a two-node Proxmox cluster, 19 ArgoCD applications, a public site, a Go API and two local LLMs. The services cost $3.33 a month. The power costs $70. Here is every line."
tags: ["homelab", "kubernetes", "sre"]
---

I went to write this post thinking my platform cost under $20 a month. I was both right and badly wrong. Every service I pay a company for comes to **$3.33 a month**. The electricity to run it comes to **$70**. I had been quoting the small number for a year without noticing the big one sitting in the wall socket.

Here is the whole bill, line by line.

## The services bill is two domain names

That is not a simplification. It is the entire list of things that send me an invoice.

| Line | What it runs | Cost |
| --- | --- | --- |
| `ark31.info` + `bp31app.com` | this site, homelab services | **$40/year** — $3.33/month |
| Cloudflare Pages | the whole static site, build and CDN | $0 — free tier |
| Cloudflare Tunnel | public ingress with no open ports at home | $0 — free tier |
| Cloudflare Turnstile | contact-form anti-spam | $0 — free tier |
| Oracle Cloud VPS | Pangolin, CrowdSec, Authentik, Pi-hole, a second Grafana | $0 — free instance |
| GCP Cloud Run | the Go contact/resume API | $0 — still inside the free tier |
| GitHub | repo, CI, auto-merge | $0 — free for public repos |
| Proxmox, K3s, ArgoCD, Terraform, Prometheus, Grafana, Loki, Longhorn, VyOS, n8n, Ollama | the actual platform | $0 — all open source |
| openjev + Ollama | both [sre-reflex](/posts/shadow-mode-zero-labels/) triage models | $0 — local inference |

**$40 a year.** Two domains. Nothing else on that list has ever charged me a cent.

The Cloud Run line is the one I expected to hurt and does not. It is a Go binary that answers a contact form and a signed-token resume endpoint. It scales to zero, so it bills for the seconds it is awake, and the free tier has never run out. I put the API on Cloud Run for exactly one reason: I did not want my home IP answering POST requests from the open internet. It turns out that decision was free.

## Then there is the part that is not free

| Line | Cost per month |
| --- | --- |
| Electricity for the homelab | **$70.00** |
| Hardware, $1,000 amortised over 5 years | **$16.67** |
| Services, from the table above | $3.33 |
| **Total** | **$90.00** |

The five-year amortisation is an accounting convenience, not a plan. I did not buy any of this gear with a payback period in mind — I bought it for fun, which is the honest answer and the one that makes the hardware line arbitrary. Over three years it is $27.78 and the total is $101.11. Pick whichever horizon you believe.

What the table says is blunt: **96% of what this platform costs me is power and hardware.** Two physical Proxmox nodes, eight K3s VMs, a VyOS router, an RTX 3060 doing local inference, and the disks under Longhorn's two-replica storage. That is 24 hours a day of draw whether or not anyone visits the site.

Seventy dollars a month is $840 a year. I have been describing this setup as "basically free" in conversation. It is not. It is a hobby with a utility bill.

## What is actually worth paying for

Three things, in order.

**Electricity, without question.** The whole value of this platform is that it is mine, physically, and that I can break it. [I run eight nodes](/posts/k3s-deep-dive/) so I can drain one and watch Longhorn reattach the volumes. You cannot learn that from a managed control plane that never lets you see the failure. The $70 buys the failures.

**The domains.** $40 a year to own the namespace my work lives at. This site has got me further professionally than anything else I have built, and it costs less than one dinner.

**The GPU, though not for the reason I bought it.** The RTX 3060 was going to be for experiments. It now runs both halves of a triage bot and serves local inference at $0 marginal cost per query. That is the line that turned out to pay for itself.

And one thing not worth paying for: managed Kubernetes, *for this purpose*. If this were a product with customers I would buy EKS the same afternoon. As a place to learn, the managed version removes the exact layer I am trying to understand.

## The cost that does not appear on any bill

The honest total is not $90, because the largest input is weekends.

[My Kubernetes platform was silently broken for 24 hours](/posts/ai-engineering-in-the-loop/) and every signup went into a black hole. I built a triage bot to catch silent failures, and then [spent twelve days with the evaluation harness silently broken](/posts/shadow-mode-zero-labels/). Nobody invoices me for those hours, so they never show up when I say "under $20 a month." They are the real price, and they are also the thing I would pay again.

So: $90 a month, of which $3.33 goes to anyone who sends an invoice. A platform built almost entirely out of open source and other people's free tiers, running on hardware in my house and power from my meter.

The lesson is one I should have learned from monitoring, not from a spreadsheet. **I had been quoting the cost I could see.** The $20 figure was the one on a credit card statement, so it was the one I believed. The $70 was on a utility bill I never connected to the rack, and the hardware was a sunk cost I had already stopped counting. I measured the thing that was easy to instrument and called it the total — which is exactly the mistake I keep writing posts about.
