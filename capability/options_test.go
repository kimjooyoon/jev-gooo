package capability

import (
	"testing"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

func TestDiscoverOptionsReturnsDeterministicAlternatives(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo\nactivity reverse_observe\nactivity generate_output\nproperty evidence_digest string\nproperty output string")
	if err != nil {
		t.Fatal(err)
	}
	options, err := DiscoverOptions("What can this language do with provenance and code generation?", declaration)
	if err != nil {
		t.Fatal(err)
	}
	if err := options.Validate(); err != nil {
		t.Fatal(err)
	}
	if len(options.Options) < 2 {
		t.Fatalf("expected alternatives, got %+v", options)
	}
	for _, option := range options.Options {
		if option.Status != StatusAvailable {
			t.Fatalf("expected declaration-bound option to be available, got %+v", option)
		}
	}
}

func TestDiscoverOptionsPreservesDeferredCandidates(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo\nactivity observe_result")
	if err != nil {
		t.Fatal(err)
	}
	options, err := DiscoverOptions("provenance and code generation", declaration)
	if err != nil {
		t.Fatal(err)
	}
	if err := options.Validate(); err != nil {
		t.Fatal(err)
	}
	if len(options.Options) < 2 {
		t.Fatalf("expected multiple deferred candidates, got %+v", options)
	}
	foundDeferred := false
	for _, option := range options.Options {
		if option.Status == StatusDeferred && len(option.MissingSignals) > 0 {
			foundDeferred = true
		}
	}
	if !foundDeferred {
		t.Fatalf("expected deferred candidate, got %+v", options)
	}
}

func TestDiscoverOptionsUnknownQueryHasBoundEmptyResult(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo")
	if err != nil {
		t.Fatal(err)
	}
	options, err := DiscoverOptions("how can I model a compiler", declaration)
	if err != nil {
		t.Fatal(err)
	}
	if err := options.Validate(); err != nil {
		t.Fatal(err)
	}
	if len(options.Options) != 0 {
		t.Fatalf("expected no options, got %+v", options.Options)
	}
}

func TestDiscoverOptionsRejectsTamperedDigest(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo\nactivity observe_result\nproperty evidence_digest string")
	if err != nil {
		t.Fatal(err)
	}
	options, err := DiscoverOptions("provenance", declaration)
	if err != nil {
		t.Fatal(err)
	}
	options.EvidenceDigest = digest("tampered")
	if err := options.Validate(); err == nil {
		t.Fatal("expected tampered options to fail validation")
	}
}