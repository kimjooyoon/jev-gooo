package capability

import (
	"testing"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

func TestBuildGuideUnknownPreservesClarification(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo")
	if err != nil {
		t.Fatal(err)
	}
	discovery, err := Discover("how can I model a compiler", declaration)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := BuildPlan(discovery)
	if err != nil {
		t.Fatal(err)
	}
	guide, err := BuildGuide(discovery, plan)
	if err != nil {
		t.Fatal(err)
	}
	if err := guide.Validate(); err != nil {
		t.Fatal(err)
	}
	if guide.Action != GuideAskClarifyingQuestion || len(guide.Questions) == 0 {
		t.Fatalf("unexpected unknown guide: %+v", guide)
	}
}

func TestBuildGuideDeferredListsMissingSignals(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo\nactivity observe_result")
	if err != nil {
		t.Fatal(err)
	}
	discovery, err := Discover("reverse observation provenance", declaration)
	if err != nil {
		t.Fatal(err)
	}
	if discovery.Status != StatusDeferred {
		t.Fatalf("expected deferred discovery, got %s", discovery.Status)
	}
	plan, err := BuildPlan(discovery)
	if err != nil {
		t.Fatal(err)
	}
	guide, err := BuildGuide(discovery, plan)
	if err != nil {
		t.Fatal(err)
	}
	if err := guide.Validate(); err != nil {
		t.Fatal(err)
	}
	if guide.Action != GuideDeclareMissingSignals || len(guide.MissingSignals) == 0 {
		t.Fatalf("unexpected deferred guide: %+v", guide)
	}
}

func TestBuildGuideAvailableRequiresVerifiedEvidenceRoute(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo\nactivity reverse_observe\nproperty evidence_digest string")
	if err != nil {
		t.Fatal(err)
	}
	discovery, err := Discover("reverse observation provenance", declaration)
	if err != nil {
		t.Fatal(err)
	}
	if discovery.Status != StatusAvailable {
		t.Fatalf("expected available discovery, got %s", discovery.Status)
	}
	plan, err := BuildPlan(discovery)
	if err != nil {
		t.Fatal(err)
	}
	guide, err := BuildGuide(discovery, plan)
	if err != nil {
		t.Fatal(err)
	}
	if err := guide.Validate(); err != nil {
		t.Fatal(err)
	}
	if guide.Action != GuideCollectVerifiedEvidence || len(guide.Questions) != 0 {
		t.Fatalf("unexpected available guide: %+v", guide)
	}
}

func TestBuildGuideRejectsTampering(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo\nactivity reverse_observe\nproperty evidence_digest string")
	if err != nil {
		t.Fatal(err)
	}
	discovery, err := Discover("reverse observation provenance", declaration)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := BuildPlan(discovery)
	if err != nil {
		t.Fatal(err)
	}
	guide, err := BuildGuide(discovery, plan)
	if err != nil {
		t.Fatal(err)
	}
	guide.Constraints[0] = "changed"
	if err := guide.Validate(); err == nil {
		t.Fatal("expected tampered guide to fail validation")
	}
}