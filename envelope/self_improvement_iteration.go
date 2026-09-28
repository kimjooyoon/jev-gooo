package envelope

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

const SelfImprovementIterationVersion = "self.improvement.iteration.v1"

type SelfImprovementIterationState string

const (
	SelfImprovementIterationProposed SelfImprovementIterationState = "PROPOSED"
	SelfImprovementIterationAccepted SelfImprovementIterationState = "ACCEPTED"
	SelfImprovementIterationRejected SelfImprovementIterationState = "REJECTED"
	SelfImprovementIterationUnknown SelfImprovementIterationState = "UNKNOWN"
)

type SelfImprovementIteration struct {
	Version            string                        `json:"version"`
	IterationID        string                        `json:"iteration_id"`
	ObservationDigest  string                        `json:"observation_digest"`
	CandidateDigest    string                        `json:"candidate_digest"`
	VerificationDigest string                        `json:"verification_digest,omitempty"`
	RollbackDigest     string                        `json:"rollback_digest,omitempty"`
	State              SelfImprovementIterationState `json:"state"`
	FirstMissingStage  string                        `json:"first_missing_stage,omitempty"`
	NextOperation      string                        `json:"next_operation,omitempty"`
	NonExecuting       bool                          `json:"non_executing"`
	NonAuthorizing     bool                          `json:"non_authorizing"`
}

func (i SelfImprovementIteration) Validate() error {
	if i.Version != SelfImprovementIterationVersion {
		return fmt.Errorf("unsupported self-improvement iteration version %q", i.Version)
	}
	if strings.TrimSpace(i.IterationID) == "" {
		return fmt.Errorf("iteration_id is required")
	}
	if !validSelfImprovementDigest(i.ObservationDigest) {
		return fmt.Errorf("observation_digest must be a sha256 digest")
	}
	if !validSelfImprovementDigest(i.CandidateDigest) {
		return fmt.Errorf("candidate_digest must be a sha256 digest")
	}
	switch i.State {
	case SelfImprovementIterationProposed:
	case SelfImprovementIterationAccepted, SelfImprovementIterationRejected:
		if !validSelfImprovementDigest(i.VerificationDigest) {
			return fmt.Errorf("verification_digest is required for state %s", i.State)
		}
		if !validSelfImprovementDigest(i.RollbackDigest) {
			return fmt.Errorf("rollback_digest is required for state %s", i.State)
		}
	case SelfImprovementIterationUnknown:
		if strings.TrimSpace(i.FirstMissingStage) == "" {
			return fmt.Errorf("first_missing_stage is required for UNKNOWN")
		}
		if strings.TrimSpace(i.NextOperation) == "" {
			return fmt.Errorf("next_operation is required for UNKNOWN")
		}
	default:
		return fmt.Errorf("unsupported self-improvement iteration state %q", i.State)
	}
	if !i.NonExecuting || !i.NonAuthorizing {
		return fmt.Errorf("self-improvement iteration must be non-executing and non-authorizing")
	}
	return nil
}

func (i SelfImprovementIteration) CanonicalDigest() (string, error) {
	if err := i.Validate(); err != nil {
		return "", err
	}
	payload, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(payload)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func validSelfImprovementDigest(value string) bool {
	if len(value) != len("sha256:")+64 || !strings.HasPrefix(value, "sha256:") {
		return false
	}
	_, err := hex.DecodeString(value[len("sha256:"):])
	return err == nil
}

