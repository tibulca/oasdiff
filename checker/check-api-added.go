package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func APIAddedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []CheckResult {
	result := make([]CheckResult, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	for _, path := range diffReport.PathsDiff.Added {
		for operation := range diffReport.PathsDiff.Revision[path].Operations() {
			op := diffReport.PathsDiff.Revision[path].Operations()[operation]
			source := (*operationsSources)[op]
			result = append(result, CheckResult{
				Id:          "api-path-added",
				Level:       INFO,
				Text:        config.i18n("api-path-added"),
				Operation:   operation,
				OperationId: op.OperationID,
				Path:        path,
				Source:      source,
			})
		}
	}

	return result
}
