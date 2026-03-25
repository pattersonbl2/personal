---
title: "How I Monitor Everything in My Homelab"
date: 2026-03-25
draft: false
---

I can tell you the CPU temperature of my Proxmox nodes, how many packets my VyOS router
dropped in the last hour, whether any Longhorn volumes are degraded, and if someone's
brute-forcing my VPS — all from the same Grafana instance. Here's how I built a monitoring
stack that covers every layer of my homelab, from bare metal to Kubernetes pods.

---

## The Stack

Everything funnels into three systems running in my K3s cluster:

- **Prometheus** — metrics from 20+ scrape targets, 15-day retention
- **Loki** — logs from every K8s pod and the VPS
- **Grafana** — 13 dashboards managed through Terraform, with Authentik SSO

Alertmanager routes critical alerts to Ntfy on my phone and Slack for everything else.
Uptime Kuma watches 16 services from the outside and pings me if any go down.

---

## What Gets Scraped

This is the part I'm most proud of. Prometheus doesn't just watch the K3s cluster — it
reaches into every corner of the infrastructure.

**Kubernetes internals** are covered by kube-prometheus-stack out of the box: node-exporter
on every K3s node, kube-state-metrics for pod/deployment status, and kubelet metrics.
ServiceMonitors auto-discover Traefik, Cilium Hubble, and Longhorn so I didn't have to
write scrape configs for those.

**The Docker host** (homelab-1) exposes three exporters. Node-exporter for system metrics,
cAdvisor for container resource usage, and a Traefik metrics endpoint for the Docker-side
reverse proxy. These are static scrape targets in the Prometheus Helm values since they're
outside the cluster.

**Proxmox** metrics come through pve-exporter running on the Docker host. It talks to the
Proxmox API on both physical nodes and exposes VM-level CPU, memory, disk, and network
stats. I can see per-VM resource usage without logging into the Proxmox UI.

**The VyOS router** gets scraped via SNMP. An SNMP exporter pod runs inside the K3s cluster
with a Cilium network policy allowing it to reach the router at 10.10.10.1:161. It pulls
interface metrics using the IF-MIB module — throughput, errors, drops per interface.
Alongside that, goflow2 collects NetFlow data from VyOS for traffic flow analysis.

**The GPU node** (openclaw-1) runs an NVIDIA GPU exporter on the RTX 3060 that's passed
through via Proxmox PCI passthrough. Temperature, utilization, memory, fan speed, power
draw — all in Prometheus.

**The VPS** is the trickiest one. It's an Oracle Cloud instance running the Pangolin tunnel
endpoint. Node-exporter runs in Docker and its metrics are exposed through a Pangolin
tunnel ingress at `vps-metrics.vps.bp31app.com`. Prometheus scrapes it over HTTPS through
the tunnel — the same path external users take. CrowdSec metrics come through the same way.

**The media stack** uses Scraparr to export metrics from Sonarr, Radarr, and the rest of
the *arr suite. Download queue sizes, library counts, health check status — not
mission-critical, but fun to have on a dashboard.

---

## Logs

Loki handles log aggregation with 30-day retention.

**Inside the cluster**, Promtail runs as a DaemonSet and ships pod logs from every K3s node.
Standard setup — nothing fancy, but having logs next to metrics in the same Grafana instance
makes correlating issues significantly faster.

**On the VPS**, a separate Promtail instance collects syslog, auth logs, kernel logs, and
Docker container logs. It pushes to a Loki ingress endpoint at `loki-push.vps.bp31app.com`
through the Pangolin tunnel. This means I can search VPS logs from my Grafana without
SSH-ing into the VPS.

---

## 29 Alert Rules

I spent time getting alerts right because bad alerts are worse than no alerts. Every rule
has a threshold, a duration, and a severity. No alerts fire on transient blips.

**Node health** is the foundation — alerts for any node being down for 2 minutes, disk
usage over 80% (warning) or 90% (critical), and memory usage over 85%.

**Kubernetes alerts** watch for pods in CrashLoopBackOff (5 min) and pods not ready for
15 minutes. I exclude Jobs from the not-ready check because they're expected to terminate.

