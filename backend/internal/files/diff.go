package files

import (
	"fmt"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// DiffText returns a string showing the differences between oldStr and newStr
// in a simple diff format with '-' for deletions and '+' for insertions.
func DiffText(oldStr, newStr string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffCleanupSemantic(dmp.DiffMain(oldStr, newStr, true))
	res := ""
	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			res += fmt.Sprintf("- %s\n", diff.Text)
		case diffmatchpatch.DiffInsert:
			res += fmt.Sprintf("+ %s\n", diff.Text)
		}
	}
	return res
}
