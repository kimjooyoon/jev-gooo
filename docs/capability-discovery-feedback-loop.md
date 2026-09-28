# Capability discovery feedback loop

The feedback receipt closes the loop between a user question and the next language investment.

## Flow

1. gooo inspects a declaration and proposes focused questions.
2. The user selects a question and records whether it was useful, not useful, or unresolved.
3. The feedback receipt preserves source, declaration, capability receipt, evidence, and query identities.
4. jev replay compares the receipt with later observations without executing or authorizing the proposed operation.
5. meta-ontology uses the result to decide whether a missing stage deserves investment.

## Invariants

- The original question is never replaced by the suggested question.
- Missing evidence remains UNKNOWN or unresolved rather than becoming success.
- A feedback receipt must remain non-executing and non-authorizing.
- A digest is required before the feedback can be compared across observations.
- A useful result is evidence about user value, not proof of language completeness.

## Investment signal

A repeated useful outcome can justify measuring an implementation investment for the related missing stage. A not-useful outcome can reduce priority or change the question catalog. An unresolved outcome preserves the boundary and requests stronger evidence. All three outcomes are valuable because they change what the language should inspect next without pretending to know more than the provenance supports.