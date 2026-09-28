package capability

import "testing"

func TestBuildImprovementLedgerPreservesUnknownBoundary(t *testing.T) {
	coverage, feedback, plan := buildLedgerFixture(t, CoverageFeedbackPreserveUnknown)
	ledger, err := BuildImprovementLedger(coverage, feedback, plan)
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
	coverage, feedback, plan := buildLedgerFixture(t, CoverageFeedbackEligibleForCatalogReview)
	ledger, err := BuildImprovementLedger(coverage, feedback, plan)
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
	coverage, feedback, plan := buildLedgerFixture(t, CoverageFeedbackPreserveDeferred)
	plan.BindingDigest = "0000000000000000000000000000000000000000000000000000000000000000"
	if _, err := BuildImprovementLedger(coverage, feedback, plan); err == nil {
		t.Fatal("expected mismatched plan to fail")
	}
}

func TestImprovementLedgerRejectsCompletenessClaim(t *testing.T) {
	coverage, feedback, plan := buildLedgerFixture(t, CoverageFeedbackPreserveUnknown)
	ledger, err := BuildImprovementLedger(coverage, feedback, plan)
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
	options, err := DiscoverOptions("reverse observation provenance", declarationForLedger())
	if err != nil {
		t.Fatal(err)
	}
	coverage, err := MeasureCoverage(options)
	if err != nil {
		t.Fatal(err)
	}
	feedback, err := BuildFeedback(disposition)
	if err != nil {
		t.Fatal(err)
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
