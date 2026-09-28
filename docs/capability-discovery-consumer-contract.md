# Capability discovery consumer contract

This package consumes a gooo capability discovery result as an observation, not as an execution request.

## Handoff

The producer is the gooo capability discovery command or LSP path. It may answer a natural question such as:

    What can gooo do with this declaration?

The consumer is the jev-capability-from-gooo command. The handoff is valid only when the producer output is bound to the exact source, declaration, catalog, and toolchain identities supplied to the consumer.

## What the receipt means

A capability discovery receipt records what can be inspected next:

- AVAILABLE means the capability description is available for inspection;
- DEFERRED means a declared stage exists but its evidence is not complete;
- UNKNOWN means the result cannot be safely classified;
- first_missing_stage identifies the earliest unresolved boundary;
- next_operation identifies an inspectable follow-up when one is justified;
- non_executing and non_authorizing must remain true.

The receipt does not execute a program, grant a capability, approve a workflow, or claim domain completeness.

## Required consumer behavior

A consumer must:

1. reject an unbound declaration or source;
2. reject a missing catalog or toolchain identity;
3. reject a receipt that has no justified next operation when the workflow requires one;
4. preserve the producer query, evidence digest, state, and first missing stage;
5. expose unresolved stages instead of converting them into success;
6. emit a canonical receipt digest for later replay.

The replay path can then compare prior evidence with a later observation and classify the result as confirmed, refuted, or unresolved.

## Why this improves gooo usability

The user can begin with a declaration and a question before a full implementation exists. gooo discovers a bounded next operation; jev makes the observation durable and replayable; meta-ontology measures whether the domain boundary is becoming more useful without hiding what remains unknown.
