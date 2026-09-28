package envelope

import "fmt"

const CapabilityDiscoveryReplayVersion = "capability.discovery.replay.v1"

type CapabilityDiscoveryReplayOutcome string

const (
	CapabilityDiscoveryReplayConfirmed  CapabilityDiscoveryReplayOutcome = "CONFIRMED"
	CapabilityDiscoveryReplayRefuted    CapabilityDiscoveryReplayOutcome = "REFUTED"
	CapabilityDiscoveryReplayUnresolved CapabilityDiscoveryReplayOutcome = "UNRESOLVED"
)

// CapabilityDiscoveryReplay links one bounded discovery receipt to a later
// observation without allowing the observation to execute or authorize work.
type CapabilityDiscoveryReplay struct {
	Version               string                               `json:"schema_version"`
	PriorEvidenceDigest   string                               `json:"prior_evidence_digest"`
	ObservationDigest     string                               `json:"observation_digest"`
	CurrentEvidenceDigest string                               `json:"current_evidence_digest"`
	Outcome               CapabilityDiscoveryReplayOutcome    `json:"outcome"`
	FirstMissingStage     string                               `json:"first_missing_stage,omitempty"`
	NextOperation         string                               `json:"next_operation"`
	NonExecuting          bool                                 `json:"non_executing"`
	NonAuthorizing        bool                                 `json:"non_authorizing"`
}

func (replay CapabilityDiscoveryReplay) Validate() error {
	if replay.Version != CapabilityDiscoveryReplayVersion {
		return fmt.Errorf("unsupported capability discovery replay version %q", replay.Version)
	}
	for name, value := range map[string]string{
		"prior evidence digest": replay.PriorEvidenceDigest,
		"observation digest": replay.ObservationDigest,
		"current evidence digest": replay.CurrentEvidenceDigest,
	} {
		if !validDigest(value) {
			return fmt.Errorf("%s is invalid", name)
		}
	}
	switch replay.Outcome {
	case CapabilityDiscoveryReplayConfirmed:
		if replay.FirstMissingStage != "" {
			return fmt.Errorf("confirmed replay cannot carry a missing stage")
		}
	case CapabilityDiscoveryReplayRefuted, CapabilityDiscoveryReplayUnresolved:
		if replay.FirstMissingStage == "" {
			return fmt.Errorf("non-confirmed replay lost first missing stage")
		}
	default:
		return fmt.Errorf("unsupported capability discovery replay outcome %q", replay.Outcome)
	}
	if replay.NextOperation == "" {
		return fmt.Errorf("capability discovery replay is missing next operation")
	}
	if !replay.NonExecuting || !replay.NonAuthorizing {
		return fmt.Errorf("capability discovery replay crossed a capability boundary")
	}
	return nil
}

func (replay CapabilityDiscoveryReplay) ValidateAgainst(receipt CapabilityDiscoveryReceipt) error {
	if err := receipt.Validate(); err != nil {
		return fmt.Errorf("prior capability discovery receipt: %w", err)
	}
	if err := replay.Validate(); err != nil {
		return err
	}
	if replay.PriorEvidenceDigest != receipt.EvidenceDigest {
		return fmt.Errorf("capability discovery replay does not continue prior evidence")
	}
	return nil
}