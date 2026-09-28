package capability

import (
	"strings"
	"testing"
)

func TestAnalyzeSourceConnectsShapeDiscoveryAndPlan(t *testing.T) {
	analysis, err := AnalyzeSource(
		"What can this language do?",
		"package jev\nnamespace example\nentity Receipt\nproperty status string\nactivity discover_capability\nproperty evidence_digest string",
	)
	if err != nil {
		t.Fatal(err)
	}
	if analysis.Discovery.Status != StatusAvailable || analysis.Plan.Status != StatusAvailable {
		t.Fatalf("unexpected analysis: %#v", analysis)
	}
	if len(analysis.DeclarationShape.Entities) != 1 {
		t.Fatalf("unexpected declaration shape: %#v", analysis.DeclarationShape)
	}
	if err := analysis.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestAnalyzeSourcePreservesUnknownPlan(t *testing.T) {
	analysis, err := AnalyzeSource("How should I choose a color palette?", "package jev\nactivity inspect")
	if err != nil {
		t.Fatal(err)
	}
	if analysis.Discovery.Status != StatusUnknown || analysis.Plan.Status != StatusUnknown {
		t.Fatalf("unexpected analysis: %#v", analysis)
	}
	if err := analysis.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestAnalyzeSourceRejectsTampering(t *testing.T) {
	analysis, err := AnalyzeSource("What can this language do?", "package jev\nentity Evidence\nproperty evidence_digest string\nactivity discover_capability")
	if err != nil {
		t.Fatal(err)
	}
	analysis.AnalysisDigest = strings.Repeat("0", 64)
	if err := analysis.Validate(); err == nil {
		t.Fatal("tampered source analysis should be rejected")
	}
}
