package capability

import (
	"fmt"
	"sort"
	"strings"
)

type AssistantObservationItem struct {
	Path      string    `json:"path"`
	Assistant Assistant `json:"assistant"`
}

type ObservationCounts struct {
	Total     int `json:"total"`
	Available int `json:"available"`
	Deferred  int `json:"deferred"`
	Unknown   int `json:"unknown"`
}

type AssistantObservationReport struct {
	Query          string                    `json:"query"`
	Files          []AssistantObservationItem `json:"files"`
	Counts         ObservationCounts        `json:"counts"`
	Constraints    []string                  `json:"constraints"`
	EvidenceDigest string                    `json:"evidence_digest"`
}

var observationReportConstraints = []string{
	"observed_scope_only",
	"not_a_language_completeness_claim",
	"unknown_is_preserved",
	"no_provider_invocation",
	"no_authorization_grant",
	"no_catalog_mutation",
	"cache_presence_is_not_semantic_evidence",
}

func BuildAssistantObservationReport(query string, files []AssistantObservationItem) (AssistantObservationReport, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return AssistantObservationReport{}, fmt.Errorf("assistant observation report query is required")
	}
	if len(files) == 0 {
		return AssistantObservationReport{}, fmt.Errorf("assistant observation report requires at least one file")
	}

	result := AssistantObservationReport{
		Query:       query,
		Files:       append([]AssistantObservationItem(nil), files...),
		Constraints: append([]string(nil), observationReportConstraints...),
	}
	if err := result.recalculate(); err != nil {
		return AssistantObservationReport{}, err
	}
	result.EvidenceDigest = result.digest()
	return result, nil
}

func (report AssistantObservationReport) Validate() error {
	if strings.TrimSpace(report.Query) == "" || len(report.Files) == 0 {
		return fmt.Errorf("assistant observation report is incomplete")
	}
	if len(report.Constraints) != len(observationReportConstraints) || !sameStrings(report.Constraints, observationReportConstraints) {
		return fmt.Errorf("assistant observation report constraints are incomplete")
	}
	if !validDigest(report.EvidenceDigest) {
		return fmt.Errorf("assistant observation report evidence digest is invalid")
	}
	recalculated := report
	if err := recalculated.recalculate(); err != nil {
		return err
	}
	if recalculated.Counts != report.Counts {
		return fmt.Errorf("assistant observation report counts do not match files")
	}
	if report.digest() != report.EvidenceDigest {
		return fmt.Errorf("assistant observation report evidence digest does not match")
	}
	return nil
}

func (report *AssistantObservationReport) recalculate() error {
	seen := make(map[string]struct{}, len(report.Files))
	report.Counts = ObservationCounts{Total: len(report.Files)}
	for index := range report.Files {
		item := &report.Files[index]
		item.Path = strings.TrimSpace(item.Path)
		if item.Path == "" {
			return fmt.Errorf("assistant observation report file %d has no path", index)
		}
		if _, ok := seen[item.Path]; ok {
			return fmt.Errorf("assistant observation report repeats path %q", item.Path)
		}
		seen[item.Path] = struct{}{}
		if err := item.Assistant.Validate(); err != nil {
			return fmt.Errorf("validate assistant for %s: %w", item.Path, err)
		}
		if item.Assistant.Query != report.Query {
			return fmt.Errorf("assistant query for %s does not match report query", item.Path)
		}
		switch item.Assistant.Discovery.Status {
		case StatusAvailable:
			report.Counts.Available++
		case StatusDeferred:
			report.Counts.Deferred++
		case StatusUnknown:
			report.Counts.Unknown++
		default:
			return fmt.Errorf("assistant status for %s is invalid", item.Path)
		}
	}
	return nil
}

func (report AssistantObservationReport) digest() string {
	parts := []string{
		"capability-assistant-observation-report",
		report.Query,
		fmt.Sprintf("%d", report.Counts.Total),
		fmt.Sprintf("%d", report.Counts.Available),
		fmt.Sprintf("%d", report.Counts.Deferred),
		fmt.Sprintf("%d", report.Counts.Unknown),
		strings.Join(report.Constraints, "|"),
	}
	paths := make([]string, 0, len(report.Files))
	for _, item := range report.Files {
		paths = append(paths, item.Path+"|"+item.Assistant.EvidenceDigest)
	}
	sort.Strings(paths)
	return digest(append(parts, paths...)...)
}
