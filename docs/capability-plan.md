# Capability plan

The plan command turns a declaration-bound discovery into explicit next steps without executing anything.

    printf '%s\n' '{"query":"Can this language run it?","declaration":"package jev\nactivity inspect"}' | go run ./cmd/gooo-capability-plan

UNKNOWN produces a clarifying-question step. DEFERRED produces one pending step for each missing declaration signal. AVAILABLE produces pending evidence and feedback-review steps. The plan is not an authorization decision and never upgrades a result on its own; its digest makes the proposed next steps inspectable.
