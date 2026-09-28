# Discover what gooo can do

The capability discovery command turns a natural-language question into a
verified catalog overview. It reports which catalog entries are available for
the .gooo declaration, which are deferred, and which declaration or evidence
signals should be supplied next.

Run it with:

    go run ./cmd/gooo-capability-discovery --query "gooo가 지원하는 기능은?"

Use a narrower question to inspect one capability family:

    go run ./cmd/gooo-capability-discovery --query "How can gooo trace provenance?"

The JSON output is bounded by the existing Overview contract. Its evidence
digest is validated before output, and the result does not claim language
completeness, invoke a provider, grant authorization, mutate the catalog, or
treat cache presence as semantic evidence.
