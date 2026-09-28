# Capability assistant

The gooo-capability-assistant command combines the existing overview,
discovery, plan, and guide stages for one natural-language question and one
.gooo declaration.

For direct use with a declaration file:

    go run ./cmd/gooo-capability-assistant --query "What can this language do?" --declaration examples/capability-discovery.gooo

The JSON input form remains available when a caller already has an envelope:

    printf '%s
' '{"query":"What can this language do with provenance?","declaration":"package jev
activity reverse_observe
property evidence_digest string"}' | go run ./cmd/gooo-capability-assistant

The JSON result preserves:

- the bounded catalog overview and its mode
- the declaration-bound AVAILABLE, DEFERRED, or UNKNOWN discovery
- observed and pending plan steps
- the next action and questions from the guide
- shared declaration, stage, and final evidence digests

The assistant is a composition boundary, not a new completeness oracle. It does
not execute a provider, grant authorization, mutate the catalog, infer missing
evidence, or treat cache presence as semantic evidence. Unknown questions remain
unknown and retain a clarification question.
