package envelope

import (
	"strings"
	"testing"
)

func TestCapabilityRequestForPreservesDeclarationBinding(t *testing.T) {
	digestValue := strings.Repeat("a", 64)
	value, err := NewCapabilityObservationEnvelope(
		"corr-3", digestValue, digestValue, digestValue, digestValue, digestValue,
		[]CapabilityObservation{{ID: "capability.query", State: CapabilityObservationAvailable}},
		-1, "",
	)
	if err != nil {
		t.Fatalf("create envelope: %v", err)
	}
	request, err := value.CapabilityRequestFor("user", "gooo", "capability.query")
	if err != nil {
		t.Fatalf("project request: %v", err)
	}
	if request.DeclarationDigest != value.DeclarationDigest {
		t.Fatalf("declaration binding was not preserved")
	}
}

func TestCapabilityRequestForRejectsUnknownObservation(t *testing.T) {
	digestValue := strings.Repeat("a", 64)
	value, err := NewCapabilityObservationEnvelope(
		"corr-4", digestValue, digestValue, digestValue, digestValue, digestValue,
		[]CapabilityObservation{{ID: "capability.network", State: CapabilityObservationUnknown}},
		2, "inspect-provider-boundary",
	)
	if err != nil {
		t.Fatalf("create envelope: %v", err)
	}
	if _, err := value.CapabilityRequestFor("user", "gooo", "capability.network"); err == nil {
		t.Fatalf("expected unknown observation to remain non-requestable")
	}
}
