---
title: "Deep Dive: How I Run Kubernetes in My Homelab"
date: 2026-03-25T20:00:00-04:00
draft: false
---

I've been running a K3s cluster in my homelab for a few months now. What started as "let me learn Kubernetes properly" turned into 19 ArgoCD-managed applications across 8 nodes, with real monitoring, network policies, and a GitOps workflow I actually trust. Here's how it all fits together.

---

## The Cluster

Two control plane nodes and six workers, all running as Proxmox VMs spread across two physical hosts. The control plane nodes are lightweight — 2 CPU, 4GB RAM each. Workers get 4 CPU and 8GB RAM. Everything sits on a dedicated lab network (10.10.10.0/24) behind a VyOS router, separate from my home network.

Terraform provisions the VMs with static IPs and cloud-init. K3s itself is installed manually — I wanted to understand the bootstrap process rather than abstracting it away. The kubeconfig lives at `~/.kube/homelab-k3s.yaml` and I reach the cluster through a WireGuard tunnel from my Mac.

```
Control Plane:
  k3s-cp-1 (10.10.10.11) — pve
  k3s-cp-2 (10.10.10.12) — pve2

Workers:
  k3s-worker-1 through k3s-worker-6
  Distributed across both Proxmox nodes
```

---

## GitOps with ArgoCD

Everything in the cluster is managed by ArgoCD using an app-of-apps pattern. A single bootstrap application watches the `apps/` directory in my gitops repo. Drop a new YAML file in there, push, and ArgoCD picks it up automatically.

Every application has `selfHeal: true` and `prune: true` enabled. If I manually change something in the cluster, ArgoCD reverts it. If I delete a manifest from git, ArgoCD deletes the resource. The cluster always matches what's in the repo.

Sync waves control the order: Cilium deploys first (wave -1) because nothing works without a CNI, then sealed-secrets (wave 0), then Traefik and the tunnel (wave 1), then everything else. This matters more than you'd think — I've had deployments fail because they tried to create ingress resources before Traefik existed.

ArgoCD manages itself too. Its own Helm values live in the gitops repo, and the ArgoCD application points at them. It's turtles all the way down.

---

## Networking

This is where the homelab gets interesting.

**Cilium** handles CNI duties. I replaced the default K3s networking with Cilium for its network policy engine, kube-proxy replacement, and Hubble observability. Every namespace has a default-deny ingress policy, with explicit rules for what's allowed. The monitoring namespace can scrape anything, Traefik can reach services it needs to route to, and everything else is locked down.

**Traefik** runs as the ingress controller on a LoadBalancer service. Cilium assigns it a floating IP on the lab network.

**Cloudflare Tunnel** connects the cluster to the internet without opening any ports. A `cloudflared` pod maintains an outbound tunnel to Cloudflare's edge, and a wildcard DNS rule routes `*.bp31app.com` through the tunnel to Traefik. When a user creates an ingress resource in the cluster, it's publicly accessible within seconds. One gotcha: QUIC breaks WebSocket connections through the tunnel, so I run it in http2 mode with chunked encoding disabled.

**Pangolin** provides a second path in through a VPS. The `newt` client runs in-cluster and connects to a Pangolin control plane on an Oracle Cloud VPS. This gives me `*.vps.bp31app.com` URLs that bypass Cloudflare entirely — useful for admin tools like Longhorn and Prometheus that I don't want on the public tunnel.

**WireGuard** is how I access the cluster from my laptop. A tunnel from my Mac to the VyOS router puts me on the lab network, so kubectl commands work directly against the API server.

---

## Monitoring

The monitoring stack is probably the part I've spent the most time on.

**Prometheus** scrapes metrics from everywhere — not just the K3s cluster, but also my Docker host, the VyOS router (via SNMP), Proxmox nodes, an RTX 3060 GPU, the VPS, and even my media stack. Retention is 15 days at 10GB. The scrape configs are a mix of in-cluster ServiceMonitors and static external targets.

**Grafana** has dashboards managed entirely through Terraform. I learned the hard way not to use Helm's `gnetId` imports — they break when upstream dashboards change. Now every dashboard is a JSON file in the `terraform-grafana` directory, applied via the Grafana Terraform provider through a port-forward. Grafana auth goes through Authentik OAuth so I get SSO across my homelab.

**Loki and Promtail** collect pod logs from every node with 30-day retention. Having logs alongside metrics in the same Grafana instance makes debugging significantly faster.

