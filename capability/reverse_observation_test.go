package capability

import "testing"

func reverseObservationInput() ReverseObservationInput {
	return ReverseObservationInput{
		SourceDigest:            "sha256:source",
		DeclarationDigest:       "sha256:declaration",
		IRDigest:                "sha256:ir",
		GeneratedArtifactDigest: "sha256:generated",
		EvidencePrefixDigest:    "sha256:prefix",
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
	input.ObservedArtifactDigest = "sha256:different"
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
