---
title: "Homelab End-to-End: A Production-Grade Platform in 2026"
date: 2026-01-28T14:00:00-05:00
draft: false
---

A lot has changed since I first wrote about building my homelab with AI assistance. What started as a learning project has evolved into a production-grade platform spanning two repositories, multiple Terraform projects, a full GitOps pipeline, and services running on both Docker and Kubernetes. This post is both a progress update and a visual tour of where things stand today.

## The Big Picture

My homelab now runs across:

- **Proxmox cluster** — 2 physical nodes hosting all VMs
- **Docker hosts** — 2 VMs running containerized services (media, networking, Gitea, etc.)
- **K3s cluster** — 2 control plane nodes + workers for Kubernetes workloads
- **VyOS router** — Software-defined networking with BGP
- **TrueNAS** — Network storage
- **Oracle Cloud VPS** — Pangolin tunnel endpoint for external access

Everything is defined in code across two repos:

| Repo | Purpose |
|------|---------|
| [homelab](https://github.com/pattersonbl2/homelab) | Terraform (VM provisioning), Ansible (config management), Docker Compose stacks |
| [gitops](https://github.com/pattersonbl2/gitops) | ArgoCD applications, Kubernetes manifests, Helm values |

---

## Infrastructure Layer

### Terraform Projects

I manage five separate Terraform projects, each with its own state:

```
infrastructure/
├── terraform/           # Docker host VMs
├── terraform-k3s/       # K3s cluster VMs
├── terraform-vyos/      # VyOS router VM
├── terraform-truenas/   # TrueNAS VM
├── terraform-authentik/ # Authentik identity provider
└── terraform-grafana/   # Grafana Cloud dashboards
```

The reusable `proxmox-vm` module handles common patterns: cloud-init, networking, SSH keys. Each project just defines the specifics.

### Ansible Automation

Ansible handles everything post-provisioning:

- **Docker installation** and configuration
- **Stack deployment** — copies Compose files, manages `.env`, runs `docker compose up -d`
- **K3s cluster upgrades** — rolling upgrades across control plane and workers
- **Tailscale setup** — for remote access to the homelab network

The key playbook is `deploy-stacks.yml`, which takes a list of stacks and deploys them idempotently. Terraform calls this automatically after VM creation.

### Local CI/CD

One of the most useful additions: **git hooks that auto-deploy on commit**. When I commit changes to a stack, the hook detects what changed and runs the appropriate Ansible deployment.

```bash
git commit -m "Update traefik config"
# → Automatically detects and deploys changed stacks
```

No manual commands, no forgetting to deploy. Just commit and push.

---

## Docker Stacks

The homelab repo defines 10+ Docker Compose stacks:

| Stack | Services |
|-------|----------|
| **core** | Traefik (reverse proxy), Portainer, node-exporter, cAdvisor, pve-exporter |
| **media** | Jellyfin, Sonarr, Radarr, Prowlarr, qBittorrent, Gluetun VPN |
| **database** | PostgreSQL |
| **networking** | Pi-hole (DNS), Netbox (IPAM) |
| **gitea** | Gitea server, container registry, CI runner |
| **ai-tools** | Open WebUI + Ollama (local LLM) |
| **projects** | Vikunja (project & task tracking) |
| **newt** | Pangolin tunnel connector |
| **pangolin** | VPS-side tunnel control plane |
| **monitoring-vps** | Prometheus, Grafana, Loki on VPS |
| **authentik-vps** | Authentik identity provider on VPS |
| **crowdsec-vps** | CrowdSec security on VPS |
| **pihole-vps** | Pi-hole on VPS for external DNS |

All stacks share a common `env` file and follow the same patterns: Traefik labels for routing, consistent logging, health checks.

---

## Kubernetes (GitOps)

The K3s cluster runs workloads managed entirely through ArgoCD. The gitops repo contains:

### ArgoCD Applications

```
apps/
├── bootstrap.yaml         # Root app that deploys everything else
├── argocd.yaml            # ArgoCD self-management
├── cilium.yaml            # CNI with BGP
├── traefik.yaml           # Ingress controller
├── cloudflare-tunnel.yaml # Secure external access
├── monitoring.yaml        # Prometheus + Grafana
├── loki.yaml              # Log aggregation
├── longhorn.yaml          # Distributed storage
├── uptime-kuma.yaml       # Status monitoring
├── ntfy.yaml              # Push notifications
├── vaultwarden.yaml       # Password manager
├── gitea-runner.yaml      # CI runner for Gitea
└── newt.yaml              # Pangolin tunnel client
```

### Key Components

- **Cilium CNI** — Replaced Flannel; provides BGP peering with VyOS for LoadBalancer IPs
- **Cloudflare Tunnel** — Subnet routing (not hostname routing) for scalable external access
- **Tailscale + Cloudflare Access** — Access restricted to Tailscale IPs only
- **Full observability** — Prometheus scrapes Docker hosts, K3s, and Proxmox; Loki aggregates logs; Grafana dashboards

### kubectl Anywhere

I can access the cluster from anywhere using three kubeconfig options:

1. **Direct** — On the home network, connect directly to control plane
2. **SSH tunnel** — Port-forward through VyOS when on home network but different subnet
3. **Cloudflare tunnel** — Works from anywhere with Tailscale connected

---

## Networking

### VyOS Router

The VyOS VM handles:

- **Inter-VLAN routing** between lab network (10.10.10.0/24) and home network
- **BGP peering** with Cilium for Kubernetes LoadBalancer IPs
- **Static routes** for Kubernetes pod and service networks

### Cloudflare Tunnel + Tailscale

External access uses **subnet routing** instead of individual hostname routes:

- Routes `10.42.0.0/16` (pod network) and `10.43.0.0/16` (service network) through the tunnel
- Cloudflare Access policy restricts access to Tailscale IPs (`100.64.0.0/10`)
- No need to configure individual services — just add them to the cluster

This pattern scales without tunnel configuration changes.

### Pangolin Tunnel

For services that need public access (without Tailscale), I run **Pangolin** — a self-hosted tunnel solution:

- Control plane runs on Oracle Cloud VPS
- `newt` connector runs on Docker host and K3s
- Provides stable public endpoints for specific services

---

## Observability

### Monitoring Stack

- **Prometheus** — Scrapes metrics from Docker hosts, K3s, Proxmox, and individual services
- **Grafana** — Dashboards for everything; some managed via Terraform
- **Loki** — Log aggregation from Docker and Kubernetes
- **Alertmanager** — Routes alerts to ntfy for push notifications
- **Uptime Kuma** — External health checks and status page

### Exporters

The `core` stack runs exporters for comprehensive coverage:

- `node-exporter` — Host-level metrics
- `cAdvisor` — Container metrics
- `pve-exporter` — Proxmox metrics

---

## What's Changed Since Last Time

Since my previous homelab post, I've added:

1. **VPS integration** — Pangolin, CrowdSec, Authentik, monitoring all running on Oracle Cloud
2. **Full monitoring pipeline** — Prometheus + Loki + Grafana with alerting to ntfy
3. **Longhorn storage** — Distributed block storage for K3s stateful workloads
4. **Uptime Kuma + Vaultwarden + Ntfy** — Self-hosted alternatives to SaaS tools
5. **Terraform for Authentik and Grafana** — Managing identity and dashboards as code
6. **K3s upgrade automation** — Ansible playbook for rolling cluster upgrades
7. **Local CI/CD with git hooks** — Auto-deploy on commit

The platform has gone from "learning project" to "I actually rely on this daily."

---

## Lessons Learned

1. **GitOps is worth the setup cost.** Once ArgoCD is running, adding services is just a YAML file and a commit.

2. **Subnet routing > hostname routing.** Configuring individual hostnames doesn't scale. Route the subnet and use DNS.

3. **Observability from day one.** I can debug issues in minutes because I have metrics, logs, and traces.

4. **Terraform state management matters.** Multiple smaller projects with isolated state beats one giant project.

5. **Document as you go.** The `CLAUDE.md` files and READMEs save me constantly when I come back to something after weeks.

6. **VPS extends homelab reach.** A free Oracle Cloud instance gives me public endpoints and redundancy without exposing my home IP.

---

## What's Next

- **Authentik SSO** — Single sign-on across all services
- **Backup automation** — TrueNAS snapshots + off-site replication
- **More K3s workloads** — Migrate remaining Docker services that benefit from HA
- **Cost tracking** — Power monitoring and cloud cost dashboards

---

If you're building a homelab, I'd encourage you to treat it like a real platform: version control everything, automate deployments, and invest in observability early. The compound benefits are real — every new service is easier to add, debug, and maintain.

*Check out the repos: [homelab](https://github.com/pattersonbl2/homelab) | [gitops](https://github.com/pattersonbl2/gitops)*
