package capability

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"unicode"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

type Status string

const (
	StatusAvailable Status = "AVAILABLE"
	StatusDeferred  Status = "DEFERRED"
	StatusUnknown   Status = "UNKNOWN"
)

type Descriptor struct {
	ID              string
	Label           string
	RequiredSignals []string
	Aliases         []string
}

type Discovery struct {
	Query            string `json:"query"`
	DeclarationDigest string `json:"declaration_digest"`
	CapabilityID      string `json:"capability_id"`
	CapabilityLabel   string `json:"capability_label"`
	Status            Status `json:"status"`
	MatchedSignals    []string `json:"matched_signals"`
	MissingSignals    []string `json:"missing_signals"`
	NextQuestions     []string `json:"next_questions"`
	EvidenceDigest    string `json:"evidence_digest"`
}

type capabilitySpec struct {
	Descriptor
	questions map[string]string
}

var catalog = []capabilitySpec{
	{
		Descriptor: Descriptor{
			ID:              "capability-discovery",
			Label:           "discover-what-gooo-can-do",
			RequiredSignals: []string{"activity", "evidence"},
			Aliases:         []string{"what can", "what can this language do", "capability", "discover", "possible", "지원", "기능", "무엇을 할 수"},
		},
		questions: map[string]string{
			"activity": "어떤 .gooo activity가 이 기능을 표현해야 하나요?",
			"evidence": "어떤 digest 또는 관찰 근거가 결과를 뒷받침하나요?",
		},
	},
	{
		Descriptor: Descriptor{
			ID:              "execution-observation",
			Label:           "observe-an-execution-boundary",
			RequiredSignals: []string{"execution", "receipt"},
			Aliases:         []string{"execute", "execution", "run", "receipt", "실행", "실행 결과"},
		},
		questions: map[string]string{
			"execution": "실행 경계와 비실행 관찰 경계를 어떻게 선언하나요?",
			"receipt":   "terminal result와 receipt digest가 있나요?",
		},
	},
	{
		Descriptor: Descriptor{
			ID:              "reverse-observation",
			Label:           "trace-result-back-to-evidence",
			RequiredSignals: []string{"reverse", "evidence"},
			Aliases:         []string{"reverse", "reverse observation", "observe", "provenance", "origin", "evidence", "역관찰", "기원"},
		},
		questions: map[string]string{
			"reverse":  "역관찰 단계와 최초 missing stage를 선언했나요?",
			"evidence": "관찰 결과와 verifier digest를 연결했나요?",
		},
	},
	{
		Descriptor: Descriptor{
			ID:              "code-generation",
			Label:           "generate-an-implementation-boundary",
			RequiredSignals: []string{"generation", "output"},
			Aliases:         []string{"generate", "generation", "codegen", "code generation", "생성", "코드 생성"},
		},
		questions: map[string]string{
			"generation": "생성 단계와 생성기 identity를 선언했나요?",
			"output":     "생성된 output digest를 보존하나요?",
		},
	},
	{
		Descriptor: Descriptor{
			ID:              "lsp-feedback",
			Label:           "turn-feedback-into-language-guidance",
			RequiredSignals: []string{"lsp", "feedback"},
			Aliases:         []string{"lsp", "language server", "feedback", "diagnostic", "언어 서버", "피드백"},
		},
		questions: map[string]string{
			"lsp":      "LSP 응답이 어떤 declaration과 연결되나요?",
			"feedback": "피드백의 source와 disposition을 보존하나요?",
		},
	},
}

func Catalog() []Descriptor {
	result := make([]Descriptor, 0, len(catalog))
	for _, item := range catalog {
		result = append(result, Descriptor{
			ID:              item.ID,
			Label:           item.Label,
			RequiredSignals: append([]string(nil), item.RequiredSignals...),
			Aliases:         append([]string(nil), item.Aliases...),
		})
	}
	return result
}

