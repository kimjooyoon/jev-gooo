package capability

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

type Option struct {
	ID              string `json:"id"`
	Label           string `json:"label"`
	Score           int    `json:"score"`
	Status          Status `json:"status"`
	RequiredSignals []string `json:"required_signals"`
	MissingSignals  []string `json:"missing_signals,omitempty"`
}

type Options struct {
	Query            string   `json:"query"`
	DeclarationDigest string   `json:"declaration_digest"`
	Options          []Option `json:"options"`
	EvidenceDigest   string   `json:"evidence_digest"`
}

func DiscoverOptions(query string, declaration envelope.Declaration) (Options, error) {
	if err := declaration.Validate(); err != nil {
		return Options{}, err
	}
	query = strings.TrimSpace(query)
	if query == "" {
		return Options{}, fmt.Errorf("capability options query is required")
	}

	queryLower := strings.ToLower(query)
	queryTerms := termSet(queryLower)
	matchedSignals := declarationSignals(declaration.Source)
	result := Options{Query: query, DeclarationDigest: declaration.Digest}
	for _, item := range catalog {
		score := 0
		for _, alias := range item.Aliases {
			score = max(score, aliasScore(queryLower, queryTerms, alias))
		}
		if score == 0 {
			continue
		}
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
		result.Options = append(result.Options, Option{
			ID:              item.ID,
			Label:           item.Label,
			Score:           score,
			Status:          status,
			RequiredSignals: append([]string(nil), item.RequiredSignals...),
			MissingSignals:  missing,
		})
	}
	sort.Slice(result.Options, func(left, right int) bool {
		if result.Options[left].Score != result.Options[right].Score {
			return result.Options[left].Score > result.Options[right].Score
		}
		return result.Options[left].ID < result.Options[right].ID
	})
	result.EvidenceDigest = result.digest()
	return result, nil
}

func (o Options) Validate() error {
	if strings.TrimSpace(o.Query) == "" || !validDigest(o.DeclarationDigest) || !validDigest(o.EvidenceDigest) {
		return fmt.Errorf("capability options evidence is incomplete")
	}
	previousScore := int(^uint(0) >> 1)
	previousID := ""
	for index, option := range o.Options {
		if strings.TrimSpace(option.ID) == "" || strings.TrimSpace(option.Label) == "" || option.Score <= 0 {
			return fmt.Errorf("capability option %d is incomplete", index)
		}
		if index > 0 && (option.Score > previousScore || option.Score == previousScore && option.ID <= previousID) {
			return fmt.Errorf("capability options are not deterministically ordered")
		}
		if len(option.RequiredSignals) == 0 {
			return fmt.Errorf("capability option %q has no required signals", option.ID)
		}
		switch option.Status {
		case StatusAvailable:
			if len(option.MissingSignals) != 0 {
				return fmt.Errorf("available capability option %q retains missing signals", option.ID)
			}
		case StatusDeferred:
			if len(option.MissingSignals) == 0 {
				return fmt.Errorf("deferred capability option %q has no missing signals", option.ID)
			}
		default:
			return fmt.Errorf("capability option %q has invalid status %q", option.ID, option.Status)
		}
		previousScore = option.Score
		previousID = option.ID
	}
	if o.digest() != o.EvidenceDigest {
		return fmt.Errorf("capability options evidence digest does not match")
	}
	return nil
}

func (o Options) digest() string {
	parts := []string{"capability-options", o.Query, o.DeclarationDigest}
	for _, option := range o.Options {
		parts = append(parts,
			option.ID,
			option.Label,
			strconv.Itoa(option.Score),
			string(option.Status),
			strings.Join(option.RequiredSignals, ","),
			strings.Join(option.MissingSignals, ","),
		)
	}
	return digest(parts...)
}