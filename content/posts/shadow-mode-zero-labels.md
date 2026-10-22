---
title: "Twelve Days of Shadow Mode, Zero Labels, and One Line of YAML"
date: 2026-11-01T09:00:00-04:00
draft: false
summary: "I promised a precision/recall table after two weeks of labelled alerts. I got an empty one, because my evaluation harness had been silently broken the whole time. Here is what the data did say."
tags: ["sre", "observability", "ai", "homelab"]
---

I ran `sre-reflex eval` after twelve days of shadow mode and got an empty table. 89 alerts, 89 states, 468 model decisions, and **0 labels**. No precision, no recall, no Brier score. The cause was one line of YAML in my Alertmanager receiver.

In [my last post](/posts/alert-triage-decision-models/) I built sre-reflex: a bot that scores every Prometheus alert with a small decision model instead of a giant LLM context window, running in shadow mode so it could not make anything worse. I ended that post with a promise — when the labelled data was in, I would write up the numbers, including the unflattering ones.

This is the unflattering one. The data never came in.

## The Headline Is a Zero

Here is what twelve days (10–22 October) actually produced:

| Table | Rows |
|---|---|
| alerts | 89 |
| states | 89 |
| decisions | 468 |
| labels | **0** |

Two separate failures hide in that last row.

The first is mine and it is boring: I did not tap **Real** or **Noise** on the notifications. Not once, in twelve days. The ground truth was supposed to be one tap per alert and I never built the habit.

The second one is a bug, and it is the interesting half. I had a backstop for exactly that laziness — an hourly job that infers "noise" for any alert that resolved inside ten minutes. It needs a `resolved_at` timestamp to do that. All 89 alerts have `resolved_at = NULL`.

The reason is one line in the Alertmanager receiver:

```yaml
sendResolved: false
```

The bot is never told that an alert resolved. So it cannot infer noise, and it also cannot grade its own "will this resolve on its own?" prediction — the one question of the three where the ground truth should have been free.

So there is no precision, recall, F1 or Brier score in this post. Those need ground truth and I do not have any. I would rather say that than publish a table I cannot back up.

## What the Data Does Support

The pipeline works. These are measurements, not quality scores — the difference matters, and conflating the two is how you end up trusting a bot you have not graded.

**Volume and mix.** 89 alerts across 12 distinct alert names. Four of them are 73% of the traffic: NodeMemoryUsageWarning (28), KubeJobFailed (16), TargetDown (11), NodeDiskUsageWarning (10). That is a finding on its own. Most of what fires at me is a handful of repeat offenders, and the long tail I was worried about missing is barely a tail.

**Context quality.** 76 of 89 states (85%) had every collector working. 12 were missing Loki logs. One was missing logs, history and metrics — a state so thin the models were guessing from the alert name. State size averaged 134 tokens, max 176, comfortably under the 250-token budget I set myself.

**Coverage.** openjev answered 75 of 89. Ollama answered 81. Both failed to return a usable answer on some states and I have not yet worked out why. That is a real gap: a triage model that silently declines 16% of alerts is a triage model with a hole in it.

**Latency and cost.** Median about 2.4 seconds for both. The tails differ — p95 is 4.2s for openjev and 7.8s for the Ollama LLM. Cost is $0 either way, since both run locally on hardware I already own.

**Behaviour.** The two models do not think alike:

| | openjev | ollama |
|---|---|---|
| "Actionable" at 0.5 or above | 75 of 75 (100%) | 66 of 81 (81%) |
| Mean actionable score | 0.65 | 0.78 |
| "Self-resolving" at 0.5 or above | 62 of 75 (83%) | 53 of 81 (65%) |
| Most common severity | 2 (61%) | 3 (60%) |

They agree on severity 40% of the time and on self-resolving 77% of the time.

## The Red Flag I Can Name Without Labels

I cannot tell you which model is right. I can tell you that openjev called **every single alert** actionable. 75 out of 75.

A triage model that never says "ignore this" is not triaging. It is forwarding.

There are two explanations and I cannot pick between them. Either the model is miscalibrated for my alert mix, or my alerts genuinely are all worth a look — which, given that four alert names make up 73% of them, I doubt. Labels would tell me which. I do not have labels.

Note also the shape of the disagreement: openjev is *more* likely to call something actionable and *more* likely to call the same thing self-resolving. Those two together describe a model hedging on both questions at once, which is what a model does when the state it is reading does not actually determine the answer.

## Observability for the Observer, Again

The last post ended on a principle I was pleased with at the time:

> A triage bot that fails silently is just a new way to repeat the 24-hour outage.

sre-reflex exposes `/healthz` and `/metrics`. It was up the whole time. Nothing was down. Every component reported healthy, and the system still produced nothing of value for twelve days, because "healthy" was measuring the wrong thing.

I monitored the bot. I did not monitor the *evaluation*. There was no alert for "zero labels in 24 hours", no alert for "100% of alerts have a null resolved_at", no alert for "the hourly inference job ran and inferred nothing, again, for the 288th time". Every one of those is a one-line Prometheus rule, and the bot that was built to catch silent failures had none of them pointed at itself.

## What I'm Changing

1. Set `sendResolved: true` on the receiver, so `resolved_at` gets recorded and inferred noise starts working.
2. Add the three rules above. If the label count stays flat for a day, I want to be told.
3. Start tapping Real / Noise. Zero hand labels in twelve days is on me, not on the YAML.
4. Re-run the eval after another two weeks and publish the real precision, recall and Brier numbers, whatever they say.

v2 stays gated. Suppression and priority routing were explicitly conditional on the eval justifying them, and an empty eval justifies nothing. The bot stays in shadow mode. That is not a setback — it is the gate doing exactly the job I built it to do.

## What This Taught Me

- **An evaluation harness needs its own monitoring.** It is a production system. It has inputs that can go to zero and nobody notices, because its output is a report that nobody reads until the day you need it.
- **Free ground truth is not free if you never wire it up.** The self-resolving signal was sitting in Alertmanager the entire time, one boolean away.
- **Measuring the model measured me instead.** I set out to learn whether a small decision model could triage my alerts. I learned that I will not do one tap a day without being nagged, and that I shipped a measurement system I never tested end to end.

I built a bot to catch silent failures. It spent twelve days quietly unable to measure itself, and the only reason I found out is that I had promised a blog post with numbers in it.

Deadlines are a monitoring system too. Not a good one.
