# Capability guide

`gooo-capability-guide` turns declaration-bound discovery into a deterministic next-use route.

- `UNKNOWN` becomes a clarifying question.
- `DEFERRED` lists the missing declaration signals and their questions.
- `AVAILABLE` routes the user to verified evidence collection.

The guide carries the discovery digest, plan digest, and its own digest. It is advisory only: it does not invoke providers, grant authorization, mutate the capability catalog, or treat cache presence as semantic evidence.

Example:

    printf '%s\n' '{"query":"What can this language do with provenance?","declaration":"package jev\nactivity reverse_observe\nproperty evidence_digest string"}' | go run ./cmd/gooo-capability-guide