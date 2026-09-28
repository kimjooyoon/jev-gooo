package envelope

import (
	"strings"
	"testing"
)

func TestCapabilityObservationEnvelopePreservesBoundary(t *testing.T) {
	digestValue := strings.Repeat("a", 64)
	value, err := NewCapabilityObservationEnvelope(
		"corr-1", digestValue, digestValue, digestValue, digestValue, digestValue,
		[]CapabilityObservation{{ID: "capability.query", State: CapabilityObservationAvailable}},
		-1, "",
	)
	if err != nil {
		t.Fatalf("create envelope: %v", err)
	}
	if err := value.Validate(); err != nil {
		t.Fatalf("validate envelope: %v", err)
	}
	value.EvidenceDigest = strings.Repeat("b", 64)
	if err := value.Validate(); err == nil {
		t.Fatalf("expected tampered evidence digest to fail")
	}
}

func TestCapabilityObservationEnvelopeRequiresNextOperation(t *testing.T) {
	digestValue := strings.Repeat("a", 64)
	if _, err := NewCapabilityObservationEnvelope(
		"corr-2", digestValue, digestValue, digestValue, digestValue, digestValue,
		[]CapabilityObservation{{ID: "capability.network", State: CapabilityObservationUnknown}},
		1, "",
	); err == nil {
		t.Fatalf("expected unknown observation without next operation to fail")
	}
}
