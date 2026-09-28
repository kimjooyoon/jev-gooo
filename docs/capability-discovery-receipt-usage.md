# Capability Discovery Receipt Usage

`jev-gooo` stores the result of a bounded capability question as a replayable receipt. The receipt describes an observation; it does not execute a provider or authorize a side effect.

## Flow

1. `gooo-lsp` answers a question such as `What can gooo do with this declaration?` and returns the capability status plus provenance digests.
2. The consumer copies the source, declaration, query, catalog, evidence, and toolchain identities into a `capability.discovery.v1` receipt.
3. `jev-capability-receipt` validates the receipt before it is stored or sent to a later observation step.
4. A later observation reuses the prior evidence digest and can confirm, refute, or leave the capability unresolved.

## Status interpretation

- `AVAILABLE`: the declaration-bound catalog supports a capability projection. This is not an execution grant.
- `DEFERRED`: an explicit external boundary or additional evidence is required.
- `UNKNOWN`: the answer is not established. Keep `first_missing_stage` and `next_operation` instead of inferring success.

## Validation example

```text
go run ./cmd/jev-capability-receipt examples/capability_discovery_available.json
```

The validator must reject execution or authorization flags, malformed digests, duplicate capability IDs, and state records that omit their required boundary fields.

## Provenance rule

A receipt is useful only when its source, declaration, query, catalog, evidence, and toolchain identities are preserved. Cache presence, capability count, and a successful editor response are not semantic evidence of domain completeness.