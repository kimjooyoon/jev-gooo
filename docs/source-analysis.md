# Source analysis

The source analysis command connects the language-facing path in one non-executing response:

    printf '%s\n' '{"query":"What can this language do?","declaration":"package jev\nactivity discover_capability\nproperty evidence_digest string"}' | go run ./cmd/gooo-analyze

The response contains the declaration shape, capability discovery, and next-step plan. Each stage keeps its own digest and the aggregate adds an analysis digest. UNKNOWN and DEFERRED remain explicit states; the aggregate never executes a provider or applies a catalog change.
