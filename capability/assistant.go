package capability

import (
	"fmt"
	"strings"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

type Assistant struct {
	Query             string   `json:"query"`
	DeclarationDigest string   `json:"declaration_digest"`
	Overview          Overview `json:"overview"`
	Discovery         Discovery `json:"discovery"`
	Plan              Plan     `json:"plan"`
	Guide             Guide    `json:"guide"`
	Constraints       []string `json:"constraints"`
	EvidenceDigest    string   `json:"evidence_digest"`
}

var assistantConstraints = []string{
	"bounded_catalog_discovery",
	"unknown_is_preserved",
	"no_provider_invocation",
	"no_authorization_grant",
	"no_catalog_mutation",
	"cache_presence_is_not_semantic_evidence",
}

func DiscoverAssistant(query string, declaration envelope.Declaration) (Assistant, error) {
	if err := declaration.Validate(); err != nil {
		return Assistant{}, err
	}
	query = strings.TrimSpace(query)
	if query == "" {
		return Assistant{}, fmt.Errorf("capability assistant query is required")
	}

	overview, err := DiscoverOverview(query, declaration)
	if err != nil {
		return Assistant{}, err
	}
	discovery, err := Discover(query, declaration)
	if err != nil {
		return Assistant{}, err
	}
	plan, err := BuildPlan(discovery)
	if err != nil {
		return Assistant{}, err
	}
	guide, err := BuildGuide(discovery, plan)
	if err != nil {
		return Assistant{}, err
	}

	result := Assistant{
		Query:             query,
		DeclarationDigest: declaration.Digest,
		Overview:          overview,
		Discovery:         discovery,
		Plan:              plan,
		Guide:             guide,
		Constraints:       append([]string(nil), assistantConstraints...),
	}
	result.EvidenceDigest = result.digest()
	return result, nil
}

func (a Assistant) Validate() error {
	if strings.TrimSpace(a.Query) == "" || !validDigest(a.DeclarationDigest) || !validDigest(a.EvidenceDigest) {
		return fmt.Errorf("capability assistant evidence is incomplete")
	}
	if len(a.Constraints) != len(assistantConstraints) || !sameStrings(a.Constraints, assistantConstraints) {
		return fmt.Errorf("capability assistant constraints are incomplete")
	}
	if err := a.Overview.Validate(); err != nil {
		return fmt.Errorf("validate assistant overview: %w", err)
	}
	if err := a.Discovery.Validate(); err != nil {
		return fmt.Errorf("validate assistant discovery: %w", err)
	}
	if err := a.Plan.Validate(); err != nil {
		return fmt.Errorf("validate assistant plan: %w", err)
	}
	if err := a.Guide.Validate(); err != nil {
		return fmt.Errorf("validate assistant guide: %w", err)
	}
	if a.Overview.Query != a.Query || a.Discovery.Query != a.Query {
		return fmt.Errorf("assistant stages do not share query identity")
	}
	if a.Overview.DeclarationDigest != a.DeclarationDigest || a.Discovery.DeclarationDigest != a.DeclarationDigest {
		return fmt.Errorf("assistant stages do not share declaration identity")
	}
	if a.Plan.DiscoveryDigest != a.Discovery.EvidenceDigest || a.Guide.DiscoveryDigest != a.Discovery.EvidenceDigest {
		return fmt.Errorf("assistant stages do not share discovery identity")
	}
	if a.Guide.PlanDigest != a.Plan.PlanDigest {
		return fmt.Errorf("assistant stages do not share plan identity")
	}
	if a.digest() != a.EvidenceDigest {
		return fmt.Errorf("capability assistant evidence digest does not match")
	}
	return nil
}

func (a Assistant) digest() string {
	return digest(
		"capability-assistant",
		a.Query,
		a.DeclarationDigest,
		a.Overview.EvidenceDigest,
		a.Discovery.EvidenceDigest,
		a.Plan.PlanDigest,
		a.Guide.GuideDigest,
		strings.Join(a.Constraints, "|"),
	)
}
