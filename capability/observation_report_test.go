package capability

import (
	"strings"
	"testing"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

func TestBuildAssistantObservationReportCountsBoundedStates(t *testing.T) {
	availableDeclaration, err := envelope.BindDeclaration("package jev\nactivity reverse_observe\nproperty evidence_digest string")
	if err != nil {
		t.Fatal(err)
	}
	deferredDeclaration, err := envelope.BindDeclaration("package demo")
	if err != nil {
		t.Fatal(err)
	}
	available, err := DiscoverAssistant("What can this language do with provenance?", availableDeclaration)
	if err != nil {
		t.Fatal(err)
	}
	deferred, err := DiscoverAssistant("What can this language do with provenance?", deferredDeclaration)
	if err != nil {
		t.Fatal(err)
	}
	report, err := BuildAssistantObservationReport("What can this language do with provenance?", []AssistantObservationItem{
		{Path: "available.gooo", Assistant: available},
		{Path: "deferred.gooo", Assistant: deferred},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := report.Validate(); err != nil {
		t.Fatal(err)
	}
	if report.Counts != (ObservationCounts{Total: 2, Available: 1, Deferred: 1, Unknown: 0}) {
		t.Fatalf("counts = %#v", report.Counts)
	}
}

func TestBuildAssistantObservationReportRejectsDuplicatePaths(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo")
	if err != nil {
		t.Fatal(err)
	}
	assistant, err := DiscoverAssistant("What can this language do?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := BuildAssistantObservationReport("What can this language do?", []AssistantObservationItem{
		{Path: "same.gooo", Assistant: assistant},
		{Path: "same.gooo", Assistant: assistant},
	}); err == nil {
		t.Fatal("duplicate paths should be rejected")
	}
}

func TestBuildAssistantObservationReportRejectsTampering(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo")
	if err != nil {
		t.Fatal(err)
	}
	assistant, err := DiscoverAssistant("What can this language do?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	report, err := BuildAssistantObservationReport("What can this language do?", []AssistantObservationItem{
		{Path: "demo.gooo", Assistant: assistant},
	})
	if err != nil {
		t.Fatal(err)
	}
	report.EvidenceDigest = strings.Repeat("0", 64)
	if err := report.Validate(); err == nil {
		t.Fatal("tampered report should be rejected")
	}
}
