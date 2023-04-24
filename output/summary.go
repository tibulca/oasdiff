package output

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

type summaryProfile struct {
}

func (p summaryProfile) format(format string, diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) (diffOutput []byte, errCode int, err error) {
	if diffOutput, _, err = marshalOutput(format, diffReport.GetSummary(), false); err != nil {
		return nil, summaryGenerationErrCode, fmt.Errorf("failed to print summary in %s format with %w", format, err)
	}
	return diffOutput, noErrCode, nil
}
