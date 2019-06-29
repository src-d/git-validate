package compliance

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// Location describes the location where the rule was triggered.
type Location interface {
	IsLocation()
	String() string
}

// ReferenceLocation the rules was triggered at a git reference.
type ReferenceLocation struct {
	Reference *plumbing.Reference
}

// IsLocation honors the Location interface.
func (ReferenceLocation) IsLocation() {}
func (loc *ReferenceLocation) String() string {
	return loc.Reference.Name().String()
}

// CommitLocation the rules was triggered at a git commit.
type CommitLocation struct {
	Commit *object.Commit
}

// IsLocation honors the Location interface.
func (CommitLocation) IsLocation() {}
func (loc *CommitLocation) String() string {
	return loc.Commit.Hash.String()[:6]
}

// FileLocation the rules was triggered at a file in a commit.
type FileLocation struct {
	Commit   *object.Commit
	Filename string
}

// IsLocation honors the Location interface.
func (FileLocation) IsLocation() {}
func (loc *FileLocation) String() string {
	return fmt.Sprintf("%s@%s", loc.Filename, loc.Commit.Hash.String()[:6])
}

// LineLocation the rules was triggered at a line of file in a commit.
type LineLocation struct {
	Commit   *object.Commit
	Filename string
	Line     int
}

// IsLocation honors the Location interface.
func (LineLocation) IsLocation() {}
func (loc *LineLocation) String() string {
	return fmt.Sprintf("%s:%d@%s", loc.Filename, loc.Line, loc.Commit.Hash.String()[:6])
}
