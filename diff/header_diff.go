package diff

import "github.com/getkin/kin-openapi/openapi3"

// HeaderDiff is a diff between header objects: https://swagger.io/specification/#header-object
type HeaderDiff struct {
	ExtensionsDiff  *ExtensionsDiff `json:"extensions,omitempty"`
	DescriptionDiff *ValueDiff      `json:"description,omitempty"`
	DeprecatedDiff  *ValueDiff      `json:"deprecated,omitempty"`
	RequiredDiff    *ValueDiff      `json:"required,omitempty"`
	ExampleDiff     *ValueDiff      `json:"example,omitempty"`
	// Examples
	SchemaDiff  *SchemaDiff  `json:"schema,omitempty"`
	ContentDiff *ContentDiff `json:"content,omitempty"`
}

func (headerDiff *HeaderDiff) empty() bool {
	return headerDiff == nil || *headerDiff == HeaderDiff{}
}

func getHeaderDiff(config *Config, header1, header2 *openapi3.Header) *HeaderDiff {
	diff := getHeaderDiffInternal(config, header1, header2)
	if diff.empty() {
		return nil
	}
	return diff
}

func getHeaderDiffInternal(config *Config, header1, header2 *openapi3.Header) *HeaderDiff {
	result := HeaderDiff{}

	result.ExtensionsDiff = getExtensionsDiff(config, header1.ExtensionProps, header2.ExtensionProps)
	result.DescriptionDiff = getValueDiff(header1.Description, header2.Description)
	result.DeprecatedDiff = getValueDiff(header1.Deprecated, header2.Deprecated)
	result.RequiredDiff = getValueDiff(header1.Required, header2.Required)
	result.SchemaDiff = getSchemaDiff(config, header1.Schema, header2.Schema)

	if config.IncludeExamples {
		result.ExampleDiff = getValueDiff(header1.Example, header2.Example)
	}

	result.ContentDiff = getContentDiff(config, header1.Content, header2.Content)

	return &result
}