func Discover(query string, declaration envelope.Declaration) (Discovery, error) {
	if err := declaration.Validate(); err != nil {
		return Discovery{}, err
	}
	query = strings.TrimSpace(query)
	if query == "" {
		return Discovery{}, fmt.Errorf("capability discovery query is required")
	}

	result := Discovery{
		Query:             query,
		DeclarationDigest: declaration.Digest,
		Status:            StatusUnknown,
		NextQuestions:     []string{"어떤 capability family를 찾고 있나요?"},
	}

	queryLower := strings.ToLower(query)
	queryTerms := termSet(queryLower)
	var selected *capabilitySpec
	bestScore := 0
	for index := range catalog {
		score := 0
		for _, alias := range catalog[index].Aliases {
			score = max(score, aliasScore(queryLower, queryTerms, alias))
		}
		if score > bestScore {
			selected = &catalog[index]
			bestScore = score
		}
	}
	if selected == nil {
		result.EvidenceDigest = result.digest()
		return result, nil
	}

	result.CapabilityID = selected.ID
	result.CapabilityLabel = selected.Label
	result.MatchedSignals = declarationSignals(declaration.Source)
	for _, required := range selected.RequiredSignals {
		if !contains(result.MatchedSignals, required) {
			result.MissingSignals = append(result.MissingSignals, required)
		}
	}
	sort.Strings(result.MatchedSignals)
	sort.Strings(result.MissingSignals)
	if len(result.MissingSignals) == 0 {
		result.Status = StatusAvailable
		result.NextQuestions = nil
	} else {
		result.Status = StatusDeferred
		result.NextQuestions = make([]string, 0, len(result.MissingSignals))
		for _, signal := range result.MissingSignals {
			result.NextQuestions = append(result.NextQuestions, selected.questions[signal])
		}
	}
	result.EvidenceDigest = result.digest()
	return result, nil
}

func (d Discovery) Validate() error {
	if strings.TrimSpace(d.Query) == "" || !validDigest(d.DeclarationDigest) || !validDigest(d.EvidenceDigest) {
		return fmt.Errorf("capability discovery evidence is incomplete")
	}
	switch d.Status {
	case StatusAvailable:
		if d.CapabilityID == "" || len(d.MissingSignals) != 0 || len(d.NextQuestions) != 0 {
			return fmt.Errorf("available capability discovery is incomplete")
		}
	case StatusDeferred:
		if d.CapabilityID == "" || len(d.MissingSignals) == 0 || len(d.NextQuestions) == 0 {
			return fmt.Errorf("deferred capability discovery must preserve missing signals")
		}
	case StatusUnknown:
		if d.CapabilityID != "" || len(d.NextQuestions) == 0 {
			return fmt.Errorf("unknown capability discovery must preserve next question")
		}
	default:
		return fmt.Errorf("capability discovery status %q is invalid", d.Status)
	}
	if d.digest() != d.EvidenceDigest {
		return fmt.Errorf("capability discovery evidence digest does not match")
	}
	return nil
}

func (d Discovery) digest() string {
	return digest(
		"capability-discovery",
		d.Query,
		d.DeclarationDigest,
		d.CapabilityID,
		d.CapabilityLabel,
		string(d.Status),
		strings.Join(d.MatchedSignals, ","),
		strings.Join(d.MissingSignals, ","),
		strings.Join(d.NextQuestions, "|"),
	)
}

func declarationSignals(source string) []string {
	terms := termSet(strings.ToLower(source))
	signals := make([]string, 0, 8)
	if hasTerm(terms, "activity") {
		signals = append(signals, "activity")
	}
	if hasAnyTerm(terms, "evidence", "provenance", "digest") {
		signals = append(signals, "evidence")
	}
	if hasAnyTerm(terms, "execute", "execution", "run") {
		signals = append(signals, "execution")
	}
	if hasAnyTerm(terms, "receipt") {
		signals = append(signals, "receipt")
	}
	if hasAnyTerm(terms, "reverse", "observation", "observe") {
		signals = append(signals, "reverse")
	}
	if hasAnyTerm(terms, "generate", "generation", "codegen") {
		signals = append(signals, "generation")
	}
	if hasAnyTerm(terms, "output", "result") {
		signals = append(signals, "output")
	}
	if hasAnyTerm(terms, "lsp", "language") {
		signals = append(signals, "lsp")
	}
	if hasAnyTerm(terms, "feedback", "diagnostic") {
		signals = append(signals, "feedback")
	}
	return signals
}

func aliasScore(query string, queryTerms map[string]struct{}, alias string) int {
	alias = strings.ToLower(strings.TrimSpace(alias))
	if alias == "" {
		return 0
	}
	if strings.Contains(query, alias) {
		return 100 + len(alias)
	}
	aliasTerms := terms(alias)
	for _, term := range aliasTerms {
		if _, ok := queryTerms[term]; !ok {
			return 0
		}
	}
	return len(aliasTerms)
}

func termSet(value string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, term := range terms(value) {
		result[term] = struct{}{}
	}
	return result
}

func terms(value string) []string {
	return strings.FieldsFunc(strings.ToLower(value), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
}

func hasTerm(set map[string]struct{}, value string) bool {
	_, ok := set[value]
	return ok
}

func hasAnyTerm(set map[string]struct{}, values ...string) bool {
	for _, value := range values {
		if hasTerm(set, value) {
			return true
		}
	}
	return false
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func digest(parts ...string) string {
	sum := sha256.Sum256([]byte(strings.Join(parts, "|")))
	return hex.EncodeToString(sum[:])
}

func validDigest(value string) bool {
	if len(value) != 64 {
		return false
	}
	_, err := hex.DecodeString(value)
	return err == nil
}
