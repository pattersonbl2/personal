---
title: "Current Resume"
url: "/resume/"

---

{{< resume-download >}}

## BRANDON PATTERSON · PLATFORM ENGINEER | SRE

### Charlotte, NC | linkedin.com/in/pattersonbl2

Platform-focused DevOps Engineer with 5+ years of experience designing and operating Kubernetes-based infrastructure across AWS and GCP. Strong background in Infrastructure as Code (Terraform, Atlantis), observability at scale, and GitOps workflows. Proven track record reducing toil, improving developer velocity, and increasing reliability for shared platforms and production services. Passionate about platform engineering, scalable systems, and continuous improvement.

CORE COMPETENCIES

**Platforms & Cloud:** Kubernetes (GKE, K3s), Docker, Proxmox, GCP (Cloud Run, BigQuery, Cloud SQL), AWS, Helm  
**Infrastructure as Code & GitOps:** Terraform, Atlantis, ArgoCD, Helm, GitHub Actions, CircleCI, Jenkins, deployment automation, release engineering  
**Observability & Reliability:** Datadog (monitors, SLOs, on-call, RUM, dashboards), Prometheus, Grafana, Loki, Alertmanager, incident response, capacity planning  
**Networking & Security:** VXLAN, Kubernetes networking, BGP, NGINX, DNS, TLS, load balancing, network security, supply-chain CI  
**Programming Languages:** Python, Go, JavaScript, TypeScript, Bash  
**Leadership & Operations:** Cross-functional collaboration, infrastructure planning, FinOps, platform modernization, technical debt reduction

PROFESSIONAL EXPERIENCE

**DevOps / Platform Engineer |** Jasper.ai | Remote | 2026 – Present

- Overhauled Datadog observability for multiple engineering teams — consolidated SLOs, standardized monitor tagging and routing, and reduced alert noise by retiring hundreds of low-signal alerts. Delivered LLM-aware golden signals, RUM dashboards, team-based on-call schedules, and critical-only alerting.
- Eliminated Argo CD sync lag after merges by automating post-merge refresh across the releases pipeline, with supporting Helm (Argo accounts/RBAC) and Terraform (Workload Identity Federation + IAP-style auth for CI automation).
- Implemented Atlantis to enable pull request–driven Terraform workflows, improving visibility, enforcing approval processes, and reducing manual infrastructure operations.
- Designed and implemented GitHub organization automation to standardize repository configuration, permissions, and CI/CD workflows, improving developer onboarding and reducing manual setup.
- Drove Terraform changes for BigQuery–Cloud SQL connectivity, replica/DNS/Helm health-check cleanup, VPC Flow Logs enablement, and default VPC reduction for improved security posture.
- Owned urgent supply-chain remediation for a compromised Trivy GitHub Action — assessment, pinning, and coordinated remediation across Security and DevOps. Supported compliance-oriented work including customer data deletion and access pattern review.
- Executed FinOps initiatives (GCS lifecycle policies, cloud cleanup, legacy platform decommission) and cross-team platform requests (secrets, IAM, BigQuery, Cloud Run, vendor integrations). Authored internal audits, runbooks, and rollout plans to align engineering and leadership on observability and infrastructure changes.

**Site Reliability Engineer |** Mozilla | Remote | 2023 – 2025

- Migrated production workloads from AWS to GCP using Terraform, cutting cloud costs by 40% while improving platform clarity and maintainability.
- Built and launched new GCP infrastructure from scratch, reducing deployment time from hours to seconds through reusable Terraform patterns and improved Kubernetes platform onboarding.
- Operated multi-tenant Kubernetes clusters supporting 50M+ users, maintaining availability, performance, and platform consistency at scale.
- Developed Helm charts and GitHub Actions pipelines that standardized deployments and eliminated manual release effort across multiple services.
- Served as incident responder for critical services, driving root cause analysis and follow-up reliability improvements.
- Owned observability, automation, release engineering, capacity planning, and on-call operations across multiple production services via MozCloud platform triage rotation.

**Software Engineer, Hubs Support |** Mozilla | Remote | 2021 – 2023

- Owned and maintained GKE-based production infrastructure, improving operational efficiency across shared environments.
- Built and improved GKE infrastructure, enabling faster deployments and a better developer experience.
- Built internal QA tooling in TypeScript to automate test workflows and reduce repetitive engineering work.
- Automated testing processes that reduced time to completion by 75%, from one day to under two hours.
- Improved CI/CD pipelines with GitHub Actions, reducing deployment friction and supporting faster delivery.

**SysOps Administrator (Azure/AWS) |** Audacious Inquiry | 2020 – 2021

- Served as identity and access management SME for Azure AD, supporting secure access for 200+ users.
- Implemented endpoint protection services to maintain compliance with data privacy and security requirements.
- Deployed AWS infrastructure with Terraform and automated onboarding workflows, reducing onboarding time by 10%.
- Managed SSL certificates and DNS for state agency systems, maintaining uptime and compliance.

**IT Administrator |** Single Stone Consulting | Richmond, VA | 2019 – 2020

- Managed day-to-day IT administration including user provisioning, software lifecycle work, and troubleshooting for consulting teams.
- Supported hybrid cloud environments spanning office systems, financial operations, and AWS-hosted infrastructure.
- Reduced downtime and improved service continuity through practical cloud operations and infrastructure support.

### CONSULTANCY & SPECIAL PROJECTS

- **Personal Homelab Platform (2024 – Present)** Multi-node Kubernetes platform (Proxmox + K3s) with VXLAN-isolated networking, GitOps delivery via ArgoCD/Terraform/Ansible, Traefik + MetalLB/Cilium ingress, and a full Prometheus/Grafana/Loki/Alertmanager observability stack. Also spans Docker, VyOS, TrueNAS, and Oracle Cloud VPS.
- **ark31.info (2025 – Present)** Designed and deployed a personal blog and portfolio platform using Hugo with a custom theme and a Go backend API on GCP Cloud Run for contact submissions, resume delivery, rate limiting, and operational controls.
- **Email Services Deployment (2023 – 2024)** Designed and deployed scalable email infrastructure for relay.firefox.com and mozmail.com using Terraform and AWS.
- **Technical Support, DLA Windows 10 Migration (2018 – 2019)** Provided technical support during the DoD enterprise Windows 10 migration, deploying systems for 500 users across secure environments while minimizing downtime.

EDUCATION

**Virginia Commonwealth University, Richmond, VA:** _Bachelor of Science_
