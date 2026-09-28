package capability

import (
	"strings"
	"testing"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

func feedbackDeclaration(t *testing.T, source string) envelope.Declaration {
	t.Helper()
	declaration, err := envelope.BindDeclaration(source)
	if err != nil {
		t.Fatal(err)
	}
	return declaration
}

func availableDiscovery(t *testing.T) Discovery {
	t.Helper()
	declaration := feedbackDeclaration(t, "package jev\nactivity discover_capability\nproperty evidence_digest string")
	discovery, err := Discover("What can this language do?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	return discovery
}

func TestObserveFeedbackRequiresEvidence(t *testing.T) {
	feedback := ObserveFeedback(availableDiscovery(t), FeedbackInput{Source: "human note"})
	if feedback.Disposition != DispositionPendingRequiredEvidence {
		t.Fatalf("unexpected disposition: %#v", feedback)
	}
	if err := feedback.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestObserveFeedbackPreservesDeferred(t *testing.T) {
	declaration := feedbackDeclaration(t, "package jev\nactivity inspect")
	discovery, err := Discover("Can this language run it?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	feedback := ObserveFeedback(discovery, FeedbackInput{
		Source:         "verified-looking note",
		EvidenceDigest: strings.Repeat("a", 64),
		Verified:       true,
	})
	if feedback.Disposition != DispositionPreserveDeferred {
		t.Fatalf("unexpected disposition: %#v", feedback)
	}
	if err := feedback.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestObserveFeedbackEligibleOnlyWithExplicitEvidence(t *testing.T) {
	feedback := ObserveFeedback(availableDiscovery(t), FeedbackInput{
		Source:         "external evaluator receipt",
		EvidenceDigest: strings.Repeat("b", 64),
		Verified:       true,
	})
	if feedback.Disposition != DispositionEligibleForCatalogReview {
		t.Fatalf("unexpected disposition: %#v", feedback)
	}
	if err := feedback.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestObserveFeedbackRejectsTampering(t *testing.T) {
	feedback := ObserveFeedback(availableDiscovery(t), FeedbackInput{Source: "human note"})
	feedback.FeedbackDigest = strings.Repeat("0", 64)
	if err := feedback.Validate(); err == nil {
		t.Fatal("tampered capability feedback should be rejected")
	}
}
