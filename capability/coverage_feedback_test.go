package capability

import (
	"testing"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

func availableCoverageFeedbackInputs(t *testing.T) (Coverage, Feedback) {
	t.Helper()
	declaration, err := envelope.BindDeclaration("package demo\nactivity reverse_observe\nproperty evidence_digest string")
	if err != nil { t.Fatal(err) }
	query := "reverse observation provenance"
	discovery, err := Discover(query, declaration)
	if err != nil { t.Fatal(err) }
	options, err := DiscoverOptions(query, declaration)
	if err != nil { t.Fatal(err) }
	coverage, err := MeasureCoverage(options)
	if err != nil { t.Fatal(err) }
	feedback := ObserveFeedback(discovery, FeedbackInput{Source: "verified receipt", EvidenceDigest: digest("verified-receipt"), Verified: true})
	return coverage, feedback
}

func TestBindCoverageFeedbackEligibleRequiresBothEvidenceStreams(t *testing.T) {
	coverage, feedback := availableCoverageFeedbackInputs(t)
	binding, err := BindCoverageFeedback(coverage, feedback)
	if err != nil { t.Fatal(err) }
	if err := binding.Validate(); err != nil { t.Fatal(err) }
	if binding.Status != CoverageFeedbackEligibleForCatalogReview || binding.MissingStage != "" {
		t.Fatalf("unexpected eligible binding: %+v", binding)
	}
}

func TestBindCoverageFeedbackPreservesDeferredState(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo\nactivity observe_result")
	if err != nil { t.Fatal(err) }
	query := "reverse observation provenance"
	discovery, err := Discover(query, declaration)
	if err != nil { t.Fatal(err) }
	options, err := DiscoverOptions(query, declaration)
	if err != nil { t.Fatal(err) }
	coverage, err := MeasureCoverage(options)
	if err != nil { t.Fatal(err) }
	feedback := ObserveFeedback(discovery, FeedbackInput{Source: "unverified", EvidenceDigest: "", Verified: false})
	binding, err := BindCoverageFeedback(coverage, feedback)
	if err != nil { t.Fatal(err) }
	if err := binding.Validate(); err != nil { t.Fatal(err) }
	if binding.Status != CoverageFeedbackPreserveDeferred || binding.MissingStage == "" {
		t.Fatalf("deferred state was not preserved: %+v", binding)
	}
}

func TestBindCoverageFeedbackDoesNotPromoteUnknownCoverage(t *testing.T) {
	coverage, feedback := availableCoverageFeedbackInputs(t)
	declaration, err := envelope.BindDeclaration("package demo\nactivity reverse_observe\nproperty evidence_digest string")
	if err != nil { t.Fatal(err) }
	unknownOptions, err := DiscoverOptions("how can I model a compiler", declaration)
	if err != nil { t.Fatal(err) }
	coverage, err = MeasureCoverage(unknownOptions)
	if err != nil { t.Fatal(err) }
	binding, err := BindCoverageFeedback(coverage, feedback)
	if err != nil { t.Fatal(err) }
	if err := binding.Validate(); err != nil { t.Fatal(err) }
	if binding.Status != CoverageFeedbackPreserveUnknown || binding.MissingStage != "coverage_match" {
		t.Fatalf("unknown coverage was promoted: %+v", binding)
	}
}

func TestBindCoverageFeedbackRejectsTampering(t *testing.T) {
	coverage, feedback := availableCoverageFeedbackInputs(t)
	binding, err := BindCoverageFeedback(coverage, feedback)
	if err != nil { t.Fatal(err) }
	binding.Status = CoverageFeedbackEligibleForCatalogReview
	if err := binding.Validate(); err == nil {
		t.Fatal("expected tampered binding to fail validation")
	}
}