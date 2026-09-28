# Capability improvement ledger

The improvement ledger binds capability coverage, explicit feedback, and a focus plan into a reviewable next-investment record.

It preserves:

- coverage, binding, plan, source, and optional evidence digests
- observed matched, available, and deferred counts
- feedback status, focus decision, and next action
- an investment class such as NO_INVESTMENT_UNTIL_EVIDENCE or REVIEW_BEFORE_CATALOG_CHANGE

It never produces a language-completeness score. completeness_claimed is always false, catalog coverage is explicitly scoped, cache presence is not semantic evidence, and the ledger cannot mutate the catalog, invoke a provider, or grant authorization.

The CLI expects JSON containing validated coverage, binding, and plan objects and returns a digest-bound ledger.
