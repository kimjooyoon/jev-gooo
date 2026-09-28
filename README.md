# jev-gooo

Provider-neutral JEV execution envelopes and Gooo provenance experiments.

## Capability discovery

The repository also provides a non-executing way to inspect what a .gooo declaration can support.

Print the explicit capability catalog:

    go run ./cmd/gooo-capability-catalog

Ask a declaration-bound natural-language question:

    printf '%s\n' '{"query":"What can this language do with provenance?","declaration":"package jev\nactivity discover_capability\nproperty evidence_digest string"}' | go run ./cmd/gooo-capability-discover

Discovery returns AVAILABLE, DEFERRED, or UNKNOWN. It preserves the declaration digest, missing declaration signals, next questions, and an evidence digest.

Capability feedback is deliberately conservative. Missing or unverified evidence remains pending, and UNKNOWN or DEFERRED discovery is never promoted automatically. Only explicit verified evidence can become eligible for catalog review.

The envelope, discovery, and feedback packages never invoke a provider, issue an authorization grant, or treat cache presence as semantic evidence.
