# Capability options

`gooo-capability-options` lists every catalog match for a declaration-bound natural-language query instead of returning only the top match.

Each option includes a deterministic score, required declaration signals, and `AVAILABLE` or `DEFERRED` status. An empty list is a valid evidence-bound result for an unknown query.

The command is advisory only. It does not execute a provider, grant authorization, mutate the catalog, or treat cache presence as semantic evidence.