package envelope

import (
	"strings"
	"testing"
)

func TestCapabilityObservationEnvelopeReportPreservesAvailableProvenance(t *testing.T) {
	value, err := NewCapabilityObservationEnvelope(
		"corr-report-available",
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
	if report.State != CapabilityObservationAvailable || !report.Requestable {
		t.Fatalf("report state/requestable = %q/%t", report.State, report.Requestable)
	}
	if report.Provenance.DeclarationDigest != strings.Repeat("2", 64) || report.ObservationDigest != value.Digest {
		t.Fatalf("report lost provenance tuple")
	}
}

func TestCapabilityObservationEnvelopeReportPreservesDeferredBoundary(t *testing.T) {
	value, err := NewCapabilityObservationEnvelope(
		"corr-report-deferred",
		strings.Repeat("1", 64),
		strings.Repeat("2", 64),
		strings.Repeat("3", 64),
		strings.Repeat("4", 64),
		strings.Repeat("5", 64),
		[]CapabilityObservation{{ID: "capability-discovery", State: CapabilityObservationDeferred}},
		2,
		"provide_missing_declaration_signal",
	)
	if err != nil {
		t.Fatalf("new deferred envelope: %v", err)
	}
	report, err := value.Report()
	if err != nil {
		t.Fatalf("deferred report: %v", err)
	}
	if report.State != CapabilityObservationDeferred || report.Requestable {
		t.Fatalf("deferred report state/requestable = %q/%t", report.State, report.Requestable)
	}
	if report.FirstMissingStage != 2 || report.NextOperation == "" {
		t.Fatalf("deferred report lost boundary")
	}
}

func TestCapabilityObservationReportRejectsTamperedRequestability(t *testing.T) {
	value, err := NewCapabilityObservationEnvelope(
		"corr-report-tampered",
		strings.Repeat("1", 64),
		strings.Repeat("2", 64),
		strings.Repeat("3", 64),
		strings.Repeat("4", 64),
		strings.Repeat("5", 64),
		[]CapabilityObservation{{ID: "capability-discovery", State: CapabilityObservationUnknown}},
		0,
		"ask_narrower_question",
	)
	if err != nil {
		t.Fatalf("new unknown envelope: %v", err)
	}
	report, err := value.Report()
	if err != nil {
		t.Fatalf("unknown report: %v", err)
	}
	report.Requestable = true
	if err := report.Validate(); err == nil {
		t.Fatal("tampered requestability should be rejected")
	}
}
