package capability

import (
	"fmt"
	"strings"
)

type GuideAction string

const (
	GuideAskClarifyingQuestion GuideAction = "ASK_CLARIFYING_QUESTION"
	GuideDeclareMissingSignals GuideAction = "DECLARE_MISSING_SIGNALS"
	GuideCollectVerifiedEvidence GuideAction = "COLLECT_VERIFIED_EVIDENCE"
)

type Guide struct {
	DiscoveryDigest string       `json:"discovery_digest"`
	PlanDigest      string       `json:"plan_digest"`
	Status          Status       `json:"status"`
	CapabilityID    string       `json:"capability_id"`
	Action          GuideAction  `json:"action"`
	MissingSignals  []string     `json:"missing_signals,omitempty"`
	Questions       []string     `json:"questions,omitempty"`
	Constraints     []string     `json:"constraints"`
	GuideDigest     string       `json:"guide_digest"`
}

var guideConstraints = []string{
	"guide_only_no_provider_invocation",
	"no_authorization_grant",
	"no_catalog_mutation",
	"cache_presence_is_not_semantic_evidence",
}

func BuildGuide(discovery Discovery, plan Plan) (Guide, error) {
	if err := discovery.Validate(); err != nil {
		return Guide{}, err
	}
	if err := plan.Validate(); err != nil {
		return Guide{}, err
	}
	if plan.DiscoveryDigest != discovery.EvidenceDigest || plan.Status != discovery.Status {
		return Guide{}, fmt.Errorf("capability guide inputs do not share discovery identity")
	}

	guide := Guide{
		DiscoveryDigest: discovery.EvidenceDigest,
		PlanDigest:      plan.PlanDigest,
		Status:          discovery.Status,
		CapabilityID:    discovery.CapabilityID,
		Constraints:     append([]string(nil), guideConstraints...),
	}
	switch discovery.Status {
	case StatusUnknown:
		guide.Action = GuideAskClarifyingQuestion
		guide.Questions = append([]string(nil), discovery.NextQuestions...)
	case StatusDeferred:
		guide.Action = GuideDeclareMissingSignals
		guide.MissingSignals = append([]string(nil), discovery.MissingSignals...)
		guide.Questions = append([]string(nil), discovery.NextQuestions...)
	case StatusAvailable:
		guide.Action = GuideCollectVerifiedEvidence
	default:
		return Guide{}, fmt.Errorf("unsupported discovery status %q", discovery.Status)
	}
	guide.GuideDigest = guide.digest()
	return guide, nil
}

func (g Guide) Validate() error {
	if !validDigest(g.DiscoveryDigest) || !validDigest(g.PlanDigest) || !validDigest(g.GuideDigest) {
		return fmt.Errorf("capability guide digest is invalid")
	}
	if len(g.Constraints) != len(guideConstraints) || !sameStrings(g.Constraints, guideConstraints) {
		return fmt.Errorf("capability guide constraints are incomplete")
	}
	switch g.Status {
	case StatusUnknown:
		if g.CapabilityID != "" || g.Action != GuideAskClarifyingQuestion || len(g.MissingSignals) != 0 || len(g.Questions) == 0 {
			return fmt.Errorf("unknown capability guide is incomplete")
		}
	case StatusDeferred:
		if g.CapabilityID == "" || g.Action != GuideDeclareMissingSignals || len(g.MissingSignals) == 0 || len(g.Questions) != len(g.MissingSignals) {
			return fmt.Errorf("deferred capability guide is incomplete")
		}
	case StatusAvailable:
		if g.CapabilityID == "" || g.Action != GuideCollectVerifiedEvidence || len(g.MissingSignals) != 0 || len(g.Questions) != 0 {
			return fmt.Errorf("available capability guide is incomplete")
		}
	default:
		return fmt.Errorf("capability guide status %q is invalid", g.Status)
	}
	if g.digest() != g.GuideDigest {
		return fmt.Errorf("capability guide digest does not match")
	}
	return nil
}

func (g Guide) digest() string {
	return digest(
		"capability-guide",
		g.DiscoveryDigest,
		g.PlanDigest,
		string(g.Status),
		g.CapabilityID,
		string(g.Action),
		strings.Join(g.MissingSignals, ","),
		strings.Join(g.Questions, "|"),
		strings.Join(g.Constraints, "|"),
	)
}

func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}