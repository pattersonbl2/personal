---
title: "Hardening K8s Learn: Security, Monitoring, and Ditching KSOPS"
date: 2026-03-16T17:00:00-05:00
draft: false
---

The [K8s Learn platform](/posts/k8s-learn-review/) launched a few days ago and immediately exposed
a handful of problems — the kind you only find when something is live and you start poking at it.
This post covers the security fixes, monitoring additions, and a secret management migration that
simplified the entire GitOps stack.

---

## Security Hardening

The original deployment had three issues worth fixing.

**Overprivileged ClusterRole.** The n8n ServiceAccount had a ClusterRole granting full CRUD on
pods, deployments, services, ingresses, secrets, serviceaccounts, and RBAC resources cluster-wide.
That's far more access than a provisioning workflow needs. I replaced it with a `n8n-namespace-manager`
ClusterRole scoped to namespace operations only — `get`, `list`, `create`, `delete` on namespaces.
Per-namespace permissions (the ones users actually need) are now granted via a Role and RoleBinding
created inside each `learn-*` namespace during provisioning. If the n8n token were compromised, the
blast radius is now limited to creating and deleting namespaces rather than full cluster access.

**Webhook bypass route.** The Cloudflare Tunnel config included a direct route to n8n
(`learn-webhook.bp31app.com`) that bypassed the nginx rate limiter entirely. Anyone could hit the
signup endpoint directly at unlimited volume. Removed the route — all traffic now flows through
the nginx proxy where rate limiting is enforced at 2 requests/minute per IP.

**Stored XSS in signup page.** The signup page read a username from `localStorage` and injected
it into the DOM without sanitization. A crafted username like `<img src=x onerror=alert(1)>` would
execute on every page load. Fixed with proper escaping before insertion.

I also removed the `api.kube.bp31app.com` Cloudflare Tunnel rule that exposed the Kubernetes API
externally. There was no reason for the API server to be reachable outside the homelab network.

---

## WebSocket Fix for Terminals

The in-browser ttyd terminals weren't connecting through Cloudflare Tunnel. The issue was
`cloudflared`'s http2 protocol mode dropping WebSocket upgrade requests with "unsupported connection
type." Adding `disableChunkedEncoding: true` to the wildcard ingress rule in the tunnel config
fixed it — WebSocket connections now proxy cleanly through the tunnel to Traefik and on to the
ttyd pods.

---

## Monitoring Dashboard

I built a Grafana dashboard for the platform, provisioned via Terraform. The top row has six stat
panels: active environments, total pods, CPU usage, memory usage, pod restarts, and failed pods.
These give an at-a-glance health check without drilling into individual namespaces.

Below that, a table merges data from four Prometheus queries to show each `learn-*` namespace with
its age, pod count, CPU, and memory side by side. The age column was initially broken — showing "56"
instead of a relative time — because `kube_namespace_created` returns epoch seconds but Grafana's
`dateTimeFromNow` expects milliseconds. Multiplying by 1000 in the query fixed it.

Time series panels break down CPU and memory by namespace, track pod status (running/pending/failed),
and chart container restart rates. Network I/O panels show receive and transmit bytes per namespace.
An "Environments by Age" bar gauge color-codes namespaces: green under a day, yellow under a week,
red over 30 days — useful for catching environments that outlive the 24-hour cleanup window.

---

## Alerting

Prometheus routes critical alerts to both [ntfy](https://ntfy.sh) (self-hosted, for push
notifications) and Slack. The alertmanager config lives in a sealed secret and includes receivers
for each channel. Non-critical alerts go to Slack only.

Getting ntfy working required adding a Bearer token to the webhook — ntfy runs with `deny-all`
default access, so unauthenticated requests were returning 403. I also disabled the kubelet
recording rules from kube-prometheus-stack. K3s doesn't expose metrics like
`kubelet_certificate_manager_client_ttl_seconds` that those rules expect, so they were generating
constant `PrometheusRuleFailures` alerts — noise, not signal.

---

## Ditching KSOPS for Sealed Secrets

This was the biggest change. The cluster used KSOPS (Kustomize + SOPS) for managing secrets in Git:
GPG-encrypted YAML files decrypted at sync time by a sidecar container on ArgoCD's repo-server. The
idea is sound — encrypt secrets with GPG, commit the ciphertext, decrypt during kustomize build. In
practice, it was the most fragile part of the stack.

The problems were layered. The KSOPS Docker image didn't include `gpg`, so the sidecar needed to use
the ArgoCD base image instead, with KSOPS binaries copied in via init container. Those binaries had
to be placed in an exact kustomize plugin directory structure
(`$KUSTOMIZE_PLUGIN_HOME/viaduct.ai/v1/ksops/ksops`). The GPG keyring needed to be writable (not
read-only) so `gpg` could create lock files during decryption. The CMP ConfigMap that defined the
plugin had to be Helm-managed to avoid ArgoCD reporting it as out of sync. Every one of these was a
separate failure mode that produced opaque errors — "unable to find plugin root," "0 successful
groups required," "Read-only file system" — none of which pointed clearly at the actual problem.

I replaced the entire setup with [Sealed Secrets](https://github.com/bitnami-labs/sealed-secrets).
The controller installs as a standard Helm chart, generates its own key pair on the cluster, and
requires no sidecar. Secrets are encrypted locally with `kubeseal`, committed as `SealedSecret` CRDs,
and decrypted into regular `Secret` objects by the controller. ArgoCD handles them natively.

The migration: decrypt each SOPS secret, pipe through `kubeseal`, replace the kustomize generator
with a plain resource reference, strip the sidecar from ArgoCD's values. Four secrets migrated
(cloudflare-tunnel credentials, alertmanager config, personal API secrets, ArgoCD repo credentials),
the GPG key deleted, and the repo-server went from a two-container pod back to one.

The repo is cleaner too. No more `secret-generator.yaml` files, no `.sops.yaml` config, no
`secrets/` directories with encrypted YAML. Each app's kustomization just lists a
`sealed-*.yaml` file as a resource, same as any other manifest.

---

## What Changed

| Before | After |
|--------|-------|
| Cluster-wide ClusterRole | Namespace-only ClusterRole + per-namespace Roles |
| Direct webhook route bypassing rate limit | All traffic through nginx |
| Unsanitized DOM insertion | Escaped localStorage values |
| K8s API exposed via tunnel | Removed |
| No monitoring dashboard | 15-panel Grafana dashboard |
| KSOPS sidecar + GPG + CMP plugin | Sealed Secrets controller |
| 2-container repo-server pod | 1-container repo-server pod |

*Check out the repo: [gitops](https://github.com/pattersonbl2/gitops)*
