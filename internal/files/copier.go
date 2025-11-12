package files

import (
	"github.com/otiai10/copy"
)

// Copier abstracts copying folder
type Copier interface {
	Copy(src, dest string) error
}

type folderCopier struct{}

// NewCopier creates and new Copier and returns it
func NewCopier() Copier {
	return folderCopier{}
}

func (folderCopier) Copy(src, dst string) error {
	return copy.Copy(src, dst)
}
