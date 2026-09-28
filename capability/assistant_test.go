package capability

import (
	"strings"
	"testing"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

func TestDiscoverAssistantBindsCapabilityStages(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package jev\nactivity reverse_observe\nproperty evidence_digest string")
	if err != nil {
		t.Fatal(err)
	}
	assistant, err := DiscoverAssistant("What can this language do with provenance?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	if err := assistant.Validate(); err != nil {
		t.Fatal(err)
	}
	if assistant.Overview.Mode != OverviewMatchedOptions || assistant.Discovery.Status != StatusAvailable {
		t.Fatalf("unexpected assistant stages: %#v", assistant)
	}
	if assistant.Guide.Action != GuideCollectVerifiedEvidence || len(assistant.Plan.Steps) != 4 {
		t.Fatalf("unexpected assistant guidance: %#v", assistant)
	}
}

func TestDiscoverAssistantPreservesUnknown(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo")
	if err != nil {
		t.Fatal(err)
	}
	assistant, err := DiscoverAssistant("How can I model a compiler?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	if err := assistant.Validate(); err != nil {
		t.Fatal(err)
	}
	if assistant.Discovery.Status != StatusUnknown || assistant.Guide.Action != GuideAskClarifyingQuestion {
		t.Fatalf("unexpected unknown assistant: %#v", assistant)
	}
	if len(assistant.Guide.Questions) == 0 {
		t.Fatal("unknown assistant must preserve a next question")
	}
}

func TestDiscoverAssistantRejectsTampering(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo")
	if err != nil {
		t.Fatal(err)
	}
	assistant, err := DiscoverAssistant("How can I model a compiler?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	assistant.EvidenceDigest = strings.Repeat("0", 64)
	if err := assistant.Validate(); err == nil {
		t.Fatal("tampered assistant should be rejected")
	}
}
