package capability

import "testing"

func TestDiscoverFromSourceBindsDeclaration(t *testing.T) {
	discovery, err := DiscoverFromSource(
		"What can this language do?",
		"package jev\nactivity discover_capability\nproperty evidence_digest string",
	)
	if err != nil {
		t.Fatal(err)
	}
	if discovery.Status != StatusAvailable {
		t.Fatalf("status = %s, want %s", discovery.Status, StatusAvailable)
	}
	if err := discovery.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestDiscoverFromSourceRejectsEmptyDeclaration(t *testing.T) {
	if _, err := DiscoverFromSource("What can this language do?", ""); err == nil {
		t.Fatal("empty declaration should be rejected")
	}
}
