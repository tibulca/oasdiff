package output

import (
	"github.com/tufin/oasdiff/diff"
)

type defaultProfile struct {
}

func (p defaultProfile) format(format string, diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) (diffOutput []byte, errCode int, err error) {
	return marshalOutput(format, diffReport, true)
}