**Network alerts** cover VyOS being unreachable, BGP peers going down, the WAN interface
dropping, and high outbound traffic (>100Mbps sustained for 15 minutes — usually means
someone's downloading something large through the VPN).

**Storage alerts** fire when Longhorn volumes are degraded or faulted, and when any PVC
is over 85% full. These have caught real issues — a Prometheus PVC filling up because I
hadn't set retention limits.

**Infrastructure alerts** watch Traefik's 5xx error rate, TLS cert expiration (7-day
warning), ArgoCD apps out of sync or degraded, and high pod restart counts.

**VPS container alerts** track memory usage for Pangolin, Authentik, and CrowdSec. These
run on a small VPS and memory pressure is a real concern.

**Cilium network alerts** fire when Hubble sees more than 100 policy drops per second.
This has been genuinely useful — it caught a misconfigured network policy that was silently
dropping traffic between Uptime Kuma and Ntfy.

---

## 13 Dashboards, All in Terraform

Every dashboard is a JSON file in my homelab repo, deployed through the Grafana Terraform
provider. No clicking around in the UI and hoping someone exported the JSON before a
reinstall. I port-forward to Grafana, run `terraform apply`, and the dashboards are there.

I learned this the hard way. I originally used Helm's `gnetId` to import community
dashboards during provisioning. That worked until an upstream dashboard changed its panel
layout and my provisioning started failing. Now every dashboard is a committed JSON file
that I control.

The dashboards are organized across five folders:

**Infrastructure** — Proxmox cluster overview, VyOS router interfaces and SNMP metrics,
NVIDIA GPU utilization and temperature, and a homelab overview that ties everything together.

**K3s Cluster** — Node resources and pod distribution, Longhorn volume health and replica
status, Cilium Hubble network flows between namespaces, and interface-level network traffic.

**Docker/homelab-1** — Container CPU and memory per stack, Scraparr media metrics, and
Docker Traefik request rates and latency.

**Applications** — A dedicated K8s Learn dashboard showing active environments, resource
consumption per user namespace, pod status, and environment age (color-coded: green under
a day, yellow under a week, red over 30 days).

**VPS** — Authentik login metrics and CrowdSec ban/decision activity.

---

## The Notification Chain

When something goes wrong, the path is:

1. Prometheus evaluates alert rules every 30 seconds
2. Firing alerts route to Alertmanager
3. Alertmanager deduplicates, groups, and sends to receivers
4. **Critical alerts** → Ntfy (push notification on my phone)
5. **Warnings** → Slack

Uptime Kuma runs independently as a second opinion. It checks 16 services every 60 seconds
from inside the cluster, validating the full routing path (DNS → tunnel → Traefik →
service). If anything is down for 3 checks, it sends its own Ntfy notification on a
separate `uptime-alerts` topic.

Having two independent systems means I'll know about an outage even if one of them is the
thing that's broken.

---

## What I'd Do Differently

**Start with alerts, not dashboards.** I built dashboards first because they're satisfying
to look at. But dashboards only help when you're already looking. Alerts are what wake you
up when something's wrong at 2 AM. I should have defined alert rules first and built
dashboards around the things I was alerting on.

**Fewer scrape targets, more intentional ones.** Some metrics I'm collecting aren't useful.
I don't need 30-second resolution on my media stack. I should tune scrape intervals per
target rather than using the same interval for everything.

**Centralize exporter deployment.** My exporters are spread across Docker Compose stacks,
Kubernetes DaemonSets, and standalone containers on different VMs. If I were starting over,
I'd use a single Ansible role to deploy exporters consistently everywhere.

---

## The Numbers

| Metric | Value |
|--------|-------|
| Scrape targets | 20+ |
| Alert rules | 29 |
| Grafana dashboards | 13 |
| Log retention | 30 days |
| Metric retention | 15 days |
| Uptime monitors | 16 |
| Notification channels | 2 (Ntfy + Slack) |

It's probably more monitoring than a homelab needs. But the whole point of a homelab is
to build things the way you'd want them in production — and in production, you can never
have too much observability.

---

*Repos:*
- *[gitops](https://github.com/pattersonbl2/gitops) — Prometheus, alerting rules, Loki, ServiceMonitors*
- *[homelab](https://github.com/pattersonbl2/homelab) — Terraform dashboards, exporters, Docker stacks*
