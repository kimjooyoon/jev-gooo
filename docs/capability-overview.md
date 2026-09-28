# Capability overview

gooo-capability-overview answers broad questions such as “What can this language do?” without pretending that a catalog is the whole language.

The result has three deterministic modes:

- CATALOG_OVERVIEW: no narrower catalog match was found, so the current catalog is shown with declaration-bound AVAILABLE or DEFERRED states.
- MATCHED_OPTIONS: a narrower query matched one or more catalog entries.
- CLARIFICATION_REQUIRED: no explicit entry matched, so the caller receives a bounded next question.

Every response carries the declaration digest, capability states, suggested questions, fixed safety constraints, and an evidence digest. The feature does not execute providers, grant authorization, mutate the catalog, or use cache presence as semantic evidence. CATALOG_OVERVIEW is explicitly not a language-completeness claim.

Example:

    printf '%s
' '{"query":"What can this language do?","declaration":"package jev
activity discover_capability
property evidence_digest string"}' | go run ./cmd/gooo-capability-overview
