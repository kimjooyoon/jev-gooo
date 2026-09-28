package envelope

import (
	"fmt"
	"strings"
)

const SelfImprovementIterationReplayVersion = "self.improvement.iteration.replay.v1"

type SelfImprovementIterationReplayOutcome string

const (
	SelfImprovementIterationReplayConfirmed   SelfImprovementIterationReplayOutcome = "CONFIRMED"
	SelfImprovementIterationReplayRefuted     SelfImprovementIterationReplayOutcome = "REFUTED"
	SelfImprovementIterationReplayUnresolved  SelfImprovementIterationReplayOutcome = "UNRESOLVED"
)

type SelfImprovementIterationReplay struct {
	Version                string                              `json:"version"`
	PriorIterationDigest   string                              `json:"prior_iteration_digest"`
	ObservationDigest      string                              `json:"observation_digest"`
	CurrentIterationDigest string                              `json:"current_iteration_digest,omitempty"`
	Outcome                SelfImprovementIterationReplayOutcome `json:"outcome"`
	FirstMissingStage      string                              `json:"first_missing_stage,omitempty"`
	NextOperation          string                              `json:"next_operation,omitempty"`
	NonExecuting           bool                                `json:"non_executing"`
	NonAuthorizing         bool                                `json:"non_authorizing"`
}

func (r SelfImprovementIterationReplay) Validate() error {
	if r.Version != SelfImprovementIterationReplayVersion {
		return fmt.Errorf("unsupported self-improvement iteration replay version %q", r.Version)
	}
	if !validSelfImprovementDigest(r.PriorIterationDigest) {
		return fmt.Errorf("prior_iteration_digest must be a sha256 digest")
	}
	if !validSelfImprovementDigest(r.ObservationDigest) {
		return fmt.Errorf("observation_digest must be a sha256 digest")
	}
	switch r.Outcome {
	case SelfImprovementIterationReplayConfirmed, SelfImprovementIterationReplayRefuted:
		if !validSelfImprovementDigest(r.CurrentIterationDigest) {
			return fmt.Errorf("current_iteration_digest is required for outcome %s", r.Outcome)
		}
	case SelfImprovementIterationReplayUnresolved:
		if strings.TrimSpace(r.FirstMissingStage) == "" {
			return fmt.Errorf("first_missing_stage is required for UNRESOLVED")
		}
		if strings.TrimSpace(r.NextOperation) == "" {
			return fmt.Errorf("next_operation is required for UNRESOLVED")
		}
	default:
		return fmt.Errorf("unsupported self-improvement iteration replay outcome %q", r.Outcome)
	}
	if !r.NonExecuting || !r.NonAuthorizing {
		return fmt.Errorf("self-improvement iteration replay must be non-executing and non-authorizing")
	}
	return nil
}

func (r SelfImprovementIterationReplay) ValidateAgainst(iteration SelfImprovementIteration) error {
	if err := r.Validate(); err != nil {
		return err
	}
	digest, err := iteration.CanonicalDigest()
	if err != nil {
		return err
	}
	if r.PriorIterationDigest != digest {
		return fmt.Errorf("prior_iteration_digest does not match the supplied iteration")
	}
	return nil
}

