# Capability assistant batch

The gooo-capability-assistant-batch command asks one natural-language
capability question across multiple .gooo declarations.

Run it with:

    go run ./cmd/gooo-capability-assistant-batch --query "What can this language do?" examples/capability-discovery.gooo examples/reverse-observation.gooo

Each result preserves the input path and a declaration-bound assistant. The
assistant overview, discovery status, plan, guide, and evidence digest remain
independent for each declaration, so one file cannot make another file appear
supported.

The command is read-only. It does not execute a provider, grant authorization,
mutate the catalog, infer missing evidence, or treat cache presence as semantic
evidence.
