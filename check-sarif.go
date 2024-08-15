package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

type Sarif struct {
	Schema  string `json:"$schema"`
	Version string `json:"version"`
	Runs    []Run  `json:"runs"`
}

type Run struct {
	Tool                     Tool                       `json:"tool"`
	VersionControlProvenance []VersionControlProvenance `json:"versionControlProvenance"`
	Artifacts                []Artifact                 `json:"artifacts"`
	Results                  []Result                   `json:"results"`
	NewlineSequences         []string                   `json:"newlineSequences"`
	ColumnKind               string                     `json:"columnKind"`
	Properties               Properties                 `json:"properties"`
}

type Tool struct {
	Driver Driver `json:"driver"`
}

type Driver struct {
	Name         string `json:"name"`
	Organization string `json:"organization"`
	Version      string `json:"version"`
	Rules        []Rule `json:"rules"`
}

type Rule struct {
	ID                   string               `json:"id"`
	Name                 string               `json:"name"`
	ShortDescription     Description          `json:"shortDescription"`
	FullDescription      Description          `json:"fullDescription"`
	DefaultConfiguration DefaultConfiguration `json:"defaultConfiguration"`
	Properties           RuleProperties       `json:"properties"`
}

type Description struct {
	Text string `json:"text"`
}

type DefaultConfiguration struct {
	Enabled bool   `json:"enabled"`
	Level   string `json:"level"`
}

type RuleProperties struct {
	Tags      []string `json:"tags"`
	Kind      string   `json:"kind"`
	Precision string   `json:"precision"`
	Severity  string   `json:"severity"`
}

type VersionControlProvenance struct {
	RepositoryUri string `json:"repositoryUri"`
	RevisionId    string `json:"revisionId"`
}

type Artifact struct {
	Location ArtifactLocation `json:"location"`
}

type ArtifactLocation struct {
	Uri       string `json:"uri"`
	UriBaseId string `json:"uriBaseId"`
	Index     int    `json:"index"`
}

type Result struct {
	RuleId              string            `json:"ruleId"`
	RuleIndex           int               `json:"ruleIndex"`
	Rule                RuleReference     `json:"rule"`
	Message             Description       `json:"message"`
	Locations           []ResultLocation  `json:"locations"`
	PartialFingerprints map[string]string `json:"partialFingerprints"`
	RelatedLocations    []RelatedLocation `json:"relatedLocations,omitempty"`
	CodeFlows           []CodeFlow        `json:"codeFlows,omitempty"`
}

type RuleReference struct {
	ID    string `json:"id"`
	Index int    `json:"index"`
}

type ResultLocation struct {
	PhysicalLocation PhysicalLocation `json:"physicalLocation"`
}

type PhysicalLocation struct {
	ArtifactLocation ArtifactLocation `json:"artifactLocation"`
	Region           Region           `json:"region"`
}

type Region struct {
	StartLine   int `json:"startLine"`
	StartColumn int `json:"startColumn"`
	EndColumn   int `json:"endColumn,omitempty"`
}

type RelatedLocation struct {
	ID               int              `json:"id"`
	PhysicalLocation PhysicalLocation `json:"physicalLocation"`
	Message          Description      `json:"message"`
}

type CodeFlow struct {
	ThreadFlows []ThreadFlow `json:"threadFlows"`
}

type ThreadFlow struct {
	Locations []ThreadFlowLocation `json:"locations"`
}

type ThreadFlowLocation struct {
	Location PhysicalLocation `json:"location"`
	Message  Description      `json:"message"`
}

type Properties struct {
	SemmleFormatSpecifier string `json:"semmle.formatSpecifier"`
	SemmleSourceLanguage  string `json:"semmle.sourceLanguage"`
}

func main() {
	// Check arguments
	if len(os.Args) < 2 {
		fmt.Println("Please provide a JSON file name as an argument.")
		os.Exit(1)
	}
	fileName := os.Args[1]

	// Read the JSON file
	sarifData, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	decoder := json.NewDecoder(strings.NewReader(string(sarifData)))

	// Disallow unknown fields in the JSON and decode the JSON into a typed struct
	decoder.DisallowUnknownFields()
	var sarif Sarif
	err = decoder.Decode(&sarif)
	if err != nil {
		slog.Error("failed to unmarshal SARIF", "error", err)
	} else {
		slog.Info("Successfully parsed JSON into Sarif struct.")
	}
	os.Exit(0)

}
