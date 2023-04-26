package output

import (
	"fmt"

	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/report"
)

type changelogProfile struct {
	diffReport        *diff.Diff
	operationsSources *diff.OperationsSourcesMap
	checks            checker.BackwardCompatibilityCheckConfig
}

func (p *changelogProfile) format(format string, diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) (diffOutput []byte, errCode int, err error) {
	p.checks = checker.GetAllChecks()
	p.checks.Localizer = *localizations.New("en", "en")
	p.operationsSources = operationsSources
	p.diffReport = diffReport

	errs := checker.BackwardCompatibilityErrors{}
	errs = append(errs, p.checkAddedPaths()...)
	errs = append(errs, p.checkDeletedPaths()...)
	errs = append(errs, p.checkModifiedPaths()...)

	if diffOutput, _, err = marshalOutput(format, errs, false); err != nil {
		return nil, summaryGenerationErrCode, fmt.Errorf("failed to print summary in %s format with %w", format, err)
	}
	return diffOutput, noErrCode, nil
}

func (p *changelogProfile) checkAddedPaths() checker.BackwardCompatibilityErrors {
	errs := checker.BackwardCompatibilityErrors{}
	for _, path := range p.diffReport.PathsDiff.Added {
		dr := &diff.Diff{
			PathsDiff: &diff.PathsDiff{
				Added:    []string{path},
				Base:     p.diffReport.PathsDiff.Base,
				Revision: p.diffReport.PathsDiff.Revision,
			},
		}
		errs = append(errs, p.check(dr)...)
	}
	return errs
}

func (p *changelogProfile) checkDeletedPaths() checker.BackwardCompatibilityErrors {
	errs := checker.BackwardCompatibilityErrors{}
	for _, path := range p.diffReport.PathsDiff.Deleted {
		dr := &diff.Diff{
			PathsDiff: &diff.PathsDiff{
				Deleted:  []string{path},
				Base:     p.diffReport.PathsDiff.Base,
				Revision: p.diffReport.PathsDiff.Revision,
			},
		}
		errs = append(errs, p.check(dr)...)
	}
	return errs
}

func (p *changelogProfile) checkModifiedPaths() checker.BackwardCompatibilityErrors {
	errs := checker.BackwardCompatibilityErrors{}
	for path, pathDiff := range p.diffReport.PathsDiff.Modified {
		errs = append(errs, p.checkModifiedRefDiff(path, pathDiff)...)
		errs = append(errs, p.checkModifiedSummaryDiff(path, pathDiff)...)
		errs = append(errs, p.checkModifiedDescriptionDiff(path, pathDiff)...)
		errs = append(errs, p.checkModifiedOperationsDiff(path, pathDiff)...)

		// todo: pathDiff.ExtensionsDiff
		// todo: pathDiff.ServersDiff
		// todo: pathDiff.ParametersDiff

		// dr := &diff.Diff{
		// 	PathsDiff: &diff.PathsDiff{
		// 		Modified: map[string]*diff.PathDiff{path: pathDiff},
		// 		Base:     p.diffReport.PathsDiff.Base,
		// 		Revision: p.diffReport.PathsDiff.Revision,
		// 	},
		// }

		// errs = append(errs, p.check(dr)...)
	}
	return errs
}

func (p *changelogProfile) checkModifiedRefDiff(path string, pathDiff *diff.PathDiff) checker.BackwardCompatibilityErrors {
	if pathDiff.RefDiff == nil {
		return checker.BackwardCompatibilityErrors{}
	}
	return p.checkModifiedPathDiff(path, &diff.PathDiff{RefDiff: pathDiff.RefDiff})
}

func (p *changelogProfile) checkModifiedSummaryDiff(path string, pathDiff *diff.PathDiff) checker.BackwardCompatibilityErrors {
	if pathDiff.SummaryDiff == nil {
		return checker.BackwardCompatibilityErrors{}
	}
	return p.checkModifiedPathDiff(path, &diff.PathDiff{SummaryDiff: pathDiff.SummaryDiff})
}

func (p *changelogProfile) checkModifiedDescriptionDiff(path string, pathDiff *diff.PathDiff) checker.BackwardCompatibilityErrors {
	if pathDiff.DescriptionDiff == nil {
		return checker.BackwardCompatibilityErrors{}
	}
	return p.checkModifiedPathDiff(path, &diff.PathDiff{DescriptionDiff: pathDiff.DescriptionDiff})
}

