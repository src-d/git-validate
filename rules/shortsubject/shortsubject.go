package shortsubject

import (
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/vbatts/git-validation/validate"
)

var (
	// ShortSubjectRule is the rule being registered
	ShortSubjectRule = validate.Rule{
		Name:        "short-message",
		Description: "commit message are strictly less than 90 (github ellipsis length)",
		Run:         ValidateShortSubject,
		Default:     true,
	}
)

func init() {
	validate.RegisterRule(ShortSubjectRule)
}

// ValidateShortSubject checks that the commit's subject is strictly less than
// 90 characters (preferably not more than 72 chars).
func ValidateShortSubject(_ *git.Repository, c *object.Commit) (vr validate.Result, err error) {
	if c.NumParents() > 1 {
		vr.Pass = true
		vr.Msg = "merge commits do not require length check"
		return vr, nil
	}

	lines := strings.SplitN(c.Message, "\n", 2)
	if len(lines[0]) >= 90 {
		vr.Pass = false
		vr.Msg = "commit subject exceeds 90 characters"
		return
	}

	vr.Pass = true
	if len(lines[0]) > 72 {
		vr.Msg = "commit subject is under 90 characters, but is still more than 72 chars"
	} else {
		vr.Msg = "commit subject is 72 characters or less! *yay*"
	}

	return
}
