package output

import (
	"fmt"
	"os"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/diff"
)

type changelogProfile struct {
}

func (p changelogProfile) format(format string, diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) (diffOutput []byte, errCode int, err error) {
	checks := checker.GetAllChecks()
	checks.Localizer = *localizations.New("en", "en")

	errs := checker.BackwardCompatibilityErrors{}
	errs = append(errs, checkAddedPaths(diffReport, checks, operationsSources)...)
	errs = append(errs, checkDeletedPaths(diffReport, checks, operationsSources)...)
	errs = append(errs, checkModifiedPaths(diffReport, checks, operationsSources)...)

	if diffOutput, _, err = marshalOutput(format, errs, false); err != nil {
		return nil, summaryGenerationErrCode, fmt.Errorf("failed to print summary in %s format with %w", format, err)
	}
	return diffOutput, noErrCode, nil
}

func checkAddedPaths(diffReport *diff.Diff, checks checker.BackwardCompatibilityCheckConfig, operationsSources *diff.OperationsSourcesMap) checker.BackwardCompatibilityErrors {
	errs := checker.BackwardCompatibilityErrors{}
	for _, path := range diffReport.PathsDiff.Added {
		dr := &diff.Diff{
			PathsDiff: &diff.PathsDiff{
				Added:    []string{path},
				Base:     diffReport.PathsDiff.Base,
				Revision: diffReport.PathsDiff.Revision,
			},
		}
		errs = append(errs, check(dr, checks, operationsSources)...)
	}
	return errs
}

func checkDeletedPaths(diffReport *diff.Diff, checks checker.BackwardCompatibilityCheckConfig, operationsSources *diff.OperationsSourcesMap) checker.BackwardCompatibilityErrors {
	errs := checker.BackwardCompatibilityErrors{}
	for _, path := range diffReport.PathsDiff.Deleted {
		dr := &diff.Diff{
			PathsDiff: &diff.PathsDiff{
				Deleted:  []string{path},
				Base:     diffReport.PathsDiff.Base,
				Revision: diffReport.PathsDiff.Revision,
			},
		}
		errs = append(errs, check(dr, checks, operationsSources)...)
	}
	return errs
}

func checkModifiedPaths(diffReport *diff.Diff, checks checker.BackwardCompatibilityCheckConfig, operationsSources *diff.OperationsSourcesMap) checker.BackwardCompatibilityErrors {
	errs := checker.BackwardCompatibilityErrors{}
	for p, d := range diffReport.PathsDiff.Modified {
		dr := &diff.Diff{
			PathsDiff: &diff.PathsDiff{
				Modified: map[string]*diff.PathDiff{p: d},
				Base:     diffReport.PathsDiff.Base,
				Revision: diffReport.PathsDiff.Revision,
			},
		}

		errs = append(errs, check(dr, checks, operationsSources)...)
	}
	return errs
}

func check(dr *diff.Diff, checks checker.BackwardCompatibilityCheckConfig, operationsSources *diff.OperationsSourcesMap) checker.BackwardCompatibilityErrors {
	pathErrs := checker.CheckBackwardCompatibility(checks, dr, operationsSources)
	if pathErrs.Len() == 0 {
		s, _ := ToJSON(dr)
		fmt.Printf("no oasdiff checker matched:\n%s", s)
		os.Exit(missingCheckerErrCode)
	}
	return pathErrs
}
