---
title: "I Gave My Agent a Typed API Instead of a Tool List, and It Found a Hole in My DNS"
date: 2026-12-01T09:00:00-05:00
draft: false
summary: "Cloudflare's Code Mode replaces a list of MCP tools with one code-execution tool and a typed API. I pointed it at my own account. The full resolved spec is 3,764 operations and about 6.6 million tokens — and the thing that found a real misconfiguration was a for-loop, not the model."
tags: ["ai", "platform-engineering", "homelab"]
---

The Cloudflare API has **3,764 operations** across 2,353 paths. I counted them myself, in the thing I am writing about. Serialised with every `$ref` resolved, those operations come to 26.4 million characters — roughly **6.6 million tokens** of schema.

There is no context window that holds that. There probably never will be, because the API grows faster than the windows do.

So the interesting question is not "how do I fit the Cloudflare API into a model." It is "why was I ever trying to?"

## The pattern

The normal way to give a model an API is MCP: one tool per operation, each with a JSON schema, all of it loaded into the context before the model has done anything. It works until the catalog gets big. Then the tool definitions crowd out the actual work, and the model spends its context remembering what it *could* do instead of doing it.

Cloudflare's [Code Mode](https://developers.cloudflare.com/agents/tools/codemode/) inverts that. The model gets **two** tools instead of 3,764:

- `search` — run a JavaScript function over the OpenAPI spec and return what you find.
- `execute` — run a JavaScript function that calls `cloudflare.request()` and return the result.

The API surface is not in the context. It is a thing the code can query. Cloudflare's own figure for this: about **1,000 tokens** to describe the whole API that way, against roughly **1,170,000 tokens** for the equivalent native MCP server with full schemas. My 6.6 million is the same measurement taken differently — I serialised every operation in full, response schemas included, which their comparison does not. Either number makes the same point by a factor of a thousand.

## Counting the API with the API

The first thing I did was ask it about itself. That is one `search` call, and the function is the whole request:

```javascript
async () => {
  let paths = 0, ops = 0, chars = 0;
  for (const [path, methods] of Object.entries(spec.paths)) {
    paths++;
    for (const [method, op] of Object.entries(methods)) {
      ops++;
      chars += JSON.stringify(op).length;
    }
  }
  return { paths, ops, chars, approxTokens: Math.round(chars / 4) };
}
```

`{ paths: 2353, ops: 3764, chars: 26387182, approxTokens: 6596796 }`

Notice what did *not* happen. The spec never entered the model's context. A loop walked 26 megabytes of JSON and four numbers came back.

## The real task

Counting endpoints is a parlour trick. Here is the task I actually wanted: **for every zone on my account, how many DNS records are there, what types, and how many are proxiable but not proxied.**

That last clause is the one that matters. A record that *could* sit behind Cloudflare but does not is an origin server with its address published in public DNS.

With one tool per operation, that job is a minimum of four model round trips — list accounts, list zones, list records for zone one, list records for zone two — and every DNS record passes through the context on the way. With Code Mode it is one function:

```javascript
async () => {
  const accts = await cloudflare.request({ method: "GET", path: "/accounts" });
  const out = [];
  for (const a of accts.result) {
    const zones = await cloudflare.request({
      method: "GET", path: "/zones", query: { "account.id": a.id, per_page: 50 }
    });
    for (const z of zones.result) {
      const recs = await cloudflare.request({
        method: "GET", path: `/zones/${z.id}/dns_records`, query: { per_page: 500 }
      });
      const byType = {};
      let proxied = 0, exposed = 0;
      for (const r of recs.result) {
        byType[r.type] = (byType[r.type] || 0) + 1;
        if (r.proxiable) r.proxied ? proxied++ : exposed++;
      }
      out.push({ zone: z.name, records: recs.result.length, byType, proxied, exposed });
    }
  }
  return out;
}
```

Four API calls, 1.1 seconds, one model round trip:

| Zone | Records | Types | Proxied | Proxiable, not proxied |
|---|---|---|---|---|
| ark31.info | 21 | 1 A, 8 CNAME, 6 MX, 6 TXT | 2 | 7 |
| bp31app.com | 5 | 4 A, 1 CNAME | 2 | 3 |

Ten records that could be behind Cloudflare and are not. Seven of them are the CNAMEs that point `mail`, `calendar`, `drive` and `groups` at Google, which is correct and must stay that way — proxying someone else's hostname breaks it.

The other three were the finding. On `bp31app.com`, a **wildcard A record pointing straight at a public IP address**, unproxied. Every hostname on that domain that I had not explicitly configured — including the Kubernetes learning platform I have [written about](/posts/k8s-learn-try-it/) and which is currently returning 404 — resolves directly to an origin box, bypassing Cloudflare completely. No proxy, no WAF, no rate limiting, origin address in public DNS for anyone who asks.

I did not know that. I have been running this domain for months.

## What actually found it

This is the part I want to be precise about, because it is easy to tell this story as "the AI found a security hole in my DNS," and that is not what happened.

The model made one decision: *the interesting quantity is `proxiable && !proxied`.* That is judgment. It is a question about what matters, it has no correct answer derivable from the data, and framing it well is the entire value the model added.

Everything after that was a for-loop. Counting records, grouping by type, comparing two booleans, assembling a table — all of it has exactly one right answer, and none of it needed a language model. If I had asked a model to read 26 records and tell me which were unproxied, it would probably have got it right, and I would have had no way to know it had.

That is [the rule I keep relearning](/posts/alert-triage-decision-models/), from building an alert triage bot earlier this year: **use the model for judgment, use code for anything with a right answer.** I wrote that line about keeping deterministic computation out of an LLM's hands. Code Mode is the same rule pointed the other way — instead of computing the facts up front and handing the model a summary, you let the model write the computation and hand *itself* the summary. The boundary lands in the same place. The model never touches a number it could get wrong.

## Where it does not help

Code Mode is [marked experimental](https://developers.cloudflare.com/agents/tools/codemode/) in Cloudflare's own docs, and the honest case against it is short:

- **One tool, one call, no thanks.** If the job is "purge this cache," a direct tool call is one round trip. Code Mode makes it two — one to write the code, one to run it — and the code is a wrapper around a single request. Pure overhead.
- **Latency is unmeasured.** Cloudflare publishes the token numbers and does not publish latency or accuracy. Generating code and executing it in a sandbox is strictly more work per step than calling a tool. It pays off by removing *later* round trips, which means it pays off on multi-step work and loses on single-step work. I have no figures beyond that, and neither, publicly, do they.
- **A generated for-loop is still generated.** Deterministic code is not correct code. My loop was right because I read it. The failure mode moves from "the model miscounted" to "the model wrote a loop that confidently miscounts," which is quieter, and quiet failures are the ones that cost me twenty-four hours last time.
- **Blast radius.** `execute` holds a credential that can create and delete real infrastructure. I kept to `GET`. Credentials never enter the model context and outbound `fetch` is blocked by default, which is the right design — but a `DELETE` typed into that box is a `DELETE`.

## What I am doing about it

The wildcard is the hole, and proxying it is the lazy fix: it would push every hostname on the domain behind Cloudflare whether that is right for it or not. So the fix was the narrower one — `learn.bp31app.com` now has its own record, a CNAME to a Cloudflare Tunnel, so that hostname gets a deliberate route instead of inheriting the wildcard, and the origin address is no longer the default answer to every name I never configured.

It is a change I should have made in February. It took a for-loop over 26 DNS records to tell me it was needed.
