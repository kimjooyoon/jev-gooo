package capability

import (
	"strings"
	"testing"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

func TestBuildImprovementLedgerPreservesUnknownBoundary(t *testing.T) {
	coverage, binding, plan := buildLedgerFixture(t, DispositionPreserveUnknown)
	ledger, err := BuildImprovementLedger(coverage, binding, plan)
	if err != nil {
		t.Fatal(err)
	}
	if err := ledger.Validate(); err != nil {
		t.Fatal(err)
	}
	if ledger.InvestmentClass != InvestmentNoInvestmentUntilEvidence || ledger.CompletenessClaimed {
		t.Fatalf("ledger = %#v", ledger)
	}
}

func TestBuildImprovementLedgerPreservesReviewBoundary(t *testing.T) {
	coverage, binding, plan := buildLedgerFixture(t, DispositionEligibleForCatalogReview)
	ledger, err := BuildImprovementLedger(coverage, binding, plan)
	if err != nil {
		t.Fatal(err)
	}
	if err := ledger.Validate(); err != nil {
		t.Fatal(err)
	}
	if ledger.InvestmentClass != InvestmentReviewBeforeCatalogChange || ledger.NextAction != "review_candidate_without_catalog_mutation" {
		t.Fatalf("ledger = %#v", ledger)
	}
}

func TestBuildImprovementLedgerRejectsMismatchedPlan(t *testing.T) {
	coverage, binding, plan := buildLedgerFixture(t, DispositionPreserveDeferred)
	plan.BindingDigest = "0000000000000000000000000000000000000000000000000000000000000000"
	if _, err := BuildImprovementLedger(coverage, binding, plan); err == nil {
		t.Fatal("expected mismatched plan to fail")
	}
}

func TestImprovementLedgerRejectsCompletenessClaim(t *testing.T) {
	coverage, binding, plan := buildLedgerFixture(t, DispositionPreserveUnknown)
	ledger, err := BuildImprovementLedger(coverage, binding, plan)
	if err != nil {
		t.Fatal(err)
	}
	ledger.CompletenessClaimed = true
	if err := ledger.Validate(); err == nil {
		t.Fatal("expected completeness claim to fail")
	}
}

func buildLedgerFixture(t *testing.T, disposition FeedbackDisposition) (Coverage, CoverageFeedback, FocusPlan) {
	t.Helper()
	lineBreak := string([]byte{10})
	query := "reverse observation provenance"
	source := strings.Join([]string{"package jev", "activity reverse_observe", "property evidence_digest string"}, lineBreak)
	switch disposition {
	case DispositionPreserveUnknown:
		query = "database migration"
		source = "package jev"
	case DispositionPreserveDeferred:
		source = strings.Join([]string{"package jev", "activity reverse_observe"}, lineBreak)
	}
	declaration, err := envelope.BindDeclaration(source)
	if err != nil {
		t.Fatal(err)
	}
	discovery, err := Discover(query, declaration)
	if err != nil {
		t.Fatal(err)
	}
	options, err := DiscoverOptions(query, declaration)
	if err != nil {
		t.Fatal(err)
	}
	coverage, err := MeasureCoverage(options)
	if err != nil {
		t.Fatal(err)
	}
	input := FeedbackInput{
		Source:         "verified receipt",
		EvidenceDigest: strings.Repeat("0", 64),
		Verified:       disposition == DispositionEligibleForCatalogReview,
	}
	feedback := ObserveFeedback(discovery, input)
	if feedback.Disposition != disposition {
		t.Fatalf("feedback disposition = %q, want %q", feedback.Disposition, disposition)
	}
	binding, err := BindCoverageFeedback(coverage, feedback)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := BuildFocusPlan(coverage, binding)
	if err != nil {
		t.Fatal(err)
	}
	return coverage, binding, plan
}
