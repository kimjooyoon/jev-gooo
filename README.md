# jev-gooo

Provider-neutral JEV execution envelopes and Gooo provenance experiments.

## Capability discovery

The repository also provides a non-executing way to inspect what a .gooo declaration can support.

Print the explicit capability catalog:

    go run ./cmd/gooo-capability-catalog

Ask a declaration-bound natural-language question:

    printf '%s\n' '{"query":"What can this language do with provenance?","declaration":"package jev\nactivity discover_capability\nproperty evidence_digest string"}' | go run ./cmd/gooo-capability-discover

Discovery returns AVAILABLE, DEFERRED, or UNKNOWN. It preserves the declaration digest, missing declaration signals, next questions, and an evidence digest.

List all matching capability options rather than only the top match:

    printf '%s\n' '{"query":"What can this language do with provenance and code generation?","declaration":"package jev\nactivity reverse_observe\nactivity generate_output\nproperty evidence_digest string\nproperty output string"}' | go run ./cmd/gooo-capability-options

Build a deterministic next-use guide from the discovery and plan:

    printf '%s\n' '{"query":"What can this language do with provenance?","declaration":"package jev\nactivity reverse_observe\nproperty evidence_digest string"}' | go run ./cmd/gooo-capability-guide

The guide turns UNKNOWN into clarification, DEFERRED into missing declaration signals, and AVAILABLE into verified-evidence collection. It never invokes a provider, grants authorization, mutates the catalog, or treats cache presence as semantic evidence.

Capability feedback is deliberately conservative. Missing or unverified evidence remains pending, and UNKNOWN or DEFERRED discovery is never promoted automatically. Only explicit verified evidence can become eligible for catalog review.

The envelope, discovery, feedback, guide, and options packages never invoke a provider, issue an authorization grant, or treat cache presence as semantic evidence.