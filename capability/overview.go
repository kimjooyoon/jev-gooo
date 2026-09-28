package capability

import (
	"fmt"
	"strings"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

type OverviewMode string

const (
	OverviewCatalog        OverviewMode = "CATALOG_OVERVIEW"
	OverviewMatchedOptions OverviewMode = "MATCHED_OPTIONS"
	OverviewClarification  OverviewMode = "CLARIFICATION_REQUIRED"
)

type OverviewCapability struct {
	ID              string   `json:"id"`
	Label           string   `json:"label"`
	Status          Status   `json:"status"`
	RequiredSignals []string `json:"required_signals"`
	MissingSignals  []string `json:"missing_signals,omitempty"`
}

type Overview struct {
	Query              string               `json:"query"`
	DeclarationDigest  string               `json:"declaration_digest"`
	Mode               OverviewMode         `json:"mode"`
	Summary            string               `json:"summary"`
	Capabilities       []OverviewCapability `json:"capabilities"`
	SuggestedQuestions []string             `json:"suggested_questions"`
	Constraints        []string             `json:"constraints"`
	EvidenceDigest     string               `json:"evidence_digest"`
}

var overviewConstraints = []string{
	"catalog_scope_only",
	"not_a_language_completeness_claim",
	"no_provider_invocation",
	"no_authorization_grant",
	"no_catalog_mutation",
	"cache_presence_is_not_semantic_evidence",
}

func DiscoverOverview(query string, declaration envelope.Declaration) (Overview, error) {
	if err := declaration.Validate(); err != nil {
		return Overview{}, err
	}
	query = strings.TrimSpace(query)
	if query == "" {
		return Overview{}, fmt.Errorf("capability overview query is required")
	}

	options, err := DiscoverOptions(query, declaration)
	if err != nil {
		return Overview{}, err
	}
	matchedSignals := declarationSignals(declaration.Source)
	result := Overview{
		Query:             query,
		DeclarationDigest: declaration.Digest,
		Constraints:       append([]string(nil), overviewConstraints...),
	}

	switch {
	case isOverviewQuery(query):
		result.Mode = OverviewCatalog
		result.Capabilities = overviewCapabilitiesFromCatalog(matchedSignals)
		result.Summary = fmt.Sprintf("gooo can describe %d catalog capabilities for this declaration; this is not a completeness claim.", len(result.Capabilities))
		result.SuggestedQuestions = overviewQuestions()
	case len(options.Options) > 0:
		result.Mode = OverviewMatchedOptions
		result.Capabilities = overviewCapabilitiesFromOptions(options.Options)
		result.Summary = fmt.Sprintf("gooo matched %d catalog capabilities; inspect each status before investing.", len(result.Capabilities))
		result.SuggestedQuestions = questionsForCapabilities(result.Capabilities)
	default:
		result.Mode = OverviewClarification
		result.Summary = "No explicit catalog capability matched; ask a narrower question so gooo can preserve the boundary."
		result.SuggestedQuestions = []string{
			"어떤 작업을 하려 하나요: 생성, 역관찰, 실행 관찰, LSP 피드백 중 무엇인가요?",
			"어떤 .gooo activity와 evidence를 선언할 수 있나요?",
		}
	}
	result.EvidenceDigest = result.digest()
	return result, nil
}

func (o Overview) Validate() error {
	if strings.TrimSpace(o.Query) == "" || !validDigest(o.DeclarationDigest) || !validDigest(o.EvidenceDigest) {
		return fmt.Errorf("capability overview evidence is incomplete")
	}
	if strings.TrimSpace(o.Summary) == "" || len(o.Constraints) != len(overviewConstraints) || !sameStrings(o.Constraints, overviewConstraints) {
		return fmt.Errorf("capability overview constraints are incomplete")
	}
	switch o.Mode {
	case OverviewCatalog:
		if len(o.Capabilities) != len(catalog) || len(o.SuggestedQuestions) == 0 {
			return fmt.Errorf("catalog capability overview is incomplete")
		}
	case OverviewMatchedOptions:
		if len(o.Capabilities) == 0 {
			return fmt.Errorf("matched capability overview is empty")
		}
	case OverviewClarification:
		if len(o.Capabilities) != 0 || len(o.SuggestedQuestions) == 0 {
			return fmt.Errorf("clarification overview must preserve questions")
		}
	default:
		return fmt.Errorf("capability overview mode %q is invalid", o.Mode)
	}
	seen := make(map[string]struct{}, len(o.Capabilities))
	for index, capability := range o.Capabilities {
		if strings.TrimSpace(capability.ID) == "" || strings.TrimSpace(capability.Label) == "" || len(capability.RequiredSignals) == 0 {
			return fmt.Errorf("capability overview item %d is incomplete", index)
		}
		if _, ok := seen[capability.ID]; ok {
			return fmt.Errorf("capability overview repeats capability %q", capability.ID)
		}
		seen[capability.ID] = struct{}{}
		switch capability.Status {
		case StatusAvailable:
			if len(capability.MissingSignals) != 0 {
				return fmt.Errorf("available overview capability %q retains missing signals", capability.ID)
			}
		case StatusDeferred:
			if len(capability.MissingSignals) == 0 {
				return fmt.Errorf("deferred overview capability %q has no missing signals", capability.ID)
			}
		default:
			return fmt.Errorf("overview capability %q has invalid status %q", capability.ID, capability.Status)
		}
	}
	if o.digest() != o.EvidenceDigest {
		return fmt.Errorf("capability overview evidence digest does not match")
	}
	return nil
}

func (o Overview) digest() string {
	parts := []string{
		"capability-overview",
		o.Query,
		o.DeclarationDigest,
		string(o.Mode),
		o.Summary,
		strings.Join(o.SuggestedQuestions, "|"),
		strings.Join(o.Constraints, "|"),
	}
	for _, capability := range o.Capabilities {
		parts = append(parts,
			capability.ID,
			capability.Label,
			string(capability.Status),
			strings.Join(capability.RequiredSignals, ","),
			strings.Join(capability.MissingSignals, ","),
		)
	}
	return digest(parts...)
}

func overviewCapabilitiesFromOptions(options []Option) []OverviewCapability {
	result := make([]OverviewCapability, 0, len(options))
	for _, option := range options {
		result = append(result, OverviewCapability{
			ID:              option.ID,
			Label:           option.Label,
			Status:          option.Status,
			RequiredSignals: append([]string(nil), option.RequiredSignals...),
			MissingSignals:  append([]string(nil), option.MissingSignals...),
		})
	}
	return result
}

func overviewCapabilitiesFromCatalog(matchedSignals []string) []OverviewCapability {
	result := make([]OverviewCapability, 0, len(catalog))
	for _, item := range catalog {
		missing := make([]string, 0, len(item.RequiredSignals))
		for _, required := range item.RequiredSignals {
			if !contains(matchedSignals, required) {
				missing = append(missing, required)
			}
		}
		status := StatusAvailable
		if len(missing) > 0 {
			status = StatusDeferred
		}
		result = append(result, OverviewCapability{
			ID:              item.ID,
			Label:           item.Label,
			Status:          status,
			RequiredSignals: append([]string(nil), item.RequiredSignals...),
			MissingSignals:  missing,
		})
	}
	return result
}

func questionsForCapabilities(capabilities []OverviewCapability) []string {
	questions := make([]string, 0, 3)
	for _, capability := range capabilities {
		for _, signal := range capability.MissingSignals {
			question := fmt.Sprintf("What declaration or evidence is needed for %s (%s)?", capability.Label, signal)
			if !contains(questions, question) {
				questions = append(questions, question)
			}
			if len(questions) == 3 {
				return questions
			}
		}
	}
	if len(questions) == 0 {
		return []string{"Which verified evidence should be collected for the matched capability?"}
	}
	return questions
}

func overviewQuestions() []string {
	questions := make([]string, 0, 3)
	for _, item := range catalog {
		questions = append(questions, fmt.Sprintf("Can gooo describe %s?", item.Label))
		if len(questions) == 3 {
			break
		}
	}
	return questions
}

func isOverviewQuery(query string) bool {
	query = strings.ToLower(strings.TrimSpace(query))
	if strings.Contains(query, "what can this language do") || strings.Contains(query, "what can gooo do") || strings.Contains(query, "무엇을 할 수") {
		return !hasAnyTerm(termSet(query), "provenance", "evidence", "generate", "generation", "codegen", "reverse", "observation", "execute", "execution", "run", "lsp", "feedback")
	}
	terms := termSet(query)
	return hasAnyTerm(terms, "가능", "지원", "기능") && !hasAnyTerm(terms, "provenance", "evidence", "generate", "generation", "codegen", "reverse", "observation", "execute", "execution", "run", "lsp", "feedback")
}
