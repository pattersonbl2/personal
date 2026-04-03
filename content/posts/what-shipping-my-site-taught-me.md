---
title: "What Shipping My Personal Site Actually Taught Me"
date: 2026-04-03
draft: false
---

When I looked back through my recent GitHub activity, I expected to find a clean story.

Something like: I redesigned my site, published a few posts, polished the homepage, and called it progress.

That is not what the commit history says.

What it actually shows is a lot messier and a lot more honest. Fixes for weird deployment behavior. Small security improvements. Frontmatter date issues. Backend cleanup. Theme tweaks. Resume token problems caused by whitespace. Little things that took time, mattered a lot, and would be almost invisible to anyone just looking at the finished site.

And honestly, that felt familiar.

That’s also what a lot of real engineering work looks like.

## I wanted a website. I ended up with a system.

From the outside, a personal site seems simple. A homepage, a blog, a resume, a contact form. Nothing too serious.

But once I owned the whole thing, it stopped feeling like a “website” and started feeling like a small platform I was responsible for.

Not in an overdramatic way. Just in the practical sense that everything had to line up:

- content had to build correctly
- preview and production deploys had to behave differently in the right places
- the contact flow had to be safe
- secrets had to stay out of the repo
- metadata had to be correct
- the design had to feel intentional enough that I’d actually want to keep using it

That’s the part people don’t really talk about when they talk about personal sites. The visible part is the writing. The invisible part is all the maintenance required to make publishing feel easy instead of annoying.

I think that invisible part is what I’ve actually been working on.

## A lot of my recent work was basically me removing future headaches

Looking through the repo, I noticed that most of my changes were not ambitious. They were preventive.

That might be the most accurate summary of how I work these days.

I spent time fixing post dates because Hugo was quietly treating some content as future-dated. I fixed token handling because a small formatting issue could turn a valid resume link into a broken one. I made deployment-related fixes because the difference between “works locally” and “works when published” is still where a lot of bugs like to hide.

None of that is exciting in the usual blog-post sense.

But I’ve started to trust that kind of work more.

There’s a phase of learning where big visible progress feels like the only kind that counts. New feature. New stack. New service. New post. Something you can point to.

Lately, I’ve found myself caring more about whether a thing is stable, understandable, and pleasant to come back to later.

That’s not as flashy, but it’s probably a better long-term instinct.

## This repo felt weirdly personal in a way I didn’t expect

What surprised me most was that the site repo didn’t just show what I built. It showed how I think.

The changes tell on me a little.

They show that I’m the kind of engineer who will spend time on details that most people never notice if I think those details will reduce friction later. They show that I care about reliability more than novelty. They show that I’d rather make the existing thing better than constantly rebuild it from scratch.

That last part matters to me.

I could absolutely turn my personal site into a permanent lab experiment. New framework every month. New theme every quarter. Constant reinvention disguised as ambition.

But I know myself well enough to know that too much churn kills momentum. If the publishing path becomes a project every time I want to write, I’ll write less. If the system becomes too clever, I’ll trust it less.

So a lot of the work in this repo is me trying to build something I’ll actually continue using.

Not the most impressive version. The most sustainable one.

## Tiny bugs are humbling, and I kind of appreciate that

One thing I’ve learned over and over, both in infrastructure work and in smaller personal projects, is that systems rarely fail in cinematic ways.

Usually it’s something dumb.

A date format. Whitespace in a token. A deployment assumption that turns out not to be true. A configuration mismatch between environments. A path that works in one context and breaks in another.

That used to frustrate me more than it does now.

Now I mostly see it as part of the craft. Software has a way of forcing you to respect details. It doesn’t really care whether the bug feels important enough to deserve your time. If it breaks the flow, it breaks the flow.

In a weird way, I think personal projects are good for learning that lesson because there’s nowhere to hide. No separate team, no handoff, no ticket disappearing into a queue. If something is brittle, I’m the one who feels it later.

That kind of feedback loop is annoying, but it’s useful.

## I’m trying to build fewer things I have to fight

I think that’s the real thread running through this repo.

Not “I built a personal website.”

More like: I’m trying to build a version of my tools, systems, and projects that I don’t have to fight every time I come back to them.

That applies to infrastructure. It applies to Kubernetes. It applies to observability work. And apparently, it also applies to my blog.

I want the site to feel like an extension of how I already work:

- simple where it should be simple
- structured where it needs structure
- secure enough that I’m not nervous about it
- lightweight enough that writing still feels easy

That probably sounds obvious, but I don’t think it is. A lot of us who like building things also like making things more elaborate than they need to be.

I definitely have that tendency.

So in some ways, this repo is also me practicing restraint.

## The finished site hides the part that probably matters most

If someone lands on the site, they’ll see posts, a homepage, maybe a resume, maybe a contact page.

What they won’t see is the accumulation of small decisions underneath it:

- the deployment fixes
- the content publishing corrections
- the security hardening
- the backend cleanup
- the UI polish
- the repeated effort to make the system a little less fragile than it was before

But that hidden layer is probably the part I’m proudest of.

Not because it’s technically advanced, but because it reflects a version of engineering I’ve come to value more over time: the kind where you make things calmer, clearer, and easier to trust.

That work doesn’t always produce a dramatic before-and-after.

Sometimes it just means the site keeps working. Publishing feels smoother. A future bug never happens. A weird edge case stops stealing your time.

That counts.

## What I want this site to be

I don’t want this site to be a frozen portfolio.

I want it to feel alive enough that when someone reads it, they’re getting a real snapshot of what I’m learning, building, debugging, and thinking about.

Not just the polished wins. The process too.

That probably means I want the writing to feel a little more personal going forward. Less like I’m presenting finished work from a distance, and more like I’m documenting what it actually feels like to build things, maintain them, get them wrong, fix them, and slowly get better.

Looking through this repo reminded me that progress usually doesn’t happen as one big leap. At least not for me.

It happens in small corrections. Small cleanups. Small improvements that make the next step easier.

That’s true of this site.

And, if I’m being honest, it’s true of me too.
