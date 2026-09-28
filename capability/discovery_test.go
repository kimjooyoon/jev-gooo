package capability

import (
	"strings"
	"testing"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

func testDeclaration(t *testing.T, source string) envelope.Declaration {
	t.Helper()
	declaration, err := envelope.BindDeclaration(source)
	if err != nil {
		t.Fatal(err)
	}
	return declaration
}

func TestDiscoverCapabilityAvailable(t *testing.T) {
	declaration := testDeclaration(t, "package jev\nactivity discover_capability\nproperty evidence_digest string")
	discovery, err := Discover("What can this language do with provenance?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	if discovery.Status != StatusAvailable || discovery.CapabilityID != "capability-discovery" {
		t.Fatalf("unexpected discovery: %#v", discovery)
	}
	if err := discovery.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestDiscoverCapabilityDeferred(t *testing.T) {
	declaration := testDeclaration(t, "package jev\nactivity inspect")
	discovery, err := Discover("Can this language run it?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	if discovery.Status != StatusDeferred || discovery.CapabilityID != "execution-observation" {
		t.Fatalf("unexpected discovery: %#v", discovery)
	}
	if len(discovery.MissingSignals) != 2 {
		t.Fatalf("missing signals = %#v", discovery.MissingSignals)
	}
	if err := discovery.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestDiscoverCapabilityUnknown(t *testing.T) {
	declaration := testDeclaration(t, "package jev\nactivity inspect")
	discovery, err := Discover("How should I choose a color palette?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	if discovery.Status != StatusUnknown || discovery.CapabilityID != "" {
		t.Fatalf("unexpected discovery: %#v", discovery)
	}
	if err := discovery.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestDiscoverCapabilityRejectsTampering(t *testing.T) {
	declaration := testDeclaration(t, "package jev\nactivity discover_capability\nproperty evidence_digest string")
	discovery, err := Discover("What can gooo do?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	discovery.EvidenceDigest = strings.Repeat("0", 64)
	if err := discovery.Validate(); err == nil {
		t.Fatal("tampered capability discovery should be rejected")
	}
}
