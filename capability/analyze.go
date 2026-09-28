package capability

import "fmt"

type SourceAnalysis struct {
	DeclarationShape DeclarationShape `json:"declaration_shape"`
	Discovery        Discovery        `json:"discovery"`
	Plan             Plan             `json:"plan"`
	AnalysisDigest   string           `json:"analysis_digest"`
}

func AnalyzeSource(query, source string) (SourceAnalysis, error) {
	shape, err := InspectDeclaration(source)
	if err != nil {
		return SourceAnalysis{}, err
	}
	discovery, err := DiscoverFromSource(query, source)
	if err != nil {
		return SourceAnalysis{}, err
	}
	plan, err := BuildPlan(discovery)
	if err != nil {
		return SourceAnalysis{}, err
	}
	analysis := SourceAnalysis{
		DeclarationShape: shape,
		Discovery:        discovery,
		Plan:             plan,
	}
	analysis.AnalysisDigest = analysis.digest()
	return analysis, nil
}

func (a SourceAnalysis) Validate() error {
	if err := a.DeclarationShape.Validate(); err != nil {
		return err
	}
	if err := a.Discovery.Validate(); err != nil {
		return err
	}
	if err := a.Plan.Validate(); err != nil {
		return err
	}
	if a.AnalysisDigest == "" || !validDigest(a.AnalysisDigest) {
		return fmt.Errorf("source analysis digest is invalid")
	}
	if a.digest() != a.AnalysisDigest {
		return fmt.Errorf("source analysis digest does not match")
	}
	return nil
}

func (a SourceAnalysis) digest() string {
	return digest(
		"gooo-source-analysis",
		a.DeclarationShape.ShapeDigest,
		a.Discovery.EvidenceDigest,
		a.Plan.PlanDigest,
	)
}
