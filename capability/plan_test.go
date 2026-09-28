package capability

import (
	"strings"
	"testing"
)

func TestBuildPlanForDeferredDiscovery(t *testing.T) {
	discovery, err := DiscoverFromSource(
		"Can this language run it?",
		"package jev\nactivity inspect",
	)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := BuildPlan(discovery)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Status != StatusDeferred || len(plan.Steps) != 4 {
		t.Fatalf("unexpected plan: %#v", plan)
	}
	if err := plan.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestBuildPlanKeepsUnknownUnknown(t *testing.T) {
	discovery, err := DiscoverFromSource(
		"How should I choose a color palette?",
		"package jev\nactivity inspect",
	)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := BuildPlan(discovery)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Status != StatusUnknown || plan.Steps[len(plan.Steps)-1].Status != PlanStepPending {
		t.Fatalf("unexpected plan: %#v", plan)
	}
	if err := plan.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestBuildPlanRejectsTampering(t *testing.T) {
	discovery, err := DiscoverFromSource(
		"What can this language do?",
		"package jev\nactivity discover_capability\nproperty evidence_digest string",
	)
	if err != nil {
		t.Fatal(err)
	}
	plan, err := BuildPlan(discovery)
	if err != nil {
		t.Fatal(err)
	}
	plan.PlanDigest = strings.Repeat("0", 64)
	if err := plan.Validate(); err == nil {
		t.Fatal("tampered plan should be rejected")
	}
}