**Alertmanager** routes alerts to Ntfy for push notifications on my phone and Slack for less urgent things. Custom rules watch for node downtime, disk pressure, memory spikes, and pod crash loops. I disabled the default kubelet recording rules from kube-prometheus-stack because K3s doesn't expose the metrics they expect, and they were generating constant false-positive alerts.

**Hubble** gives me network flow visibility at the Cilium level. I can see which pods are talking to each other, which DNS queries are happening, and which connections are being dropped by network policies. The UI is exposed behind Pangolin auth.

**Uptime Kuma** monitors 16 services across both Cloudflare and Pangolin tunnel paths. Each check validates the full routing chain — DNS, tunnel, Traefik, service — so if any link breaks, I know about it. Alerts go to Ntfy.

---

## Security

Security in a homelab is mostly about building good habits for production.

**Network policies** are the foundation. Default-deny ingress on every namespace means I have to explicitly allow traffic. This has caught real issues — like when I added Ntfy notifications to Uptime Kuma and they silently failed because I hadn't allowed cross-namespace traffic yet.

**Sealed Secrets** handle sensitive data in git. I encrypt secrets locally with `kubeseal`, commit the ciphertext, and the controller decrypts them in-cluster. This replaced a KSOPS setup that was the most fragile part of the stack — GPG keyrings, kustomize plugin directories, ArgoCD sidecars. Sealed Secrets is one Helm chart and it just works.

**RBAC** is scoped tightly. The n8n provisioner that creates K8s Learn environments has a ClusterRole limited to namespace operations, with per-namespace Roles for the resources it actually needs. ArgoCD uses SAML SSO through Google Workspace. Grafana authenticates through Authentik.

**CrowdSec** runs on the VPS as an intrusion detection layer for anything hitting the Pangolin tunnel.

---

## Storage

Longhorn handles persistent storage with 2 replicas across nodes. Every stateful service — Prometheus, Grafana, Loki, Vaultwarden, Uptime Kuma, Ntfy, Ollama — gets a Longhorn PVC. The Longhorn UI is exposed behind Pangolin for volume management and snapshot operations.

It's not the fastest storage solution, but it's resilient. I've drained nodes for maintenance and the volumes reattach on other nodes without data loss.

---

## What's Actually Running

Here's what the cluster does day to day:

| Service | What it does |
|---------|-------------|
| **K8s Learn** | Interactive Kubernetes training platform with per-user namespaces and browser terminals |
| **Vaultwarden** | Self-hosted password manager (Bitwarden-compatible) |
| **Uptime Kuma** | Status monitoring for all 16 public services |
| **Ntfy** | Push notifications for alerts and events |
| **Ollama** | Local LLM inference (qwen2.5:14b on RTX 3060) |
| **Gitea Runner** | CI runner for the Gitea instance on my Docker host |
| **Newt** | Pangolin tunnel client with OpenTelemetry metrics |

Plus the infrastructure services: ArgoCD, Traefik, Cloudflare Tunnel, Prometheus, Grafana, Loki, Alertmanager, Cilium, Longhorn, and Sealed Secrets.

---

## Lessons Learned

**GitOps is worth the setup cost.** The first few days of converting everything to ArgoCD were painful. But now I never SSH into a node to fix something. Every change goes through git, gets reviewed in a PR, and syncs automatically. When something breaks, `git log` tells me exactly what changed.

**Network policies catch real bugs.** They're annoying to configure and easy to skip. Don't skip them. The first time a default-deny policy surfaces a missing connection you didn't know about, you'll understand why production clusters require them.

**Monitor everything, alert selectively.** I have dashboards for things I'll never look at proactively — and that's fine. But alerts need to be actionable. Disabling the noisy kubelet recording rules and only alerting on things I can actually fix reduced alert fatigue to near zero.

**Keep it simple until it breaks.** I started with `latest` image tags and no resource limits. That worked until it didn't. Now images are pinned, every pod has requests and limits, and PVCs have explicit storage classes. But I only added each constraint when I hit a real problem, not preemptively.

---

*Both repos are public if you want to dig into the details:*
- *[gitops](https://github.com/pattersonbl2/gitops) — ArgoCD apps and K8s manifests*
- *[homelab](https://github.com/pattersonbl2/homelab) — Terraform, Ansible, Docker stacks*
