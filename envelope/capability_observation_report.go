package envelope

import "fmt"

type CapabilityObservationProvenance struct {
	SourceDigest      string `json:"source_digest"`
	DeclarationDigest string `json:"declaration_digest"`
	CatalogDigest     string `json:"catalog_digest"`
	SchemaDigest      string `json:"schema_digest"`
	EvidenceDigest    string `json:"evidence_digest"`
}

type CapabilityObservationReport struct {
	Version           string                          `json:"version"`
	CorrelationID     string                          `json:"correlation_id"`
	State             CapabilityObservationState      `json:"state"`
	CapabilityIDs     []string                        `json:"capability_ids"`
	FirstMissingStage int                             `json:"first_missing_stage"`
	NextOperation     string                          `json:"next_operation"`
	Provenance        CapabilityObservationProvenance `json:"provenance"`
	ObservationDigest string                          `json:"observation_digest"`
	Requestable       bool                            `json:"requestable"`
	Executed          bool                            `json:"executed"`
	Authorizing       bool                            `json:"authorizing"`
}

func (e CapabilityObservationEnvelope) Report() (CapabilityObservationReport, error) {
	if err := e.Validate(); err != nil {
		return CapabilityObservationReport{}, err
	}
	state := CapabilityObservationAvailable
	for _, capability := range e.Capabilities {
		switch capability.State {
		case CapabilityObservationDeferred:
			state = CapabilityObservationDeferred
		case CapabilityObservationUnknown:
			if state == CapabilityObservationAvailable {
				state = CapabilityObservationUnknown
			}
		}
	}
	ids := make([]string, 0, len(e.Capabilities))
	for _, capability := range e.Capabilities {
		ids = append(ids, capability.ID)
	}
	report := CapabilityObservationReport{
		Version:           e.Version,
		CorrelationID:     e.CorrelationID,
		State:             state,
		CapabilityIDs:     ids,
		FirstMissingStage: e.FirstMissingStage,
		NextOperation:     e.NextOperation,
		Provenance: CapabilityObservationProvenance{
			SourceDigest:      e.SourceDigest,
			DeclarationDigest: e.DeclarationDigest,
			CatalogDigest:     e.CatalogDigest,
			SchemaDigest:      e.SchemaDigest,
			EvidenceDigest:    e.EvidenceDigest,
		},
		ObservationDigest: e.Digest,
		Requestable:       state == CapabilityObservationAvailable,
		Executed:          e.Executed,
		Authorizing:       e.Authorizing,
	}
	if err := report.Validate(); err != nil {
		return CapabilityObservationReport{}, err
	}
	return report, nil
}

func (r CapabilityObservationReport) Validate() error {
	if r.Version != CapabilityObservationEnvelopeVersion || r.CorrelationID == "" {
		return fmt.Errorf("capability observation report identity is incomplete")
	}
	switch r.State {
	case CapabilityObservationAvailable, CapabilityObservationDeferred, CapabilityObservationUnknown:
	default:
		return fmt.Errorf("capability observation report state %q is invalid", r.State)
	}
	if len(r.CapabilityIDs) == 0 {
		return fmt.Errorf("capability observation report has no capability ids")
	}
	for _, id := range r.CapabilityIDs {
		if id == "" {
			return fmt.Errorf("capability observation report has an empty capability id")
		}
	}
	for name, value := range map[string]string{
		"source":      r.Provenance.SourceDigest,
		"declaration": r.Provenance.DeclarationDigest,
		"catalog":     r.Provenance.CatalogDigest,
		"schema":      r.Provenance.SchemaDigest,
		"evidence":    r.Provenance.EvidenceDigest,
		"observation": r.ObservationDigest,
	} {
		if !validDigest(value) {
			return fmt.Errorf("capability observation report %s digest is invalid", name)
		}
	}
	if r.Executed || r.Authorizing {
		return fmt.Errorf("capability observation report crossed an execution or authorization boundary")
	}
	if r.Requestable != (r.State == CapabilityObservationAvailable) {
		return fmt.Errorf("capability observation report requestable state is inconsistent")
	}
	if r.State == CapabilityObservationAvailable {
		if r.FirstMissingStage != -1 || r.NextOperation != "" {
			return fmt.Errorf("available capability observation report carries a missing boundary")
		}
	} else if r.FirstMissingStage < 0 || r.NextOperation == "" {
		return fmt.Errorf("incomplete capability observation report lost its boundary")
	}
	return nil
}