func (p *changelogProfile) checkModifiedOperationsDiff(path string, pathDiff *diff.PathDiff) checker.BackwardCompatibilityErrors {
	if pathDiff.OperationsDiff == nil {
		return checker.BackwardCompatibilityErrors{}
	}

	errs := checker.BackwardCompatibilityErrors{}

	for _, op := range pathDiff.OperationsDiff.Added {
		errs = append(errs, p.checkModifiedPathDiff(path, &diff.PathDiff{
			OperationsDiff: &diff.OperationsDiff{
				Added: []string{op},
			},
		})...)
	}
	for _, op := range pathDiff.OperationsDiff.Deleted {
		errs = append(errs, p.checkModifiedPathDiff(path, &diff.PathDiff{
			OperationsDiff: &diff.OperationsDiff{
				Deleted: []string{op},
			},
		})...)
	}

	for op, opDiff := range pathDiff.OperationsDiff.Modified {
		/*
			ParametersDiff   *ParametersDiff           `json:"parameters,omitempty" yaml:"parameters,omitempty"`
			RequestBodyDiff  *RequestBodyDiff          `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
			ResponsesDiff    *ResponsesDiff            `json:"responses,omitempty" yaml:"responses,omitempty"`
			CallbacksDiff    *CallbacksDiff            `json:"callbacks,omitempty" yaml:"callbacks,omitempty"`
			SecurityDiff     *SecurityRequirementsDiff `json:"securityRequirements,omitempty" yaml:"securityRequirements,omitempty"`
			ServersDiff      *ServersDiff              `json:"servers,omitempty" yaml:"servers,omitempty"`
			ExternalDocsDiff *ExternalDocsDiff
		*/

		if opDiff.DeprecatedDiff != nil {
			errs = append(errs, p.checkModifiedPathDiff(path, &diff.PathDiff{
				OperationsDiff: &diff.OperationsDiff{
					Modified: map[string]*diff.MethodDiff{op: {DeprecatedDiff: opDiff.DeprecatedDiff}},
				},
			})...)
		}

		if opDiff.OperationIDDiff != nil {
			errs = append(errs, p.checkModifiedPathDiff(path, &diff.PathDiff{
				OperationsDiff: &diff.OperationsDiff{
					Modified: map[string]*diff.MethodDiff{op: {OperationIDDiff: opDiff.OperationIDDiff}},
				},
			})...)
		}

		if opDiff.TagsDiff != nil {
			errs = append(errs, p.checkModifiedPathDiff(path, &diff.PathDiff{
				OperationsDiff: &diff.OperationsDiff{
					Modified: map[string]*diff.MethodDiff{op: {TagsDiff: opDiff.TagsDiff}},
				},
			})...)
		}

		if opDiff.ParametersDiff != nil {
			if opDiff.ParametersDiff.Added != nil {
				errs = append(errs, p.checkModifiedPathDiff(path, &diff.PathDiff{
					OperationsDiff: &diff.OperationsDiff{
						Modified: map[string]*diff.MethodDiff{op: {ParametersDiff: &diff.ParametersDiff{Added: opDiff.ParametersDiff.Added}}},
					},
				})...)
			}

			if opDiff.ParametersDiff.Deleted != nil {
				errs = append(errs, p.checkModifiedPathDiff(path, &diff.PathDiff{
					OperationsDiff: &diff.OperationsDiff{
						Modified: map[string]*diff.MethodDiff{op: {ParametersDiff: &diff.ParametersDiff{Deleted: opDiff.ParametersDiff.Deleted}}},
					},
				})...)
			}

			for paramLocation, params := range opDiff.ParametersDiff.Modified {
				for param, paramDiff := range params {
					errs = append(errs, p.checkModifiedPathDiff(path, &diff.PathDiff{
						OperationsDiff: &diff.OperationsDiff{
							Modified: map[string]*diff.MethodDiff{op: {
								ParametersDiff: &diff.ParametersDiff{
									Modified: map[string]diff.ParamDiffs{
										paramLocation: map[string]*diff.ParameterDiff{param: paramDiff},
									},
								},
							}},
						},
					})...)
				}
			}
		}

		if opDiff.ResponsesDiff != nil {
			if len(opDiff.ResponsesDiff.Added) > 0 {
				errs = append(errs, p.checkModifiedPathDiff(path, &diff.PathDiff{
					OperationsDiff: &diff.OperationsDiff{
						Modified: map[string]*diff.MethodDiff{op: {ResponsesDiff: &diff.ResponsesDiff{Added: opDiff.ResponsesDiff.Added}}},
					},
				})...)
			}

			if len(opDiff.ResponsesDiff.Deleted) > 0 {
				errs = append(errs, p.checkModifiedPathDiff(path, &diff.PathDiff{
					OperationsDiff: &diff.OperationsDiff{
						Modified: map[string]*diff.MethodDiff{op: {ResponsesDiff: &diff.ResponsesDiff{Deleted: opDiff.ResponsesDiff.Deleted}}},
					},
				})...)
			}

			for statusCode, respDiff := range opDiff.ResponsesDiff.Modified {
				if respDiff.ContentDiff == nil {
					continue
				}

				if len(respDiff.ContentDiff.MediaTypeAdded) > 0 {
					errs = append(errs, p.checkModifiedPathDiff(path, &diff.PathDiff{
						OperationsDiff: &diff.OperationsDiff{
							Modified: map[string]*diff.MethodDiff{op: {
								Base:     opDiff.Base,
								Revision: opDiff.Revision,
								ResponsesDiff: &diff.ResponsesDiff{
									Modified: map[string]*diff.ResponseDiff{
										statusCode: {
											ContentDiff: &diff.ContentDiff{MediaTypeAdded: respDiff.ContentDiff.MediaTypeAdded},
										},
									},
								}}},
						},
					})...)
				}

				if len(respDiff.ContentDiff.MediaTypeDeleted) > 0 {
					errs = append(errs, p.checkModifiedPathDiff(path, &diff.PathDiff{
						OperationsDiff: &diff.OperationsDiff{
							Modified: map[string]*diff.MethodDiff{op: {
								Base:     opDiff.Base,
								Revision: opDiff.Revision,
								ResponsesDiff: &diff.ResponsesDiff{
									Modified: map[string]*diff.ResponseDiff{
										statusCode: {
											ContentDiff: &diff.ContentDiff{MediaTypeDeleted: respDiff.ContentDiff.MediaTypeDeleted},
										},
									},
								}}},
						},
					})...)
				}

				for mediaType, mediaTypeDiff := range respDiff.ContentDiff.MediaTypeModified {
					// todo: need to go deeper to mediaTypeDiff.SchemaDiff
					// todo: need to go deeper to mediaTypeDiff.SchemaDiff.PropertiesDiff.Modified
					// todo: need to go deeper to mediaTypeDiff.SchemaDiff.PropertiesDiff.Added/Deleted/Modified
					errs = append(errs, p.checkModifiedPathDiff(path, &diff.PathDiff{
						OperationsDiff: &diff.OperationsDiff{
							Modified: map[string]*diff.MethodDiff{op: {
								Base:     opDiff.Base,
								Revision: opDiff.Revision,
								ResponsesDiff: &diff.ResponsesDiff{
									Modified: map[string]*diff.ResponseDiff{
										statusCode: {
											ContentDiff: &diff.ContentDiff{
												MediaTypeModified: map[string]*diff.MediaTypeDiff{mediaType: mediaTypeDiff},
											},
										},
									},
								}}},
						},
					})...)
				}
			}
		}
		// TODO: CallbacksDiff ?

		errs = append(errs, p.checkModifiedPathDiff(path, &diff.PathDiff{
			OperationsDiff: &diff.OperationsDiff{
				Modified: map[string]*diff.MethodDiff{op: opDiff},
			},
		})...)
	}

	return errs
}

func (p *changelogProfile) checkModifiedPathDiff(path string, pathDiff *diff.PathDiff) checker.BackwardCompatibilityErrors {
	pathDiff.Base = p.diffReport.PathsDiff.Modified[path].Base
	pathDiff.Revision = p.diffReport.PathsDiff.Modified[path].Revision
	return p.check(&diff.Diff{
		PathsDiff: &diff.PathsDiff{
			Modified: map[string]*diff.PathDiff{path: pathDiff},
			Base:     p.diffReport.PathsDiff.Base,
			Revision: p.diffReport.PathsDiff.Revision,
		},
	})
}

func (p *changelogProfile) check(dr *diff.Diff) checker.BackwardCompatibilityErrors {
	pathErrs := checker.CheckBackwardCompatibility(p.checks, dr, p.operationsSources)
	if pathErrs.Len() == 0 {
		textOutput := report.GetTextReportAsString(dr)
		// s, _ := ToYAML(dr)
		// fmt.Printf("no oasdiff checker matched:\n%s\n%s", s, textOutput)
		fmt.Printf("USING text output:\n%s", textOutput)
		//os.Exit(missingCheckerErrCode)
	}
	return pathErrs
}
