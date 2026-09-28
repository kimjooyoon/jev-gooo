package envelope

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
)

// DomainCapabilityMeasurementEnvelope preserves an observation-based domain
// measurement for later replay. It is not a semantic proof, completion score,
// execution plan, or authorization grant.
type DomainCapabilityMeasurementEnvelope struct {
	Version              string   `json:"version"`
	Domain               string   `json:"domain"`
	ExpectedCapabilities []string `json:"expected_capabilities"`
	ObservedCapabilities []string `json:"observed_capabilities"`
	MissingCapabilities  []string `json:"missing_capabilities"`
	CoverageRatio        float64  `json:"coverage_ratio"`
	UsefulRatio          float64  `json:"useful_ratio"`
	UnresolvedRatio      float64  `json:"unresolved_ratio"`
	ObservationCount     int      `json:"observation_count"`
	DecisionSignal       string   `json:"decision_signal"`
	MeasurementSemantics string   `json:"measurement_semantics"`
	EvidenceDigests      []string `json:"evidence_digests"`
	NonExecuting         bool     `json:"non_executing"`
	NonAuthorizing       bool     `json:"non_authorizing"`
}

// Validate checks structural provenance and safety invariants without claiming
// that the measured domain is complete or correct.
func (m DomainCapabilityMeasurementEnvelope) Validate() error {
	if strings.TrimSpace(m.Version) == "" {
		return fmt.Errorf("measurement version is required")
	}
	if strings.TrimSpace(m.Domain) == "" {
		return fmt.Errorf("measurement domain is required")
	}
	if len(m.ExpectedCapabilities) == 0 {
		return fmt.Errorf("expected capabilities are required")
	}
	if !sortedUniqueNonEmpty(m.ExpectedCapabilities) {
		return fmt.Errorf("expected capabilities must be sorted, unique, and non-empty")
	}
	if !sortedUniqueNonEmptyOrEmpty(m.ObservedCapabilities) || !sortedUniqueNonEmptyOrEmpty(m.MissingCapabilities) {
		return fmt.Errorf("observed and missing capabilities must be sorted and unique")
	}
	expected := make(map[string]struct{}, len(m.ExpectedCapabilities))
	for _, capability := range m.ExpectedCapabilities {
		expected[capability] = struct{}{}
	}
	observed := make(map[string]struct{}, len(m.ObservedCapabilities))
	for _, capability := range m.ObservedCapabilities {
		if _, ok := expected[capability]; !ok {
			return fmt.Errorf("observed capability %q is outside the expected set", capability)
		}
		observed[capability] = struct{}{}
	}
	wantMissing := make([]string, 0, len(m.ExpectedCapabilities))
	for _, capability := range m.ExpectedCapabilities {
		if _, ok := observed[capability]; !ok {
			wantMissing = append(wantMissing, capability)
		}
	}
	if !equalStrings(wantMissing, m.MissingCapabilities) {
		return fmt.Errorf("missing capabilities do not match the expected minus observed set")
	}
	if math.IsNaN(m.CoverageRatio) || math.IsInf(m.CoverageRatio, 0) || m.CoverageRatio < 0 || m.CoverageRatio > 1 {
		return fmt.Errorf("coverage ratio must be finite and between 0 and 1")
	}
	if math.IsNaN(m.UsefulRatio) || math.IsInf(m.UsefulRatio, 0) || m.UsefulRatio < 0 || m.UsefulRatio > 1 {
		return fmt.Errorf("useful ratio must be finite and between 0 and 1")
	}
	if math.IsNaN(m.UnresolvedRatio) || math.IsInf(m.UnresolvedRatio, 0) || m.UnresolvedRatio < 0 || m.UnresolvedRatio > 1 {
		return fmt.Errorf("unresolved ratio must be finite and between 0 and 1")
	}
	if m.ObservationCount < 0 {
		return fmt.Errorf("observation count cannot be negative")
	}
	if m.ObservationCount > 0 && len(m.EvidenceDigests) == 0 {
		return fmt.Errorf("observations require evidence digests")
	}
	if !sortedUniqueNonEmptyOrEmpty(m.EvidenceDigests) {
		return fmt.Errorf("evidence digests must be sorted, unique, and non-empty")
	}
	switch m.DecisionSignal {
	case "COLLECT_MORE_OBSERVATIONS", "REVIEW_MISSING_CAPABILITIES", "REVIEW_INVESTMENT", "NO_ACTIONABLE_SIGNAL":
	default:
		return fmt.Errorf("unsupported decision signal %q", m.DecisionSignal)
	}
	if strings.TrimSpace(m.MeasurementSemantics) == "" {
		return fmt.Errorf("measurement semantics are required")
	}
	if !m.NonExecuting || !m.NonAuthorizing {
		return fmt.Errorf("measurement envelope must remain non-executing and non-authorizing")
	}
	return nil
}

// CanonicalDigest returns a stable digest for a validated envelope.
func (m DomainCapabilityMeasurementEnvelope) CanonicalDigest() (string, error) {
	if err := m.Validate(); err != nil {
		return "", err
	}
	encoded, err := json.Marshal(m)
	if err != nil {
		return "", fmt.Errorf("marshal measurement envelope: %w", err)
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}

func sortedUniqueNonEmpty(values []string) bool {
	if len(values) == 0 {
		return false
	}
	return sortedUniqueNonEmptyOrEmpty(values)
}

func sortedUniqueNonEmptyOrEmpty(values []string) bool {
	for index, value := range values {
		if strings.TrimSpace(value) == "" {
			return false
		}
		if index > 0 && values[index-1] >= value {
			return false
		}
	}
	return true
}

func equalStrings(left, right []string) bool {
	return len(left) == len(right) && sort.Strings(append([]string(nil), left...)) == nil && strings.Join(left, "\x00") == strings.Join(right, "\x00")
}
