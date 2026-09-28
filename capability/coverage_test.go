package capability

import (
	"testing"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

func TestMeasureCoverageSeparatesObservedSurfaceFromCompleteness(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo\nactivity reverse_observe\nproperty evidence_digest string")
	if err != nil { t.Fatal(err) }
	options, err := DiscoverOptions("provenance and code generation", declaration)
	if err != nil { t.Fatal(err) }
	coverage, err := MeasureCoverage(options)
	if err != nil { t.Fatal(err) }
	if err := coverage.Validate(); err != nil { t.Fatal(err) }
	if coverage.MatchedCount != coverage.AvailableCount+coverage.DeferredCount || coverage.MatchedCount == 0 { t.Fatalf("unexpected coverage counts: %+v", coverage) }
	if coverage.Interpretation != CoverageInterpretation || coverage.Scope != CoverageScope { t.Fatalf("coverage overclaims its scope: %+v", coverage) }
}

func TestMeasureCoverageRecordsUnknownQueryWithoutInventingCapability(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo")
	if err != nil { t.Fatal(err) }
	options, err := DiscoverOptions("how can I model a compiler", declaration)
	if err != nil { t.Fatal(err) }
	coverage, err := MeasureCoverage(options)
	if err != nil { t.Fatal(err) }
	if err := coverage.Validate(); err != nil { t.Fatal(err) }
	if !coverage.UnknownQuery || coverage.MatchedCount != 0 || coverage.AvailableCount != 0 || coverage.DeferredCount != 0 { t.Fatalf("unknown query was not preserved: %+v", coverage) }
}

func TestMeasureCoverageRejectsTampering(t *testing.T) {
	declaration, err := envelope.BindDeclaration("package demo\nactivity reverse_observe\nproperty evidence_digest string")
	if err != nil { t.Fatal(err) }
	options, err := DiscoverOptions("provenance", declaration)
	if err != nil { t.Fatal(err) }
	coverage, err := MeasureCoverage(options)
	if err != nil { t.Fatal(err) }
	coverage.AvailableCount++
	if err := coverage.Validate(); err == nil { t.Fatal("expected tampered coverage to fail validation") }
}