---
title: "Defining SRE Observability at Jasper — 42% Through Our Datadog Optimization"
date: 2026-01-28T10:00:00-05:00
draft: false
---

I’ve been leading an effort to mature observability at Jasper: we’re about **42% through a Datadog optimization project** aimed at standardizing how we instrument, monitor, and alert on our systems. This post is a short update on how we’re defining SRE observability practices and the approach that’s working so far.

## Why we’re doing this

After a period of fast growth, we had a lot of ad‑hoc dashboards, inconsistent naming, and gaps in coverage. The goal is to turn that into a clear, repeatable model: consistent metrics, logging, and alerting so that both SRE and product engineering can understand system behavior and respond to issues quickly.

## How we’re approaching it

**Working closely with the dev team.** Observability only works if it’s aligned with how services are built and deployed. I’ve been sitting with engineers to understand service boundaries, critical paths, and where we need better visibility. That partnership is what makes the standards stick.

**Leveraging codebase knowledge and AI.** We have a large, evolving codebase. I’ve been combining deep familiarity with the repo and AI-assisted tooling to draft **templated patterns**: standard dashboard layouts, metric naming conventions, and alert definitions that teams can reuse. The templates keep things consistent without each team reinventing the wheel; the AI helps iterate on drafts and explore edge cases quickly.

**Defining the practices, not just the tools.** The real deliverable isn’t “more Datadog” — it’s a set of SRE observability practices we can document and follow: what we measure, how we name it, when we page, and how we review and tune over time. Datadog is the implementation; the practices are the standard.

## Where we are

We’re a little past the halfway mark in terms of project scope. So far we’ve:

- Aligned on tagging and naming conventions across services
- Rolled out template dashboards and started migrating existing ones
- Tightened alert definitions and reduced noise
- Documented the patterns so new services can onboard consistently

The remaining work is mainly rolling these patterns out to the rest of the stack and refining based on feedback from incidents and weekly reviews.

I’ll share more as we get closer to completion — and hopefully a short “lessons learned” post once we’re there. If you’re doing something similar (Datadog, observability standardization, or SRE practice definition), I’d be interested to hear how you’re structuring it.
