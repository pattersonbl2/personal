---
title: "Consulting"
url: "/consulting/"
hideMeta: true
description: "Consulting services from Brandon Patterson — digital readiness, SRE / platform engineering, and web development & hosting."
faq:
  - q: "What does Brandon Patterson do?"
    a: "Brandon is a Platform Engineer and SRE with 6+ years building and operating Kubernetes infrastructure on GCP and AWS. He focuses on GitOps delivery, Terraform-based infrastructure as code, and observability at scale, currently at Jasper AI and previously at Mozilla."
  - q: "Is Brandon available for hire or consulting?"
    a: "Yes. Brandon is open to senior Platform/SRE roles and takes select consulting engagements — short-term projects, monthly retainers, or hourly advisory. He works remotely from Charlotte, NC across time zones."
  - q: "What kinds of consulting engagements do you take?"
    a: "Three shapes: fixed-scope short-term projects (migrations, build-outs, audits, observability rollouts), monthly retainers for ongoing platform support and on-call backup, and pay-as-you-go hourly advisory for architecture decisions or unblocking a specific question."
  - q: "What is your experience with Kubernetes?"
    a: "Brandon has operated multi-tenant Kubernetes clusters supporting 50M+ users, runs a multi-node Proxmox + K3s homelab with 19 ArgoCD-managed applications, and works across GKE, EKS, and K3s with Helm, Cilium, MetalLB, and VXLAN networking."
  - q: "How do I get in touch?"
    a: "Send a message through the contact form at ark31.info/contact, or email bpatterson@ark31.info. Brandon typically responds within a day or two."
---

I help small and mid-sized teams stand up production-grade infrastructure, ship faster with less toil, and put a credible presence on the web. Five years of hands-on platform engineering across AWS, GCP, and Kubernetes — most recently at Jasper.ai and Mozilla.

If any of the work below sounds like what you need, [send a message](/contact/) and I'll get back to you.

## Services

### Digital Readiness

A focused assessment of where you are today and a concrete plan for where you want to go. Best for teams modernizing legacy systems, planning a cloud migration, or trying to figure out which problems to solve first.

- Infrastructure and cloud-spend audit — what you're running, what it costs, what can be cut
- Migration planning across AWS, GCP, and on-prem (including Proxmox / K3s homelab-style setups)
- Security posture review — IAM, network segmentation, secret management, supply-chain CI
- Observability gap analysis — what you're not seeing today and what to wire up first
- Roadmap with prioritized, sized work items your team can actually execute

### SRE & Platform Engineering

Practical reliability and platform work, not slideware. I've operated multi-tenant Kubernetes clusters at 50M+ user scale and rebuilt observability stacks from the ground up.

- Kubernetes operations — GKE, EKS, K3s; Helm packaging; cluster upgrades; networking (Cilium, MetalLB, VXLAN)
- GitOps delivery with ArgoCD, Atlantis, and Terraform Workload Identity Federation
- Observability buildout — Datadog, Prometheus, Grafana, Loki, Alertmanager; SLOs, golden signals, on-call routing
- Incident response design and runbook authoring
- CI/CD pipelines on GitHub Actions, CircleCI, and Jenkins
- FinOps cleanup — lifecycle policies, idle-resource reclamation, right-sizing

### Web Development & Hosting

End-to-end builds for personal sites, portfolios, small-business sites, and lightweight SaaS landing pages. This site itself is a working reference — Hugo + Cloudflare Pages with a Go backend on GCP Cloud Run.

- Static sites with Hugo, custom themes, and responsive design
- Cloudflare Pages, GCP Cloud Run, or AWS hosting with auto-deploys from GitHub
- Go or Python backend APIs for contact forms, downloads, or light app functionality
- Custom domains, DNS, TLS, and email setup
- Ongoing hosting, monitoring, and patching if you want it handled for you

## How I work

I take on engagements in three shapes — pick whichever fits the problem.

- **Short-term project** — fixed scope, fixed timeline. Migrations, build-outs, audits, observability rollouts. Best when the deliverable is clear.
- **Retainer** — a set number of hours per month for ongoing platform support, on-call backup, or hosting and maintenance. Best for teams without a dedicated SRE.
- **Hourly advisory** — pay-as-you-go calls and async review for architecture decisions, hiring panels, or unblocking a specific question.

## Why work with me

- 5+ years in production SRE / DevOps roles at Mozilla and Jasper.ai
- Operated Kubernetes at 50M+ user scale; migrated production workloads AWS → GCP and cut cloud spend 40%
- Built and run a multi-node homelab on Proxmox / K3s with full GitOps and observability — the same patterns I apply to client work
- Comfortable working remotely across time zones and embedding with existing engineering teams

## Frequently asked

{{< faq >}}

## Get in touch

[Send a message →](/contact/)
