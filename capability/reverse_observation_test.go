package capability

import (
	"strings"
	"testing"
)

func reverseObservationInput() ReverseObservationInput {
	return ReverseObservationInput{
		SourceDigest:            strings.Repeat("1", 64),
		DeclarationDigest:       strings.Repeat("2", 64),
		IRDigest:                strings.Repeat("3", 64),
		GeneratedArtifactDigest: strings.Repeat("4", 64),
		EvidencePrefixDigest:    strings.Repeat("5", 64),
	}
}

func TestObserveReverseObservationBindsGeneratedArtifact(t *testing.T) {
	input := reverseObservationInput()
	input.ObservedArtifactDigest = input.GeneratedArtifactDigest
	receipt := ObserveReverseObservation(input)
	if receipt.Status != ReverseObservationObserved || receipt.FirstMismatch != "" || receipt.MissingStage != "" {
		t.Fatalf("unexpected observed receipt: %+v", receipt)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("observed receipt should validate: %v", err)
	}
}

func TestObserveReverseObservationPreservesMismatch(t *testing.T) {
	input := reverseObservationInput()
	input.ObservedArtifactDigest = strings.Repeat("6", 64)
	receipt := ObserveReverseObservation(input)
	if receipt.Status != ReverseObservationMismatch || receipt.FirstMismatch != "reverse_observation" {
		t.Fatalf("unexpected mismatch receipt: %+v", receipt)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("mismatch receipt should validate: %v", err)
	}
}

func TestObserveReverseObservationDefersMissingObservation(t *testing.T) {
	receipt := ObserveReverseObservation(reverseObservationInput())
	if receipt.Status != ReverseObservationDeferred || receipt.MissingStage != "reverse_observation" {
		t.Fatalf("unexpected deferred receipt: %+v", receipt)
	}
	if err := receipt.Validate(); err != nil {
		t.Fatalf("deferred receipt should validate: %v", err)
	}
}
