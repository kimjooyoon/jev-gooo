package envelope

import (
	"strings"
	"testing"
)

func TestCapabilityObservationInvestmentKeepsAvailableReviewBounded(t *testing.T) {
	value, err := NewCapabilityObservationEnvelope(
		"corr-investment-available",
		strings.Repeat("1", 64),
		strings.Repeat("2", 64),
		strings.Repeat("3", 64),
		strings.Repeat("4", 64),
		strings.Repeat("5", 64),
		[]CapabilityObservation{{ID: "capability-discovery", State: CapabilityObservationAvailable}},
		-1,
		"",
	)
	if err != nil {
		t.Fatalf("new envelope: %v", err)
	}
	report, err := value.Report()
	if err != nil {
		t.Fatalf("report: %v", err)
	}
	investment, err := report.BuildInvestment()
	if err != nil {
		t.Fatalf("build investment: %v", err)
	}
	if investment.Decision != "CATALOG_REVIEW_READY" || investment.NextAction != "review_candidate_without_catalog_mutation" || investment.CompletenessClaimed {
		t.Fatalf("available investment = %#v", investment)
	}
}

func TestCapabilityObservationInvestmentPreservesUnknownNoInvestment(t *testing.T) {
	value, err := NewCapabilityObservationEnvelope(
		"corr-investment-unknown",
		strings.Repeat("1", 64),
		strings.Repeat("2", 64),
		strings.Repeat("3", 64),
		strings.Repeat("4", 64),
		strings.Repeat("5", 64),
		[]CapabilityObservation{{ID: "capability.discovery", State: CapabilityObservationUnknown}},
		0,
		"ask_narrower_question",
	)
	if err != nil {
		t.Fatalf("new envelope: %v", err)
	}
	report, err := value.Report()
	if err != nil {
		t.Fatalf("report: %v", err)
	}
	investment, err := report.BuildInvestment()
	if err != nil {
		t.Fatalf("build unknown investment: %v", err)
	}
	if investment.Decision != "DISCOVERY_REQUIRED" || investment.InvestmentClass != "NO_INVESTMENT_UNTIL_EVIDENCE" {
		t.Fatalf("unknown investment = %#v", investment)
	}
}

func TestCapabilityObservationInvestmentRejectsTamperedDigest(t *testing.T) {
	value, err := NewCapabilityObservationEnvelope(
		"corr-investment-tampered",
		strings.Repeat("1", 64),
		strings.Repeat("2", 64),
		strings.Repeat("3", 64),
		strings.Repeat("4", 64),
		strings.Repeat("5", 64),
		[]CapabilityObservation{{ID: "capability-discovery", State: CapabilityObservationDeferred}},
		1,
		"provide_missing_declaration_signal",
	)
	if err != nil {
		t.Fatalf("new envelope: %v", err)
	}
	report, err := value.Report()
	if err != nil {
		t.Fatalf("report: %v", err)
	}
	investment, err := report.BuildInvestment()
	if err != nil {
		t.Fatalf("build investment: %v", err)
	}
	investment.PlanDigest = strings.Repeat("0", 64)
	if err := investment.Validate(); err == nil {
		t.Fatal("tampered investment digest should be rejected")
	}
}
