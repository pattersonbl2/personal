---
title: "Learn Kubernetes by Breaking Things (On My Cluster)"
date: 2026-03-25
draft: false
---

I built a free Kubernetes playground that gives you a real cluster, a browser terminal, and five tasks that take you from "what's a Pod?" to deploying an app on the internet. No installs, no credit card, no Docker Desktop fighting your laptop for RAM.

**Try it now:** [learn.bp31app.com](https://learn.bp31app.com)

---

## How It Works

You type a name and email. Thirty seconds later, you have your own isolated namespace on a real K3s cluster, a browser-based terminal with kubectl ready to go, and a personal subdomain where you'll deploy your first app.

That's it. No VM to spin up, no YAML to copy-paste from a blog post before you understand what it does. You get a terminal and five tasks that build on each other.

---

## The Five Tasks

Each one teaches a core Kubernetes concept by having you do it yourself:

1. **Run a Pod** — Start a container. Watch it run. Understand that Kubernetes won't restart it if it dies.
2. **Create a Deployment** — Wrap that container in something that self-heals. Delete a pod and watch it come back.
3. **Expose it with a Service** — Give your pods a stable address. Learn why you can't just talk to pod IPs.
4. **Deploy to the Internet** — Create an Ingress and watch your app go live at your own subdomain within seconds.
5. **Scale and Observe** — Scale to five replicas, kill a few, and watch Kubernetes do what Kubernetes does.

By the end you've built something that's actually running on the internet, not a localhost demo that vanishes when you close the tab.

---

## What's New

Since the initial launch, I've made some solid improvements:

**Single binary architecture.** The platform used to run nginx pods for static pages plus a separate auth proxy. Now it's a single Go binary that serves the signup page, task reference, and handles authentication. Fewer moving parts, faster deploys.

**Terminal-inspired UI.** The signup and tasks pages got a full retheme — dark background, monospace fonts, the kind of aesthetic that feels right for a platform about running commands.

**Auth proxy with one-click login.** When you get your welcome email, the link includes a one-time token that logs you straight into your terminal. No password prompt on first visit.

**Better security.** Tighter RBAC, network policies via Cilium, rate limiting on signups, and automatic 24-hour cleanup so nobody's forgotten namespace is eating cluster resources.

**AI help button.** Stuck on a task? There's a "?" button that routes your question to a local LLM running on a GPU in my homelab. It knows the five tasks and will nudge you in the right direction without giving away the answer.

---

## Why I Built This

Most Kubernetes tutorials have you run commands against a pre-configured cluster. You type `kubectl apply -f deployment.yaml` and it works, but you don't really know why. You never see what happens when things go wrong because the tutorial environment doesn't let things go wrong.

I wanted the opposite. A real cluster where you create real objects, break things, fix them, and build intuition for how Kubernetes actually works. The kind of learning that sticks because you earned it.

---

## Give It a Shot

It takes 30 seconds to sign up and about an hour to work through all five tasks. Your environment lasts 24 hours, so you can come back to it if you need a break.

**[learn.bp31app.com](https://learn.bp31app.com)**

If you have feedback or want to see the code, the infrastructure is open source: [github.com/pattersonbl2/gitops](https://github.com/pattersonbl2/gitops)
