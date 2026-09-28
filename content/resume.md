---
title: "Resume"
url: "/resume/"
description: "Resume of Brandon Patterson — Platform / SRE engineer focused on production observability for AI products, Kubernetes, GitOps, and Terraform."

---

{{< resume-download >}}

## BRANDON PATTERSON · PLATFORM ENGINEER | SRE

### Charlotte, NC | bpatterson@ark31.info | ark31.info | linkedin.com/in/pattersonbl2

Platform / SRE engineer focused on production observability for AI products. 6+ years operating Kubernetes on AWS and GCP; currently building Datadog SLOs and LangSmith agent tracing at Jasper.ai so teams can debug latency, cost, and failure modes in live agent workflows. Strong GitOps (ArgoCD), Terraform, and incident response background from Mozilla-scale systems (50M+ users). Proven track record reducing toil, improving developer velocity, and increasing reliability for shared platforms and production services.

CORE COMPETENCIES

**Platforms & IaC:** Kubernetes (GKE, K3s), Docker, Terraform, ArgoCD, Helm, GitHub Actions, GCP (Cloud Run, BigQuery, Cloud SQL), AWS  
**Observability & AI:** Datadog (SLOs, on-call, RUM), Prometheus/Grafana/Loki, LangSmith agent tracing, LLM golden signals (latency, cost, quality)  
**Networking & Security:** VXLAN, Kubernetes networking, BGP, NGINX, DNS/TLS, supply-chain CI  
**Languages:** Python, Go, JavaScript/TypeScript, Bash

PROFESSIONAL EXPERIENCE

**DevOps / Platform Engineer |** Jasper.ai | Remote | 2026 – Present

- Building observability for Jasper's AI agent stack with LangSmith — production tracing across agent runs, surfacing latency/cost/error/quality signals, and connecting agent telemetry into Datadog SLOs and on-call.
- Overhauled Datadog observability across multiple engineering teams — consolidated SLOs, standardized monitor tagging/routing, and cut alert noise by retiring hundreds of low-signal alerts; shipped RUM dashboards and team-based on-call schedules.
- ArgoCD SME — own the GitOps delivery pipeline; resolved repo-server timeouts and CMP stability issues causing production sync failures, and improved reliability via automated post-merge sync refresh, Helm-based RBAC, and Terraform Workload Identity Federation.
- Implemented Atlantis for PR-driven Terraform workflows, improving visibility and enforcing approval processes.
- Led urgent supply-chain remediation for a compromised Trivy GitHub Action, coordinating assessment and fix across Security and DevOps.

**Site Reliability Engineer |** Mozilla | Remote | 2023 – 2025

- Migrated production workloads from AWS to GCP using Terraform, cutting cloud costs 40% while improving platform maintainability.
- Built new GCP infrastructure from scratch, cutting deployment time from hours to seconds via reusable Terraform patterns.
- Operated multi-tenant Kubernetes clusters supporting 50M+ users, maintaining availability and platform consistency at scale.
- Built Helm charts and GitHub Actions pipelines standardizing deployments; served as incident responder driving root-cause analysis and reliability follow-through.

**Software Engineer, Hubs Support |** Mozilla | Remote | 2021 – 2023

- Owned GKE-based production infrastructure, improving operational efficiency across shared environments.
- Built internal QA tooling in TypeScript, cutting test cycle time 75% (one day to under two hours).
- Improved CI/CD pipelines with GitHub Actions, reducing deployment friction and supporting faster delivery.

**SysOps Administrator (Azure/AWS) |** Audacious Inquiry | 2020 – 2021

- Served as IAM SME for Azure AD, supporting secure access for 200+ users.
- Deployed AWS infrastructure with Terraform and automated onboarding workflows, cutting onboarding time 10%.

**Earlier:** IT Administrator, Single Stone Consulting (2019–2020) — hybrid cloud IT operations for consulting teams. Technical Support, DoD Windows 10 Migration (2018–2019) — deployed systems for 500 users.

### CONSULTANCY & SPECIAL PROJECTS

- **Personal Homelab Platform & K8s Learn (2024 – Present)** Multi-node Kubernetes platform (Proxmox/K3s) with GitOps (ArgoCD/Terraform), full Prometheus/Grafana/Loki observability, plus a self-hosted K8s training platform with a live cluster and browser terminal.
- **ark31.info (2025 – Present)** Personal site and portfolio: Hugo with a custom theme and a Go backend on GCP Cloud Run.

EDUCATION

**Virginia Commonwealth University, Richmond, VA:** _Bachelor of Science, Sociology_
