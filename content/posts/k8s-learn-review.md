---
title: "K8s Learn: A Self-Hosted Kubernetes Training Platform"
date: 2026-03-15T12:00:00-05:00
draft: false
---

K8s Learn is a hands-on Kubernetes training platform I built because most K8s tutorials give
you a pre-configured cluster and tell you to run commands — I wanted something that forces you
to reason about what you're doing and why. It's live at [learn.bp31app.com](https://learn.bp31app.com)
if you want to try it yourself. Deployed on my homelab K3s cluster via ArgoCD, the stack is
K3s, ArgoCD, n8n, Ollama, and Cloudflare Tunnel.

---

## The Frontend

I wanted a distraction-free training environment — no heavy framework, no SPA overhead, just
tasks and a terminal. The frontend is static HTML/CSS served by nginx:alpine inside K8s. Five
progressive tasks walk you through the fundamentals: Pod → Deployment → Service → Ingress → Scale.

The task sequence is deliberate. You start with a bare Pod to understand that K8s doesn't manage
it for you, then wrap it in a Deployment to see why self-healing matters, expose it with a Service
to learn the selector model, route external traffic through an Ingress, and finally scale it to
internalize how replicas and scheduling interact. Each task builds on the last rather than
standing alone.

The HTML is ConfigMap-mounted via kustomize `configMapGenerator` — no image rebuilds when content
changes, just an ArgoCD sync. Resources are deliberately minimal: 10m CPU / 16Mi RAM requests,
50m CPU / 32Mi RAM limits. It's not that the content is simple; it's that static HTML doesn't
need headroom.

---

## User Provisioning

Every user gets an isolated namespace with scoped RBAC — but provisioning that automatically
required a pattern I hadn't used before: an external service managing K8s resources directly via
the API. n8n runs outside the cluster (Docker on VM 110) and authenticates with a long-lived
ServiceAccount token stored in n8n's credential store. When a signup webhook fires, n8n creates
the namespace, ServiceAccount, RoleBinding, and kubeconfig in one workflow — no Terraform, no
Helm, just direct API calls.

The ClusterRole (`n8n-k8s-learn-admin`) grants full CRUD on namespaces, pods, deployments,
services, ingresses, secrets, serviceaccounts, and RBAC resources. This pattern — a workflow engine holding K8s credentials
and reacting to events — turns out to be a clean way to bridge external triggers with cluster
operations without running an in-cluster operator. The signup endpoint is rate-limited at nginx:
2 requests/minute per IP with burst=3. The signup page has a multi-step progress bar —
Creating namespace → Deploying terminal → Environment ready — backed by polling so users see
real status rather than a spinner while n8n does its work.

---

## Per-User Routing

Each user gets their own subdomain for the Ingress task — the point being that they create the
Ingress object themselves and watch it go live. The Cloudflare Tunnel catch-all for `*.bp31app.com`
routes everything to Traefik inside the cluster, so when a user runs:

```bash
kubectl create ingress my-app --rule="learn-alice-app.bp31app.com/*=my-service:80"
```

Traefik picks it up automatically and the subdomain resolves within seconds. No tunnel config
changes, no DNS propagation wait.

One constraint worth noting: Cloudflare's free plan only supports one level of wildcard subdomain.
That means user subdomains follow the pattern `learn-[name]-app.bp31app.com` rather than something
cleaner like `alice.learn.bp31app.com`. It's a real architectural constraint that shaped the
subdomain naming convention — and a useful lesson in working within platform limits.

---

## AI Help Widget

There's a floating "?" button on the tasks page. Click it, ask a question, and the request routes
through an n8n webhook to Ollama running `qwen2.5:14b` on a dedicated GPU node — an RTX 3060
passed through via Proxmox PCI passthrough.

The system prompt is scoped to the five tasks; off-topic questions get a redirect rather than a
hallucinated answer. If the model is slow, nginx holds the connection open for up to 60 seconds
before timing out — long enough for most responses, and the UI shows a loading state so it
doesn't feel broken. If the rate limit kicks in (6 requests/minute per IP, burst=2), the user
gets a clear "slow down" message rather than a silent failure. The GPU node handles one request
at a time, so the rate limiting also protects inference latency for concurrent users.

---

## Cleanup

User environments auto-expire after 24 hours. A scheduled n8n workflow handles deletion —
namespace, RBAC, everything — using the same ClusterRole permissions that provisioned them.
No manual intervention, no orphaned resources accumulating on the cluster.

---

## What's Next

I'd like to add more task tracks — storage, ConfigMaps, secrets management with external-secrets —
and eventually progress persistence so users can pick up where they left off. An OAuth signup flow
(Authentik is already running in the homelab) would replace the current webhook-based provisioning
and give users a real identity rather than just a namespace name. A live leaderboard showing who
completed which tasks fastest would also make it more engaging for workshop use.

*Check out the repo: [gitops](https://github.com/pattersonbl2/gitops)*
