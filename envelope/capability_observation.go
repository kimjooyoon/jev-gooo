package envelope

import (
	"fmt"
	"strconv"
	"strings"
)

const CapabilityObservationEnvelopeVersion = "capability-observation.v1"

type CapabilityObservationState string

const (
	CapabilityObservationAvailable CapabilityObservationState = "AVAILABLE"
	CapabilityObservationDeferred  CapabilityObservationState = "DEFERRED"
	CapabilityObservationUnknown   CapabilityObservationState = "UNKNOWN"
)

type CapabilityObservation struct {
	ID    string
	State CapabilityObservationState
}

// CapabilityObservationEnvelope carries declaration-bound capability discovery
// across the JEV boundary without executing or authorizing an operation.
type CapabilityObservationEnvelope struct {
	Version           string
	CorrelationID     string
	SourceDigest      string
	DeclarationDigest string
	CatalogDigest     string
	SchemaDigest      string
	EvidenceDigest    string
	Capabilities      []CapabilityObservation
	FirstMissingStage int
	NextOperation     string
	Executed          bool
	Authorizing       bool
	Digest            string
}

func NewCapabilityObservationEnvelope(
	correlationID, sourceDigest, declarationDigest, catalogDigest, schemaDigest, evidenceDigest string,
	capabilities []CapabilityObservation,
	firstMissingStage int,
	nextOperation string,
) (CapabilityObservationEnvelope, error) {
	envelope := CapabilityObservationEnvelope{
		Version:           CapabilityObservationEnvelopeVersion,
		CorrelationID:     correlationID,
		SourceDigest:      sourceDigest,
		DeclarationDigest: declarationDigest,
		CatalogDigest:     catalogDigest,
		SchemaDigest:      schemaDigest,
		EvidenceDigest:    evidenceDigest,
		Capabilities:      append([]CapabilityObservation(nil), capabilities...),
		FirstMissingStage: firstMissingStage,
		NextOperation:     nextOperation,
	}
	envelope.Digest = envelope.evidenceDigest()
	return envelope, envelope.Validate()
}

func (e CapabilityObservationEnvelope) Validate() error {
	if e.Version != CapabilityObservationEnvelopeVersion {
		return fmt.Errorf("unsupported capability observation envelope version %q", e.Version)
	}
	if strings.TrimSpace(e.CorrelationID) == "" {
		return fmt.Errorf("correlation id is required")
	}
	for name, value := range map[string]string{
		"source digest":      e.SourceDigest,
		"declaration digest": e.DeclarationDigest,
		"catalog digest":     e.CatalogDigest,
		"schema digest":      e.SchemaDigest,
		"evidence digest":    e.EvidenceDigest,
	} {
		if !validDigest(value) {
			return fmt.Errorf("%s is invalid", name)
		}
	}
	if e.Executed {
		return fmt.Errorf("capability observation cannot record execution")
	}
	if e.Authorizing {
		return fmt.Errorf("capability observation cannot authorize an operation")
	}
	if e.FirstMissingStage < -1 {
		return fmt.Errorf("first missing stage must be -1 or greater")
	}
	if len(e.Capabilities) == 0 {
		return fmt.Errorf("at least one capability observation is required")
	}

	seen := make(map[string]struct{}, len(e.Capabilities))
	incomplete := false
	for _, capability := range e.Capabilities {
		if strings.TrimSpace(capability.ID) == "" {
			return fmt.Errorf("capability id is required")
		}
		if _, exists := seen[capability.ID]; exists {
			return fmt.Errorf("duplicate capability id %q", capability.ID)
		}
		seen[capability.ID] = struct{}{}
		switch capability.State {
		case CapabilityObservationAvailable:
		case CapabilityObservationDeferred, CapabilityObservationUnknown:
			incomplete = true
		default:
			return fmt.Errorf("unsupported capability state %q", capability.State)
		}
	}
	if incomplete {
		if e.FirstMissingStage < 0 {
			return fmt.Errorf("incomplete capability observation requires first missing stage")
		}
		if strings.TrimSpace(e.NextOperation) == "" {
			return fmt.Errorf("incomplete capability observation requires next operation")
		}
	} else if e.FirstMissingStage != -1 || strings.TrimSpace(e.NextOperation) != "" {
		return fmt.Errorf("complete capability observation cannot carry a missing stage")
	}
	if !validDigest(e.Digest) {
		return fmt.Errorf("observation envelope digest is invalid")
	}
	if e.evidenceDigest() != e.Digest {
		return fmt.Errorf("observation envelope digest does not match evidence")
	}
	return nil
}

func (e CapabilityObservationEnvelope) evidenceDigest() string {
	parts := []string{
		e.Version,
		e.CorrelationID,
		e.SourceDigest,
		e.DeclarationDigest,
		e.CatalogDigest,
		e.SchemaDigest,
		e.EvidenceDigest,
		strconv.Itoa(e.FirstMissingStage),
		e.NextOperation,
		strconv.FormatBool(e.Executed),
		strconv.FormatBool(e.Authorizing),
	}
	for _, capability := range e.Capabilities {
		parts = append(parts, capability.ID, string(capability.State))
	}
	return digest("capability-observation-envelope", parts...)
}
