---
title: "Closing the Gap: Eight Months From Laid Off to ArgoCD SME"
date: 2026-06-11T09:00:00-04:00
draft: false
summary: "In October 2025 I left Mozilla feeling behind on development skills. Eight months later I'm the ArgoCD SME at Jasper AI. Here's what actually closed the gap."
tags: ["career", "platform-engineering", "homelab"]
---

In October 2025 I left Mozilla, and I'll be honest about where my head was: the market was rough, and I knew the biggest hole in my skill set was active development. I had spent years operating infrastructure at scale without being asked to write much code, and I'd let that slide for too long. It felt like a huge mountain to climb.

Eight months later I'm a DevOps engineer at Jasper AI, the team's ArgoCD SME, and I've shipped a [free Kubernetes training platform](/posts/k8s-learn-try-it/) that strangers on the internet actually use. This post is about what happened in between — not because the story is special, but because the approach is repeatable.

## Step one: build something real, in public

Instead of grinding LeetCode, I built the platform I'd want to operate professionally. My homelab went from a simple Docker setup to a [multi-node K3s cluster](/posts/k3s-deep-dive/) with GitOps delivery through ArgoCD, Terraform and Ansible for provisioning, VXLAN-isolated networking, and a [full Prometheus/Grafana/Loki observability stack](/posts/homelab-monitoring-deep-dive/).

The point wasn't the hardware. The point was that every skill I felt shaky on — writing code, structuring Terraform, debugging Kubernetes from first principles — got exercised daily, with real consequences when I got it wrong.

## Step two: use AI as a multiplier, not a crutch

I leaned hard on AI assistance to accelerate the learning, and I [wrote about that honestly](/posts/ai-engineering-in-the-loop/) — including the time it helped me silently break my own platform for 24 hours. The lesson that stuck: AI raises the ceiling for engineers who stay in the loop and verify everything. It doesn't replace the judgment you build by operating systems yourself.

## Step three: the gap closes faster than you think

In January 2026 I joined Jasper AI. The homelab work transferred directly: within months I had taken platform ownership of our ArgoCD delivery pipeline, led a Datadog observability overhaul across multiple engineering teams, and handled an urgent supply-chain security remediation.

The development gap I was so worried about? Building [K8s Learn](https://learn.bp31app.com) — a Go-backed training platform with a browser terminal, running on my own cluster — answered that question for me.

## If you're in the October-2025 version of this story

Three things I'd tell anyone staring at the same mountain:

1. **Build the thing you want to be hired to run.** A homelab with real GitOps, real monitoring, and real incidents teaches more than any course.
2. **Write it down publicly.** Every post on this site doubled as interview material.
3. **The gap is smaller than it feels.** Eight months of deliberate, hands-on work changed my trajectory completely.

If you're working through something similar and want to talk platforms, GitOps, or career transitions — [get in touch](/contact/).
