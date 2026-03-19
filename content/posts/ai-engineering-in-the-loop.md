---
title: "AI Didn't Fix My Broken Platform — We Did"
date: 2026-03-18T20:00:00-05:00
draft: false
---

My Kubernetes learning platform was silently broken for 24 hours. Every signup went into a black hole. No errors surfaced. No alerts fired. No one knew.

The irony? AI helped me build the platform. AI helped me harden the security. And AI helped me break it — without either of us realizing it.

Here's what happened, what it took to fix it, and why I think the engineers who thrive in the AI era won't be the ones who hand off tasks and walk away.

## The Platform

I built [K8s Learn](https://learn.bp31app.com) as a hands-on Kubernetes learning environment. You sign up, get an isolated namespace with a browser-based terminal, and work through five tasks: run a Pod, create a Deployment, expose a Service, set up an Ingress, and scale with rolling updates. An AI tutor is available if you get stuck.

Under the hood, the architecture looked like this:

```
User signup → nginx (rate limiting) → n8n webhook → Gitea Actions →
DinD runner → provision.sh → K8s API → namespace + terminal + email
```

It worked. Then I asked AI to help me harden the security. Lock down RBAC, add network policies, restrict container capabilities, remove the public K8s API tunnel. Standard stuff. Every change looked correct in isolation.

Then signups stopped working. Silently.

## Lesson 1: AI Doesn't Know What It Doesn't Know

I came back a day later and asked a simple question: "Has anyone signed up?"

The answer was zero user namespaces. Not because nobody tried — because the provisioning pipeline was broken. The n8n workflow was failing silently, the Gitea runner couldn't reach the K8s API, and there were no alerts for any of it.

AI had no way to catch this on its own. It didn't monitor the system after deploying changes. It didn't think to set up a synthetic test. It completed the security hardening, reported success, and moved on. From its perspective, the task was done.

This is the most fundamental gap in AI-assisted engineering: **AI operates on tasks, not outcomes.** It can tell you the RBAC change was applied successfully. It can't tell you that the change broke a downstream system that depends on permissions it just removed.

That requires someone who understands the system end-to-end. Someone who asks "is it actually working?" not just "did the apply succeed?"

## Lesson 2: Instrument Before You Debug

Once we knew provisioning was broken, AI started investigating. It checked n8n logs, Gitea workflow runs, container networking. It tried different fixes — changing Docker networks, swapping container images, adjusting environment variables.

After a few rounds of this, I stopped it: **"Add logging to the flow to help yourself."**

This is engineering discipline that AI doesn't naturally have. When you're staring at a failing system with no visibility, the instinct to start guessing is strong. AI feels this even more acutely — it wants to try things, see what happens, iterate. But without observability, you're just thrashing.

We added structured logging to the provisioning script:

```
[k8s-learn-provision] [INFO]  2026-03-18T05:07:31Z Testing cluster connectivity
[k8s-learn-provision] [ERROR] 2026-03-18T05:07:31Z Cannot connect to K8s cluster
```

We added an ERR trap that sends push notifications on failure. We added step numbering (Step 1/7 through Step 7/7) so you could see exactly where things died. We added HTTP status code logging for email sends.

After that, every problem was immediately obvious. The cluster connectivity error told us it was a networking issue. The RBAC error showed us the exact permission that was missing. The 401 on email send told us the API key was expired.

**The twenty minutes we spent on logging saved hours of guessing.** AI needed to be told this. An experienced engineer knows it instinctively: you don't debug what you can't see.

## Lesson 3: Domain Knowledge Cuts Through Complexity

Here's where it got interesting. The provisioning script ran inside a Docker-in-Docker container spawned by the Gitea Actions runner. That container couldn't reach the K8s API at `10.10.10.11:6443` because it was on an isolated Docker network.

AI went deep on this. It tried configuring the Docker network to `host` mode. It tried changing the runner config. It tried different API endpoints. Each attempt was logical in isolation, but none of them worked because DinD adds a layer of network isolation that's fundamentally hard to bridge.

I knew something AI didn't: **I had a Gitea runner already running inside the K3s cluster.** It was deployed months ago as a Kubernetes deployment with a Docker-in-Docker sidecar. Same problem, different location — the DinD containers inside the K3s pod also couldn't reach the node network.

So I made a call AI wouldn't have made on its own: **"Scrap the Gitea Actions approach entirely. Have n8n create a Kubernetes Job directly."**

```
# Before: 5 hops
signup → nginx → n8n → Gitea API → DinD runner → K8s API

# After: 3 hops
signup → nginx → n8n → K8s API (create Job) → runs in-cluster
```

The Job runs as a pod with the n8n service account. It has direct cluster access. No Docker-in-Docker. No runner networking. No Gitea Actions. The provisioning script went from failing for 24 hours to completing in 17 seconds.

This is the kind of decision that requires knowing your infrastructure intimately. AI can optimize within an architecture. Humans decide when to change the architecture.

## Lesson 4: AI Optimizes for Speed; Engineers Optimize for Consequences

During the debugging, we hit a wall with n8n's sandboxed code runner. The workflow used `process.env.GITEA_API_TOKEN` to authenticate with Gitea, but n8n 2.x runs Code nodes in a sandbox that blocks environment variable access.

AI's first suggestion: disable the sandbox by setting `N8N_RUNNERS_DISABLED=true`.

I pushed back: **"This would break security."**

That sandbox exists for a reason. n8n workflows can be created by users, and giving them access to environment variables means giving them access to API tokens, database credentials, and everything else in the container's environment.

We found the right fix: replace the Code node entirely with n8n's built-in HTTP Request node, which handles authentication through n8n's encrypted credential store. The token never appears in workflow code. It's stored encrypted in the database and only decrypted at request time.

The quick fix would have worked. It would have taken 30 seconds. And it would have opened a security hole that could persist for months before anyone noticed. **AI gave me the fastest path. I chose the safest one.** That distinction matters more than people think.

## Lesson 5: AI Doesn't Internalize Operational Patterns

This one was subtle and cost us the most time.

Our Kubernetes manifests are managed by ArgoCD with `selfHeal: true`. This means ArgoCD continuously watches the cluster and reverts any changes that don't match what's in the Git repository. It's a core GitOps principle: the repo is the source of truth.

AI applied changes directly with `kubectl apply` three separate times. Each time, ArgoCD reverted the changes within seconds. AI would verify the change was applied, move on to the next step, and then the previous step would silently undo itself.

I had to intervene: **"Stop. ArgoCD keeps reverting that. You need to commit to the repo and let ArgoCD sync."**

This isn't a knowledge gap — AI knows what GitOps is. It's an operational pattern gap. In the heat of debugging, with a broken system and mounting pressure to fix it, the instinct is to apply changes directly and see if they work. That's fine in a non-GitOps environment. In ours, it's actively counterproductive.

After I flagged it, AI saved this as a persistent rule: "Never use kubectl apply on resources managed by ArgoCD." It won't make this mistake again in our project. But it needed a human to catch the pattern and explain why it mattered.

## Lesson 6: The Failure Notification That Should Have Existed

The thing that bothers me most about this incident isn't the technical complexity. It's that the system was broken for 24 hours and nobody knew.

When we finally added failure notifications — a push notification via ntfy that fires whenever the provisioning script hits an error — I realized this should have been there from the start. Not just for provisioning, but for everything.

AI built the happy path beautifully. The signup form, the progress bar, the welcome email, the terminal deployment. But it didn't build the sad path. No failure notifications. No dead letter queue. No synthetic monitoring. No "hey, zero signups succeeded today" alert.

This is a pattern I see across AI-assisted development: **the happy path gets built to production quality. The failure modes get ignored.** Not because AI can't build them — it absolutely can. But because building failure handling requires imagining what can go wrong, and AI tends to be optimistic about code it just wrote.

## What I Actually Think About AI in Engineering

I don't think AI is going to replace engineers. I also don't think it's a toy. After tonight, here's my honest assessment:

**AI is extraordinary at:** writing provisioning scripts, RBAC policies, network policies, Dockerfiles, Helm charts, sealed secrets, nginx configs, structured logging, kustomize overlays. It wrote hundreds of lines of correct, production-quality infrastructure code in minutes. Things that would have taken me hours.

**AI is bad at:** knowing when the system is actually working. Recognizing when to stop fixing and start observing. Knowing your infrastructure well enough to make architectural pivots. Prioritizing security over expediency. Internalizing operational patterns that aren't in the code.

**The engineers who will thrive** are the ones who stay close enough to the work to catch the moments where judgment matters. Not micromanaging every line of code — that defeats the purpose. But staying in the loop enough to say:

- "Is it actually working, or did we just deploy successfully?"
- "Stop guessing. Add logging first."
- "I know a simpler path. Let's change the architecture."
- "That fix is fast but it breaks our security model."
- "That won't stick — our GitOps controller will revert it."

These aren't coding skills. They're engineering skills. And they're more valuable now than ever, precisely because AI handles so much of the coding.

## The Result

By the end of the session, the platform was fully operational:

- **17-second provisioning** (down from broken)
- **Structured logging** with step-by-step visibility and push notifications on failure
- **Hardened security** — non-root containers, read-only filesystems, network policies, scoped RBAC
- **Simplified architecture** — eliminated Gitea Actions, DinD, and two layers of networking complexity
- **Full GitOps persistence** — everything managed by ArgoCD, sealed secrets for credentials

AI wrote most of the code. I made most of the decisions. That's the partnership that works.
