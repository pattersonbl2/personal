---
title: "What Shipping My Personal Site Actually Taught Me"
date: 2026-04-03
draft: false
---

Most of the interesting work on my personal site this year was not glamorous.

It wasn’t a flashy redesign or a giant feature launch. It was the kind of work that tends to disappear into commit history: fixing deployment edge cases, tightening security, cleaning up frontmatter, wiring up backend behavior, and making the site feel a little more intentional each time I touched it.

Looking back through my recent GitHub activity, that’s the real pattern I see. A lot of small changes, but all in service of the same goal: make the site easier to trust, easier to maintain, and less likely to surprise me later.

That feels worth writing about, because it’s a more honest picture of engineering than “I built a thing” posts usually capture.

## The work behind a personal site is mostly systems work

A personal site looks simple from the outside. A homepage, a few blog posts, a resume, a contact page.

But once you actually own the thing end to end, it turns into a small platform:

- content needs to build correctly
- deployment needs to behave differently in preview and production
- form endpoints need rate limiting and input validation
- secrets need to stay out of the repo
- dates, URLs, and metadata need to be right or the site quietly does the wrong thing
- the visual layer still needs to feel coherent

That’s what I see in the repo history. Not one big "site rewrite," but a steady stream of practical improvements:

- deployment workflow updates
- Cloudflare Pages fixes
- security hardening
- backend additions for contact and resume handling
- favicon and theme polish
- blog post publishing fixes after Hugo date issues

None of those changes are individually dramatic. Together, they’re the work of turning a personal website into something that behaves like a maintained system instead of a static side project.

## Small bugs on personal projects still teach real lessons

One pattern that jumped out immediately was how many fixes were about correctness at the edges.

For example, I had to correct post frontmatter dates because Hugo was treating some content as future-dated and filtering it out. That’s a tiny detail until it isn’t. You can write the post, commit the file, push the change, and still end up with missing content because a date format was just slightly wrong.

That kind of issue is a good reminder that publishing systems are opinionated in ways that don’t always announce themselves. A lot of engineering is just learning where those sharp edges live.

Same story with environment variables and backend behavior. One recent fix trimmed whitespace from a `RESUME_TOKEN` because an invisible extra character was enough to turn a valid URL into a broken one. Again: not glamorous, very real.

I like this kind of bug because it reinforces something I’ve learned repeatedly in infrastructure and platform work:

**systems rarely fail in impressive ways.**

More often, they fail through formatting, assumptions, sequencing, or glue code.

That applies just as much to a personal website as it does to a production service.

## Personal sites are a good place to practice restraint

The repo also reflects something I’m trying to get better at: not overbuilding.

It would be very easy to turn a personal site into an endless framework experiment. New theme. New stack. New CMS. New frontend. New rewrite every two months.

Instead, most of the recent changes are incremental:

- improve the existing design
- add a dark/light mode toggle
- switch the accent color
- clean up headers and icons
- make backend pieces safer
- improve deploy behavior without redoing everything

I think that’s the right instinct.

A personal site does not need to be a monument to technical complexity. It needs to be clear, durable, and easy to update. If I’m going to keep writing, the publishing path needs to stay lightweight enough that posting feels easier than postponing.

That constraint matters. The best blog engine in the world is useless if it adds friction every time you want to hit publish.

## The invisible work is usually the work that matters

One reason I wanted to write this is that GitHub activity can be misleading if you only look for dramatic milestones.

A lot of the meaningful work in my repos lately has been invisible to anyone who only scans titles:

- adjusting deployment assumptions so preview and production behave properly
- tightening the site’s security posture instead of chasing quick shortcuts
- improving content publishing reliability
- adding backend structure around forms and resume delivery
- cleaning up the site’s presentation so it feels more finished

That kind of work doesn’t usually get applause. Nobody sees a cleanly handled environment variable and thinks, “wow.” But those are exactly the details that make software feel solid.

I’ve started to think that one sign of engineering maturity is being willing to spend real time on work that mostly prevents future annoyance.

Not because it’s exciting in the moment, but because you’re the one who will have to live with the platform later.

## A repo tells the truth about how you work

If I had to summarize what this site repo says about me right now, it’s probably this:

I care a lot less about novelty than I do about reliability and iteration.

The recent commit history is not the history of someone trying to show off. It’s the history of someone trying to make a system a little cleaner every time they touch it.

That includes writing posts, but it also includes all the surrounding maintenance work:

- getting build behavior right
- keeping secrets out of source
- reducing weird edge cases
- making the UI feel more intentional
- creating a setup I can continue to use without fighting it

Honestly, that’s how I approach most engineering work now. Build the thing, yes — but also improve the path around the thing. The path to deploy it. The path to debug it. The path to update it a month later when context is gone.

## What I want from this site going forward

I don’t want this site to become a portfolio fossil.

I want it to stay alive enough to reflect what I’m actually doing: platform engineering, infrastructure work, debugging, automation, observability, and all the weird practical lessons in between.

That probably means more posts that are grounded in real work and fewer posts that try to sound like announcements.

And it definitely means continuing to treat the site itself like a real system:

- simple enough to use
- structured enough to maintain
- secure enough to trust
- polished enough to feel deliberate

That’s the main lesson I got from looking back through the repo.

The value wasn’t in one giant change. It was in the accumulation of a lot of smaller decisions that made the site more real.

And, honestly, that’s true of most good engineering.
