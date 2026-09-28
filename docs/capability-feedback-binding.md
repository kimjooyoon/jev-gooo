# Coverage feedback binding

`gooo-capability-feedback-binding` joins the observed capability surface with explicit feedback evidence.

The binding remains `PRESERVE_UNKNOWN`, `PRESERVE_DEFERRED`, or `PENDING_REQUIRED_EVIDENCE` until both streams are valid. Even `ELIGIBLE_FOR_CATALOG_REVIEW` is only a review handoff; it never mutates the catalog, invokes a provider, or grants authorization.