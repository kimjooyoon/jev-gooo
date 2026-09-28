# jev-gooo

Provider-neutral JEV execution envelopes and Gooo provenance experiments.

## Capability discovery

The repository also provides a non-executing way to inspect what a .gooo declaration can support.

Print the explicit capability catalog:

    go run ./cmd/gooo-capability-catalog

Ask a declaration-bound natural-language question:

    printf '%s
' '{"query":"What can this language do with provenance?","declaration":"package jev
activity discover_capability
property evidence_digest string"}' | go run ./cmd/gooo-capability-discover

Discovery returns AVAILABLE, DEFERRED, or UNKNOWN. It preserves the declaration digest, missing declaration signals, next questions, and an evidence digest.

Ask a broad “what can gooo do?” question and receive a bounded catalog overview:

    printf '%s
' '{"query":"What can this language do?","declaration":"package jev
activity discover_capability
property evidence_digest string"}' | go run ./cmd/gooo-capability-overview

The overview keeps broad discovery natural while labeling the result CATALOG_OVERVIEW, preserving declaration-bound capability states, suggested narrower questions, and the explicit constraint that the catalog is not a language-completeness claim.

List all matching capability options rather than only the top match:

    printf '%s
' '{"query":"What can this language do with provenance and code generation?","declaration":"package jev
activity reverse_observe
activity generate_output
property evidence_digest string
property output string"}' | go run ./cmd/gooo-capability-options

Measure the observed catalog surface without treating it as language completeness:

    printf '%s
' '{"query":"What can this language do with provenance and code generation?","declaration":"package jev
activity reverse_observe
property evidence_digest string"}' | go run ./cmd/gooo-capability-coverage

Bind coverage and explicit feedback into a review-safe handoff:

    printf '%s
' '{"query":"reverse observation provenance","declaration":"package jev
activity reverse_observe
property evidence_digest string","source":"verified receipt","evidence_digest":"0000000000000000000000000000000000000000000000000000000000000000","verified":true}' | go run ./cmd/gooo-capability-feedback-binding

Map the handoff to a bounded next investment action:

    printf '%s
' '{"query":"reverse observation provenance","declaration":"package jev
activity reverse_observe
property evidence_digest string","source":"verified receipt","evidence_digest":"0000000000000000000000000000000000000000000000000000000000000000","verified":true}' | go run ./cmd/gooo-capability-focus

Build a deterministic next-use guide from the discovery and plan:

    printf '%s
' '{"query":"What can this language do with provenance?","declaration":"package jev
activity reverse_observe
property evidence_digest string"}' | go run ./cmd/gooo-capability-guide

The guide turns UNKNOWN into clarification, DEFERRED into missing declaration signals, and AVAILABLE into verified-evidence collection. It never invokes a provider, grants authorization, mutates the catalog, or treats cache presence as semantic evidence.

Capability feedback is deliberately conservative. Missing or unverified evidence remains pending, and UNKNOWN or DEFERRED discovery is never promoted automatically. Only explicit verified evidence can become eligible for catalog review.

The envelope, discovery, feedback, guide, options, coverage, binding, and focus packages never invoke a provider, issue an authorization grant, or treat cache presence as semantic evidence.
