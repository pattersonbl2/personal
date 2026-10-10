---
title: "Alert Triage Doesn't Need a Giant Context Window"
date: 2026-10-10T02:37:53-04:00
draft: false
summary: "Building sre-reflex: a shadow-mode bot that scores Prometheus alerts with a small decision model instead of a giant LLM context window, with human labels as ground truth."
tags: ["sre", "observability", "ai", "homelab"]
---

In [my last post](/posts/ai-engineering-in-the-loop/) I wrote about a platform that was silently broken for 24 hours, and the thing that bothered me most: *the failure notification that should have existed.*

I fixed that one specific gap. But it left me with a bigger problem. Once you add alerts to everything, you get a different failure mode: you stop reading them.

So I built **sre-reflex**, a small bot that scores every Prometheus alert before I have to decide whether to care. It doesn't use a big LLM with a giant context window. It uses a *decision model*, and it runs in shadow mode so it can't make anything worse while I find out whether it's any good.

The code is on GitHub: [sre-reflex](https://github.com/pattersonbl2/sre-reflex).

## The Problem With "Just Ask the LLM"

The obvious approach to alert triage is to dump everything into a prompt (the alert, the logs, the dashboards, the runbook) and ask a big model "what should I do?"

That works for a demo. For a stream of alerts it has three problems:

- **It's slow and expensive per alert**, and alerts arrive in bursts, exactly when you least want latency.
- **The answer is prose.** You can't threshold prose, graph it, or measure its calibration.
- **You can't tell if it's any good.** Without ground truth, "the summary looked reasonable" is the whole evaluation.

What I actually want from triage is three narrow answers:

| Question | Type |
|---|---|
| A human needs to take action on this alert. | yes/no probability |
| How severe is the impact? | 1–5 |
| This alert will resolve on its own. | yes/no probability |

Those are classification questions, not generation tasks. That's what a decision model is for.

## The Design

```
Alertmanager ──webhook──▶ sre-reflex ──▶ collectors (alert, Prometheus, Loki)
                              │
                              ├──▶ DecisionModel adapters ──▶ open-jev (GPU) / Ollama
                              ├──▶ Postgres (states, decisions, labels)
                              └──▶ ntfy (scores + ✅ Real / 🔇 Noise buttons)
```

When an alert fires, the bot builds a **pre-digested state of under 250 tokens**:

- the alert itself and how long it's been firing
- its 7-day history, e.g. *"fired 14 times, resolved 13, median duration 3m"*
- a short metric trend, e.g. *"error rate rose from 0.2% to 8.1%, sustained 6m"*
- deduplicated, scrubbed error log lines from Loki

Then it asks the three questions of two models, side by side:

- **openjev**, an Apache-2.0 decision model ([open-jev-deberta-v3-large](https://huggingface.co/com-kotobalabs/open-jev-deberta-v3-large)), a reproduction of the idea behind TypeSafe's Jev
- **ollama**, a local LLM baseline returning the same answers as JSON

Both models see the identical text. That's the whole point of the experiment.

## Decision 1: Do the Arithmetic in Code

The small decision model is weak at arithmetic. So the model never sees raw time series and has to work out "is this rising?" Instead, code computes *"error rate rose from 0.2% to 8.1%"* and hands the model a sentence.

This is a general rule I keep relearning: **use the model for judgment, use code for anything with a right answer.** Every number in the state is computed deterministically and the model only has to decide what it means.

The same thinking shaped the token cap. The model reads at most 256 state tokens, so the builder enforces 250 and trims log lines first, then metrics. If the two models saw different amounts of context, the comparison would be meaningless.

## Decision 2: Shadow Mode

sre-reflex never suppresses, reroutes, or edits an alert. My existing alerting is untouched. The bot posts its scores to a *separate* ntfy topic.

That gives me two properties I care about:

1. **If the bot dies, I lose nothing.** It's an observer, not a dependency.
2. **I don't have to trust it to start collecting data.** Trust is earned from the evaluation, not assumed.

Suppression and priority routing are explicitly v2, gated on the eval results justifying them.

## Decision 3: Make Labelling a Single Tap

An evaluation needs ground truth, and ground truth is the expensive part. So every score notification has two buttons: **✅ Real** and **🔇 Noise**.

The buttons are signed links (HMAC over alert id, value, and expiry), so the label endpoint can be exposed publicly without letting anyone write labels. Signatures expire after 7 days and can only be used once.

An hourly job also infers labels where it's safe: an alert that resolved within 10 minutes and that I never labelled gets marked as inferred noise. Inferred labels are stored separately from hand labels, so I can report on them separately and never silently mix the two.

## Decision 4: Interfaces at the Boundary

Every model sits behind a `DecisionModel` interface, and every adapter has to pass the same contract tests: one answer per question, distributions summing to 1, values in range. A deterministic fake adapter runs in CI.

That's not over-engineering for a homelab. It's what lets me add a third model later (TypeSafe's hosted Jev, when I get access) without touching the pipeline, and it's what makes the comparison fair.

## What the Evaluation Will Measure

`sre-reflex eval` replays stored states through each model and reports:

- precision, recall, and F1 for "actionable" against my hand labels
- **Brier score**, because a probability that's confidently wrong is worse than one that's honestly uncertain
- p50/p95 latency and cost per model
- performance **with and without incomplete context**, since collectors can time out and the model has to degrade gracefully
- agreement between the two models on questions with no ground truth

## Where It Stands

I'm not going to invent a results table. At the time of writing, sre-reflex is collecting shadow-mode data and I'm labelling alerts as they come in. The comparison report will be published in the repo's README after two weeks of labelled alerts.

What I can say now is what building it taught me, and it's the same theme as last time:

- **AI wrote a lot of the plumbing, and I still had to decide what "good" means.** Choosing the three questions, the 250-token budget, and the shadow-mode constraint were all judgment calls that code generation couldn't make for me.
- **Observability for the observer.** A triage bot that fails silently is just a new way to repeat the 24-hour outage. It exposes `/healthz` and `/metrics`, and it's built so that its failure is visible.
- **Measure before you trust.** I'd rather publish a report saying "the small model was no better than the LLM" than ship a triage layer on vibes.

When the labelled data is in, I'll write up the numbers, including the unflattering ones.
