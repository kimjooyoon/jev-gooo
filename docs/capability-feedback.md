# Capability feedback

Capability feedback is the first explicit self-improvement boundary for discovery.

ObserveFeedback does not turn a user note into a language feature. It preserves the current discovery state and classifies the feedback:

- PENDING_REQUIRED_EVIDENCE: source, evidence digest, or explicit verification is missing.
- PRESERVE_UNKNOWN: the query did not match the catalog.
- PRESERVE_DEFERRED: declaration signals are still missing.
- REQUIRES_DECLARATION_BINDING: the discovery result is invalid or unbound.
- ELIGIBLE_FOR_CATALOG_REVIEW: only an available discovery with explicit verified evidence reaches catalog review.

Every result carries a source digest, optional evidence digest, missing stage, and feedback digest. This keeps a metric from becoming a guessed correctness score.
