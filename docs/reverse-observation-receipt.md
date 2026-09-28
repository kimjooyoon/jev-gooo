# Reverse observation receipt

`capability.ObserveReverseObservation` binds a generated artifact back to its source, `.gooo` declaration, IR, and immutable evidence prefix.

- `OBSERVED` means the reverse-observed artifact digest equals the generated artifact digest.
- `DEFERRED` means the generation chain is bound but reverse observation is missing.
- `MISMATCH` means reverse observation exists but disagrees with generation.
- `UNKNOWN` preserves the first missing provenance stage.

The receipt is read-only. It does not execute generated code, invoke a provider, grant authorization, or mutate the capability catalog.
