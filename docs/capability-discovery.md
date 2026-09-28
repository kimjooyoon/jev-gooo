# Capability discovery

capability.Discover is a non-executing way to ask what a .gooo declaration can support.

It does four things:

1. Matches a natural-language query to a small, explicit capability catalog.
2. Reads only declared signals such as activity, evidence, reverse, generation, lsp, and feedback.
3. Returns AVAILABLE, DEFERRED, or UNKNOWN.
4. Preserves missing signals, next questions, the declaration digest, and an evidence digest.

It never executes code, infers authorization, or turns a cache hit into semantic evidence.

Example input for go run ./cmd/gooo-capability-discover:

    {
      "query": "What can this language do with provenance?",
      "declaration": "package jev\nactivity discover_capability\nproperty evidence_digest string"
    }

The output is deterministic for the same query and declaration. A DEFERRED result identifies the exact declaration signals still needed. An UNKNOWN result means the query did not match the catalog; it remains a useful boundary rather than a guessed answer.
