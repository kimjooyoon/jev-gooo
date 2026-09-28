# Capability observation report

The gooo-capability-report command summarizes a capability-assistant question
across multiple .gooo declarations.

Run it with:

    go run ./cmd/gooo-capability-report --query "What can this language do with provenance?" examples/capability-discovery.gooo examples/reverse-observation.gooo

The report preserves each input path and assistant evidence digest, then records
only the observed counts of AVAILABLE, DEFERRED, and UNKNOWN states. These counts
describe the inspected files; they are not a language completeness score and
must not be used as one.

The report rejects duplicate paths, mismatched queries, invalid assistant
evidence, and tampered aggregate digests. It does not execute a provider, grant
authorization, mutate the catalog, infer missing evidence, or treat cache
presence as semantic evidence.
