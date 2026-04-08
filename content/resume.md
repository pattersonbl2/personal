---
title: "Current Resume"
url: "/resume/"

---

{{< resume-download >}}

## BRANDON PATTERSON · PLATFORM ENGINEER | SRE

### Charlotte, NC | linkedin.com/in/pattersonbl2

Platform and reliability engineer with experience building and operating Kubernetes platforms, cloud infrastructure, CI/CD systems, and internal developer workflows across GCP and AWS. Strong background in Terraform, GitOps, observability, and production operations at scale. Proven track record reducing toil, improving developer velocity, modernizing infrastructure, and increasing reliability for shared platforms and production services. Passionate about platform engineering, scalable systems, and continuous improvement.

CORE COMPETENCIES

**Platform & Cloud Infrastructure:** Kubernetes, Terraform, GCP, AWS, Docker, Helm, Proxmox, Cloud Run  
**CI/CD, GitOps & Developer Enablement:** GitHub Actions, CircleCI, Jenkins, ArgoCD, Helm, deployment automation, release engineering  
**Observability & Reliability:** Datadog (monitors, SLOs, on-call, RUM, dashboards), Prometheus, Grafana, Loki, Alertmanager, incident response, capacity planning  
**Networking & Security:** Kubernetes networking, cloud networking, BGP, NGINX, DNS, TLS, network security, supply-chain CI  
**Programming Languages:** Python, Go, JavaScript, TypeScript, Bash  
**Leadership & Operations:** Cross-functional collaboration, infrastructure planning, FinOps, platform modernization, technical debt reduction

PROFESSIONAL EXPERIENCE

**Mid-Level DevOps Engineer |** Jasper.ai | Remote | 2026 – Present

- Led a Datadog observability overhaul: team-based on-call (schedules, escalations, Slack routing), runbooks, SLO consolidation, and standardization of hundreds of monitors with team tags, routing, and noise reduction. Delivered LLM-aware golden signals, RUM dashboards, and critical-only alerting.
- Cut Argo CD sync lag after merges by automating post-merge refresh across the releases pipeline, with supporting Helm (Argo accounts/RBAC) and Terraform (Workload Identity Federation + IAP-style auth for CI automation).
- Drove Terraform changes for BigQuery–Cloud SQL connectivity, replica/DNS/Helm health-check cleanup, VPC Flow Logs enablement, and default VPC reduction for improved security posture.
- Owned urgent supply-chain remediation for a compromised Trivy GitHub Action — assessment, pinning, and coordinated remediation across Security and DevOps. Supported compliance-oriented work including customer data deletion and access pattern review.
- Executed FinOps initiatives (GCS lifecycle policies, cloud cleanup, legacy platform decommission) and cross-team platform requests (secrets, IAM, BigQuery, Cloud Run, vendor integrations). Authored internal audits, runbooks, and rollout plans to align engineering and leadership on observability and infrastructure changes.

**Site Reliability Engineer |** Mozilla | Remote | 2023 – 2025

- Designed and launched new GCP infrastructure with Terraform, improving scalability and reducing manual operational overhead.
- Built reusable infrastructure patterns for shared environments, improved onboarding to the shared Kubernetes platform, and reduced deployment time from hours to seconds.
- Migrated production workloads from AWS to GCP, reducing cloud costs by 40% while improving platform clarity and maintainability.
- Operated multi-tenant Kubernetes clusters supporting 50M+ users with a focus on availability, performance, and platform consistency.
- Developed Helm charts and GitHub Actions pipelines to standardize deployments, improve developer workflows, and eliminate significant manual release effort.
- Served as incident responder for critical services, driving root cause analysis and follow-up reliability improvements.
- Participated in MozCloud platform triage rotation, supporting production services and shared infrastructure.
- Led initiatives across observability, automation, release engineering, capacity planning, and on-call operations for multiple production services.

**Software Engineer, Hubs Support |** Mozilla | Remote | 2021 – 2023

- Supported GKE-based production infrastructure and improved operational efficiency across shared environments.
- Helped build and improve GKE infrastructure, enabling faster deployments and a better developer experience.
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

- **Personal Homelab Platform (2024 – Present)** Built and operate a production-grade homelab platform spanning Proxmox, K3s, Docker, VyOS, TrueNAS, and Oracle Cloud VPS. Managed with GitOps via ArgoCD, Terraform for infrastructure, Ansible for configuration, and a full observability stack using Prometheus, Grafana, Loki, and Alertmanager.
- **ark31.info (2025 – Present)** Designed and deployed a personal blog and portfolio platform using Hugo with a custom theme and a Go backend API on GCP Cloud Run for contact submissions, resume delivery, rate limiting, and operational controls.
- **GCP Migration (2024 – 2025)** Led migration from single-tenant environments to shared Kubernetes clusters, improving resource efficiency, platform consistency, and operational clarity.
- **Email Services Deployment (2023 – 2024)** Designed and deployed scalable email infrastructure for relay.firefox.com and mozmail.com using Terraform and AWS.
- **Technical Support, DLA Windows 10 Migration (2018 – 2019)** Provided technical support during the DoD enterprise Windows 10 migration, deploying systems for 500 users across secure environments while minimizing downtime.

EDUCATION

**Virginia Commonwealth University, Richmond, VA:** _Bachelor of Science_

Certifications

**AWS Certified Solutions Architect - Associate:** 2020 – 2023 _(Expired)_
