package report

import (
	"bytes"

	"github.com/tufin/oasdiff/diff"
)

// GetTextReportAsString returns a textual diff report as a string
// The report is compatible with Github markdown
func GetTextReportAsString(d *diff.Diff) string {
	return output(d).String()
}

// GetTextReportAsBytes returns a textual diff report as bytes
// The report is compatible with Github markdown
func GetTextReportAsBytes(d *diff.Diff) []byte {
	return output(d).Bytes()
}

func output(d *diff.Diff) *bytes.Buffer {
	var buf bytes.Buffer
	r := report{
		Writer: &buf,
	}
	r.output(d)
	return &buf
}
