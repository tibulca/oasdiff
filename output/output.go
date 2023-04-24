package output

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/report"
	"gopkg.in/yaml.v3"
)

type OutputFormat string

const (
	YAML OutputFormat = "yaml"
	JSON OutputFormat = "json"
	Text OutputFormat = "text"
	HTML OutputFormat = "html"
)

type OutputProfile string

const (
	DefaultProfile   OutputProfile = "default"
	SummaryProfile   OutputProfile = "summary"
	ChangelogProfile OutputProfile = "changelog"
)

const (
	noErrCode                = 0
	summaryGenerationErrCode = 105
	yamlMarshalErrCode       = 106
	jsonMarshalErrCode       = 106
	htmlGenerationErrCode    = 107
	unknownFormatErrCode     = 108
	missingCheckerErrCode    = 109
)

type Profile interface {
	format(format string, diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) (diffOutput []byte, errCode int, err error)
}

// todo: lang param
func Format(format string, profile string, diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) (diffOutput []byte, errCode int, err error) {
	var outputProfile Profile
	switch OutputProfile(profile) {
	case DefaultProfile:
		outputProfile = defaultProfile{}
	case SummaryProfile:
		outputProfile = summaryProfile{}
	case ChangelogProfile:
		outputProfile = changelogProfile{}
	default:
		return nil, unknownFormatErrCode, fmt.Errorf("unknown output profile %s", profile)
	}

	return outputProfile.format(format, diffReport, operationsSources)
}

func marshalOutput(format string, data interface{}, wrapError bool) (diffOutput []byte, errCode int, err error) {
	switch OutputFormat(format) {
	case YAML:
		if diffOutput, err = ToYAML(data); err != nil {
			return nil, yamlMarshalErrCode, fmt.Errorf("failed to print diff YAML with %w", err)
		}
	case JSON:
		if diffOutput, err = ToJSON(data); err != nil {
			return nil, jsonMarshalErrCode, fmt.Errorf("failed to print diff JSON with %w", err)
		}
	case Text:
		diffReport, ok := data.(*diff.Diff)
		if !ok {
			// todo: implement text format for summary and changelog profiles
			return nil, unknownFormatErrCode, fmt.Errorf("output format %s not implemented for selected profile", format)
		}
		return report.GetTextReportAsBytes(diffReport), noErrCode, nil
	case HTML:
		diffReport, ok := data.(*diff.Diff)
		if !ok {
			// todo: implement text format for summary and changelog profiles
			return nil, unknownFormatErrCode, fmt.Errorf("output format %s not implemented for selected profile", format)
		}
		if diffOutput, err = report.GetHTMLReportAsBytes(diffReport); err != nil {
			return nil, htmlGenerationErrCode, fmt.Errorf("failed to generate HTML diff report with %w", err)
		}
	default:
		return nil, unknownFormatErrCode, fmt.Errorf("unknown output format %q", format)
	}
	return diffOutput, noErrCode, nil
}

func ToYAML(output interface{}) ([]byte, error) {
	if reflect.ValueOf(output).IsNil() {
		return nil, nil
	}
	return yaml.Marshal(output)
}

func ToJSON(output interface{}) ([]byte, error) {
	if reflect.ValueOf(output).IsNil() {
		return nil, nil
	}
	return json.MarshalIndent(output, "", "  ")
}
