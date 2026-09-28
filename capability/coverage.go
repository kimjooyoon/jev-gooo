package capability

import (
	"fmt"
	"strconv"
)

const (
	CoverageScope = "explicit_catalog_match_surface_only"
	CoverageInterpretation = "not_language_completeness"
)

type Coverage struct {
	OptionsDigest  string `json:"options_digest"`
	CatalogSize    int    `json:"catalog_size"`
	MatchedCount   int    `json:"matched_count"`
	AvailableCount int    `json:"available_count"`
	DeferredCount  int    `json:"deferred_count"`
	UnknownQuery   bool   `json:"unknown_query"`
	Scope          string `json:"scope"`
	Interpretation string `json:"interpretation"`
	CoverageDigest string `json:"coverage_digest"`
}

func MeasureCoverage(options Options) (Coverage, error) {
	if err := options.Validate(); err != nil {
		return Coverage{}, err
	}
	result := Coverage{
		OptionsDigest:  options.EvidenceDigest,
		CatalogSize:    len(catalog),
		Scope:          CoverageScope,
		Interpretation: CoverageInterpretation,
	}
	for _, option := range options.Options {
		result.MatchedCount++
		switch option.Status {
		case StatusAvailable:
			result.AvailableCount++
		case StatusDeferred:
			result.DeferredCount++
		}
	}
	result.UnknownQuery = result.MatchedCount == 0
	result.CoverageDigest = result.digest()
	return result, nil
}

func (c Coverage) Validate() error {
	if !validDigest(c.OptionsDigest) || !validDigest(c.CoverageDigest) {
		return fmt.Errorf("capability coverage digest is invalid")
	}
	if c.CatalogSize != len(catalog) || c.CatalogSize < 0 {
		return fmt.Errorf("capability coverage catalog size is invalid")
	}
	if c.MatchedCount < 0 || c.AvailableCount < 0 || c.DeferredCount < 0 || c.MatchedCount != c.AvailableCount+c.DeferredCount || c.MatchedCount > c.CatalogSize {
		return fmt.Errorf("capability coverage counts are inconsistent")
	}
	if c.UnknownQuery != (c.MatchedCount == 0) {
		return fmt.Errorf("capability coverage unknown state is inconsistent")
	}
	if c.Scope != CoverageScope || c.Interpretation != CoverageInterpretation {
		return fmt.Errorf("capability coverage scope is invalid")
	}
	if c.digest() != c.CoverageDigest {
		return fmt.Errorf("capability coverage digest does not match")
	}
	return nil
}

func (c Coverage) digest() string {
	return digest(
		"capability-coverage",
		c.OptionsDigest,
		strconv.Itoa(c.CatalogSize),
		strconv.Itoa(c.MatchedCount),
		strconv.Itoa(c.AvailableCount),
		strconv.Itoa(c.DeferredCount),
		strconv.FormatBool(c.UnknownQuery),
		c.Scope,
		c.Interpretation,
	)
}