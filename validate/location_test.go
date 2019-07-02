package validate

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func TestReferenceLocation(t *testing.T) {
	var loc Location
	loc = &ReferenceLocation{
		Reference: plumbing.NewReferenceFromStrings("foo", plumbing.ZeroHash.String()),
	}

	assert.Equal(t, "foo", loc.String())
}

func TestCommitLocation(t *testing.T) {
	var loc Location
	loc = &CommitLocation{
		Commit: &object.Commit{},
	}

	assert.Equal(t, "000000", loc.String())
}

func TestBlobLocation(t *testing.T) {
	var loc Location
	loc = &BlobLocation{
		Blob: &object.Blob{},
	}

	assert.Equal(t, "000000", loc.String())
}

func TestFileLocation(t *testing.T) {
	var loc Location
	loc = &FileLocation{
		Filename: "foo",
		Commit:   &object.Commit{},
	}

	assert.Equal(t, "foo@000000", loc.String())
}

func TestLineLocation(t *testing.T) {
	var loc Location
	loc = &LineLocation{
		Commit:   &object.Commit{},
		Filename: "foo",
		Line:     42,
	}

	assert.Equal(t, "foo:42@000000", loc.String())
}
