# Capability catalog

The catalog command prints the explicit capability families understood by discovery.

    go run ./cmd/gooo-capability-catalog

The catalog is intentionally small and inspectable. It describes aliases and required declaration signals; it is not an authorization policy and it does not claim that a provider can execute anything. A query becomes AVAILABLE only when its catalog family and required signals are both present.
