# JEV + Gooo envelope

This repository is a provider-neutral experiment for binding a `.gooo` declaration to an observed execution envelope.

The `envelope` package only records evidence. It never invokes a provider, issues a capability grant, or treats a JEV decision as authorization.

The receipt remains `UNKNOWN` with its first `missing_stage` until the declaration, capability grant, workload identity, non-authorizing decision receipt, and terminal result are all present. The `.gooo` contract in `contracts/jev_execution_envelope.gooo` is the source-level identity for that boundary.
