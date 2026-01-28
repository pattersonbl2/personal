---
title: "Using AI to Deepen My Platform Building Journey"
date: 2026-01-27T12:00:00-05:00
draft: false
---

As someone transitioning from infrastructure operations to a more development-focused SRE role, I've been on a mission to build a production-grade homelab that serves as both a learning environment and a practical platform. What started as a simple Docker setup has evolved into a full GitOps infrastructure with Kubernetes, ArgoCD, monitoring, and secure external access. But the real game-changer has been integrating AI assistance—specifically Claude—into my workflow to accelerate learning and deepen my understanding of platform engineering.

## The Challenge: Building While Learning

When I left Mozilla in October 2025, I realized I had a significant gap: while I was comfortable managing infrastructure at scale, I hadn't been writing code regularly for years. My homelab project became the perfect bridge—a place to practice infrastructure-as-code, learn Kubernetes deeply, and build the kind of platform I'd want to operate professionally.

The challenge? Building a sophisticated platform while simultaneously learning the technologies. I needed to:
- Set up a K3s cluster with proper networking (Cilium with BGP)
- Implement GitOps with ArgoCD
- Configure secure external access via Cloudflare Tunnel
- Build monitoring and observability
- Write Terraform and Ansible automation
- Learn Python for SRE/backend engineering

That's a lot to juggle, especially when you're trying to understand not just *how* to do something, but *why*.

## Enter AI: My Pair Programming Partner

I started using Claude (via claude.ai/code) as a pair programming partner for my infrastructure work. The key was creating a `CLAUDE.md` file in my homelab repository that provides context about my setup, architecture, and common patterns. This allows Claude to understand my specific environment and provide relevant, actionable guidance.

### Infrastructure Context and Guidance

My `CLAUDE.md` file documents:
- Repository structure and purpose
- Common Terraform and Ansible commands
- Architecture patterns (deployment flow, network layout, K3s cluster structure)
- Key configuration locations

When I'm working on a new feature or debugging an issue, I can share the relevant files and get context-aware suggestions. For example, when setting up Cilium BGP routing, Claude helped me understand:
- Why we disable Flannel in K3s (`--flannel-backend=none`)
- How BGP peering works with my VyOS router
- The relationship between pod networks (10.42.0.0/16) and service networks (10.43.0.0/16)

This isn't just copy-paste solutions—it's learning *why* things work the way they do.

### GitOps Infrastructure Deep Dive

My GitOps repository is a complete production-like setup with:
- **K3s** cluster with Cilium CNI
- **ArgoCD** for continuous delivery
- **Traefik** as ingress controller
- **Cloudflare Tunnel** with subnet routing
- **Crossplane** for infrastructure-as-code
- **Prometheus/Grafana** monitoring stack

When I was setting up Cloudflare Tunnel with subnet routing (instead of individual hostname routes), Claude helped me understand:
- The advantages of subnet routing for Tailscale-only access
- How to configure routes in Cloudflare Zero Trust
- The relationship between Kubernetes service IPs and tunnel routing
- Troubleshooting TLS handshake issues (ArgoCD runs on HTTP internally, Cloudflare handles TLS)

The AI didn't just give me config snippets—it explained the architecture decisions and trade-offs.

### Learning Python While Building

One of my biggest gaps was active development skills. I created a `python-sre-guide` in my homelab docs with structured learning paths, but I also use AI to:
- Explain Python concepts in the context of SRE work
- Review my code and suggest improvements
- Help me understand algorithms and data structures
- Practice coding interview patterns

The key is asking for explanations, not just solutions. "Why does this pattern work?" "What are the trade-offs?" "How would this scale?"

## The Workflow: AI-Assisted Platform Building

Here's how I typically work:

1. **Plan the Feature**: I outline what I want to build (e.g., "Add monitoring to my K3s cluster")

2. **Research with AI**: I ask Claude to explain the concepts and best practices, referencing my existing setup

