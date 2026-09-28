# Capability observation envelope example

The `gooo-capability-observation` command connects declaration-bound discovery to
the JEV envelope without invoking a provider.

Input requires a query, declaration, correlation ID, catalog digest, and schema
digest. `first_missing_stage` is required for `DEFERRED` and `UNKNOWN` results.

The command emits three layers:

- the discovery result with its evidence digest;
- the observation envelope with source, declaration, catalog, schema, and evidence
  identities;
- an optional request-shaped value only when the observed capability is
  `AVAILABLE`.

The request-shaped value is still not an authorization grant. `DEFERRED` and
`UNKNOWN` remain non-requestable and preserve their next operation and missing
stage. A denominator, cache hit, or catalog match is not a language-completeness
claim.
