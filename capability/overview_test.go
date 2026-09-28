package capability

import (
	"strings"
	"testing"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

func TestDiscoverOverviewExpandsBroadQuestion(t *testing.T) {
	declaration, err := envelope.BindDeclaration(`package jev
activity discover_capability
property evidence_digest string`)
	if err != nil {
		t.Fatal(err)
	}
	overview, err := DiscoverOverview("What can this language do?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	if err := overview.Validate(); err != nil {
		t.Fatal(err)
	}
	if overview.Mode != OverviewCatalog || len(overview.Capabilities) != len(catalog) {
		t.Fatalf("overview = %#v", overview)
	}
	if !strings.Contains(overview.Summary, "not a completeness claim") {
		t.Fatalf("summary = %q", overview.Summary)
	}
}

func TestDiscoverOverviewKeepsMatchedQuestionNarrow(t *testing.T) {
	declaration, err := envelope.BindDeclaration(`package jev
activity reverse_observe
property evidence_digest string`)
	if err != nil {
		t.Fatal(err)
	}
	overview, err := DiscoverOverview("What can this language do with provenance?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	if err := overview.Validate(); err != nil {
		t.Fatal(err)
	}
	if overview.Mode != OverviewMatchedOptions || len(overview.Capabilities) != 1 || overview.Capabilities[0].ID != "reverse-observation" {
		t.Fatalf("overview = %#v", overview)
	}
}

func TestDiscoverOverviewRequiresNarrowerQuestion(t *testing.T) {
	declaration, err := envelope.BindDeclaration(`package jev
activity unrelated`)
	if err != nil {
		t.Fatal(err)
	}
	overview, err := DiscoverOverview("Can this help with databases?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	if err := overview.Validate(); err != nil {
		t.Fatal(err)
	}
	if overview.Mode != OverviewClarification || len(overview.Capabilities) != 0 {
		t.Fatalf("overview = %#v", overview)
	}
}

func TestDiscoverOverviewRejectsTampering(t *testing.T) {
	declaration, err := envelope.BindDeclaration(`package jev
activity discover_capability
property evidence_digest string`)
	if err != nil {
		t.Fatal(err)
	}
	overview, err := DiscoverOverview("What can this language do?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	overview.Summary += " tampered"
	if err := overview.Validate(); err == nil {
		t.Fatal("expected tampered overview to fail validation")
	}
}
