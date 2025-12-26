package mapper

import (
	"omar-kada/autonas/api"
	"omar-kada/autonas/models"
)

// DiffMapper is a mapper that converts models.FileDiff to api.FileDiff.
type DiffMapper struct{}

// Map converts a models.FileDiff to an api.FileDiff.
func (DiffMapper) Map(file models.FileDiff) api.FileDiff {
	return api.FileDiff{
		Diff:    file.Diff,
		NewFile: file.NewFile,
		OldFile: file.OldFile,
	}
}
