# Convert gooo capability discovery to a JEV receipt

`jev-capability-from-gooo` consumes the JSON document emitted by the gooo
capability discovery command and produces a declaration-bound
`CapabilityDiscoveryReceipt`.

The converter requires exact values for the source artifact, capability
catalog, and toolchain. Pass them explicitly with `-source-digest`,
`-catalog-digest`, and `-toolchain-identity`. If the input has no bound
 declaration, or no next non-executing operation, conversion is rejected rather
than inferred. Use `-next-operation` only when the missing operation is known
from an external observation.

The converter preserves the discovery status, query/evidence digests,
capability identifiers, first missing stage, and the non-executing and
non-authorizing boundary. It does not invoke a provider, execute a plan, or
create an authorization grant.

Example CI pipeline shape:

```text
gooo-capabilities -json declaration.gooo > capability-discovery.json
jev-capability-from-gooo \
  -input capability-discovery.json \
  -source-digest <source-digest> \
  -catalog-digest <catalog-digest> \
  -toolchain-identity <toolchain-identity> \
  > capability-discovery-receipt.json
jev-capability-receipt capability-discovery-receipt.json
```

A later observation can be linked with `jev-capability-replay`; a receipt that
is `UNKNOWN` or `DEFERRED` remains an unresolved observation until new evidence
supports a transition.