3. **Implement with Guidance**: As I write Terraform/Ansible/Kubernetes configs, I get real-time feedback on:
   - Best practices
   - Security considerations
   - Integration with existing infrastructure
   - Troubleshooting when things don't work

4. **Document and Learn**: I update my README files and `CLAUDE.md` with what I learned, creating a knowledge base for future work

## Real Examples: What AI Helped Me Build

### 1. GitOps Bootstrap Setup

Setting up ArgoCD to manage itself and other applications required understanding:
- Application CRDs and sync policies
- Helm chart dependencies
- Kustomization patterns

AI helped me understand the GitOps pattern and how to structure my applications for proper dependency management.

### 2. Cloudflare Tunnel Subnet Routing

Instead of configuring individual hostname routes, I wanted subnet routing for better scalability. AI helped me:
- Understand the difference between hostname and subnet routing
- Configure routes in Cloudflare Zero Trust
- Set up Tailscale-only access policies
- Troubleshoot connection issues

### 3. Cilium BGP Configuration

Getting LoadBalancer services to work with my VyOS router required BGP peering. AI helped me:
- Understand BGP fundamentals in the Kubernetes context
- Configure Cilium BGP values
- Debug BGP peer status
- Understand the relationship between router and cluster networking

### 4. Local CI/CD with Git Hooks

I built an automated deployment system that detects changes and deploys via Ansible. AI helped me:
- Design the change detection logic
- Write bash scripts for git hooks
- Structure Ansible playbooks for idempotency
- Handle edge cases and error scenarios

## The Learning Multiplier

The real value isn't just getting things done faster—it's understanding *why* things work. When AI explains a concept in the context of my actual infrastructure, I internalize it better than reading generic documentation.

For example, when troubleshooting why ArgoCD wasn't accessible via Cloudflare Tunnel, I learned:
- How Kubernetes services expose HTTP vs HTTPS
- How Cloudflare Tunnel handles TLS termination
- The difference between proxied and DNS-only Cloudflare records
- How subnet routing simplifies access patterns

This deep understanding makes me better at:
- Debugging issues in production-like environments
- Making architectural decisions
- Explaining concepts to others
- Interviewing for roles that require both infrastructure and development skills

## Building a Learning Platform

My homelab isn't just infrastructure—it's a learning platform. Every component teaches me something:
- **Terraform** → Infrastructure-as-code patterns
- **Ansible** → Configuration management and idempotency
- **K3s** → Kubernetes internals without cloud provider abstractions
- **ArgoCD** → GitOps and continuous delivery
- **Cilium** → CNI, networking, and BGP
- **Cloudflare Tunnel** → Secure access patterns
- **Monitoring** → Observability and alerting

AI helps me connect these dots and understand how they fit together in a production system.

## The Future: AI as a Force Multiplier

As I continue building, I'm using AI to:
- Generate documentation from my infrastructure code
- Create runbooks for common operations
- Write tests for my automation
- Design new features with best practices in mind
- Learn new technologies faster

The goal isn't to replace understanding with AI—it's to use AI to accelerate learning and deepen comprehension. Every interaction teaches me something, and that knowledge compounds over time.

## Takeaways

If you're building a homelab or learning platform engineering:

1. **Document your setup**: Create context files (like `CLAUDE.md`) that help AI understand your environment
2. **Ask for explanations**: Don't just ask for solutions—ask *why* things work
3. **Build incrementally**: Use AI to understand each component before moving to the next
4. **Connect concepts**: Ask how different technologies relate to each other
5. **Learn by doing**: Use AI guidance to implement, then understand what you built

The combination of hands-on building and AI-assisted learning has been transformative. I'm not just copying configs—I'm building a deep understanding of platform engineering that will serve me well as I continue my career transition.

My homelab is now a production-grade platform that I understand inside and out, and AI has been an essential partner in that journey. The infrastructure is the foundation, but the knowledge and understanding are what will carry forward.

---

*Want to see the infrastructure? Check out my [homelab](https://github.com/pattersonbl2/homelab) and [gitops](https://github.com/pattersonbl2/gitops) repositories.*
