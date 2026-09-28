package capability

import (
	"testing"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

func TestBuildFocusPlanKeepsCatalogReviewBounded(t *testing.T) {
	coverage, feedback := availableCoverageFeedbackInputs(t)
	binding, err := BindCoverageFeedback(coverage, feedback)
	if err != nil { t.Fatal(err) }
	plan, err := BuildFocusPlan(coverage, binding)
	if err != nil { t.Fatal(err) }
	if err := plan.Validate(); err != nil { t.Fatal(err) }
	if plan.Decision != FocusCatalogReviewReady || plan.CompletenessClaimed || plan.NextAction != "review_candidate_without_catalog_mutation" {
		t.Fatalf("unexpected bounded focus plan: %+v", plan)
	}
}

func TestBuildFocusPlanRoutesUnknownToDiscovery(t *testing.T) {
	coverage, feedback := availableCoverageFeedbackInputs(t)
	declaration, err := envelope.BindDeclaration("package demo\nactivity reverse_observe\nproperty evidence_digest string")
	if err != nil { t.Fatal(err) }
	unknownOptions, err := DiscoverOptions("how can I model a compiler", declaration)
	if err != nil { t.Fatal(err) }
	coverage, err = MeasureCoverage(unknownOptions)
	if err != nil { t.Fatal(err) }
	binding, err := BindCoverageFeedback(coverage, feedback)
	if err != nil { t.Fatal(err) }
	plan, err := BuildFocusPlan(coverage, binding)
	if err != nil { t.Fatal(err) }
	if err := plan.Validate(); err != nil { t.Fatal(err) }
	if plan.Decision != FocusDiscoveryRequired || plan.NextAction != "clarify_or_extend_query" {
		t.Fatalf("unknown query did not route to discovery: %+v", plan)
	}
}

func TestBuildFocusPlanRoutesDeferredToDeclarationEvidence(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo\nactivity observe_result")
	if err != nil { t.Fatal(err) }
	query := "reverse observation provenance"
	discovery, err := Discover(query, declaration)
	if err != nil { t.Fatal(err) }
	options, err := DiscoverOptions(query, declaration)
	if err != nil { t.Fatal(err) }
	coverage, err := MeasureCoverage(options)
	if err != nil { t.Fatal(err) }
	feedback := ObserveFeedback(discovery, FeedbackInput{Source: "unverified", Verified: false})
	binding, err := BindCoverageFeedback(coverage, feedback)
	if err != nil { t.Fatal(err) }
	plan, err := BuildFocusPlan(coverage, binding)
	if err != nil { t.Fatal(err) }
	if err := plan.Validate(); err != nil { t.Fatal(err) }
	if plan.Decision != FocusDeclarationEvidenceRequired || plan.NextAction != "declare_missing_signals" {
		t.Fatalf("deferred query did not route to declaration evidence: %+v", plan)
	}
}

func TestBuildFocusPlanRejectsTampering(t *testing.T) {
	coverage, feedback := availableCoverageFeedbackInputs(t)
	binding, err := BindCoverageFeedback(coverage, feedback)
	if err != nil { t.Fatal(err) }
	plan, err := BuildFocusPlan(coverage, binding)
	if err != nil { t.Fatal(err) }
	plan.CompletenessClaimed = true
	if err := plan.Validate(); err == nil { t.Fatal("expected completeness overclaim to fail validation") }
}