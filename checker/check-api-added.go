package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func APIAddedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	for _, path := range diffReport.PathsDiff.Added {
		for operation := range diffReport.PathsDiff.Revision[path].Operations() {
			op := diffReport.PathsDiff.Revision[path].Operations()[operation]
			source := (*operationsSources)[op]
			result = append(result, BackwardCompatibilityError{
				Id:        "api-path-added",
				Level:     INFO,
				Text:      config.i18n("api-path-added"),
				Operation: operation,
				Path:      path,
				Source:    source,
			})
		}
	}

	return result
}
