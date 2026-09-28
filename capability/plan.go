package capability

import (
	"fmt"
	"strings"
)

type PlanStepStatus string

const (
	PlanStepObserved PlanStepStatus = "OBSERVED"
	PlanStepPending  PlanStepStatus = "PENDING"
)

type PlanStep struct {
	ID            string         `json:"id"`
	Status        PlanStepStatus `json:"status"`
	MissingSignal string         `json:"missing_signal,omitempty"`
	Question      string         `json:"question,omitempty"`
}

type Plan struct {
	DiscoveryDigest string     `json:"discovery_digest"`
	Status          Status     `json:"status"`
	Steps           []PlanStep `json:"steps"`
	PlanDigest      string     `json:"plan_digest"`
}

func BuildPlan(discovery Discovery) (Plan, error) {
	if err := discovery.Validate(); err != nil {
		return Plan{}, err
	}
	plan := Plan{
		DiscoveryDigest: discovery.EvidenceDigest,
		Status:          discovery.Status,
		Steps: []PlanStep{
			{ID: "match-capability", Status: PlanStepObserved},
			{ID: "bind-declaration", Status: PlanStepObserved},
		},
	}
	switch discovery.Status {
	case StatusUnknown:
		plan.Steps = append(plan.Steps, PlanStep{
			ID:       "ask-clarifying-question",
			Status:   PlanStepPending,
			Question: discovery.NextQuestions[0],
		})
	case StatusDeferred:
		for index, signal := range discovery.MissingSignals {
			question := ""
			if index < len(discovery.NextQuestions) {
				question = discovery.NextQuestions[index]
			}
			plan.Steps = append(plan.Steps, PlanStep{
				ID:            "declare-missing-signal-" + signal,
				Status:        PlanStepPending,
				MissingSignal: signal,
				Question:      question,
			})
		}
	case StatusAvailable:
		plan.Steps = append(plan.Steps,
			PlanStep{ID: "collect-verified-evidence", Status: PlanStepPending},
			PlanStep{ID: "review-feedback", Status: PlanStepPending},
		)
	default:
		return Plan{}, fmt.Errorf("unsupported discovery status %q", discovery.Status)
	}
	plan.PlanDigest = plan.digest()
	return plan, nil
}

func (p Plan) Validate() error {
	if !validDigest(p.DiscoveryDigest) || !validDigest(p.PlanDigest) || len(p.Steps) == 0 {
		return fmt.Errorf("capability plan is incomplete")
	}
	switch p.Status {
	case StatusAvailable, StatusDeferred, StatusUnknown:
	default:
		return fmt.Errorf("capability plan status %q is invalid", p.Status)
	}
	for _, step := range p.Steps {
		if strings.TrimSpace(step.ID) == "" {
			return fmt.Errorf("capability plan step id is required")
		}
		if step.Status != PlanStepObserved && step.Status != PlanStepPending {
			return fmt.Errorf("capability plan step status %q is invalid", step.Status)
		}
	}
	if p.digest() != p.PlanDigest {
		return fmt.Errorf("capability plan digest does not match")
	}
	return nil
}

func (p Plan) digest() string {
	parts := []string{"capability-plan", p.DiscoveryDigest, string(p.Status)}
	for _, step := range p.Steps {
		parts = append(parts, step.ID, string(step.Status), step.MissingSignal, step.Question)
	}
	return digest(parts...)
}
