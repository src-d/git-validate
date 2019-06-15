package dco

import (
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/vbatts/git-validation/validate"
)

func init() {
	validate.RegisterRule(DCORule)
}

var (
	// ValidDCO is the regexp for signed off DCO
	ValidDCO = regexp.MustCompile(`^Signed-off-by: ([^<]+) <([^<>@]+@[^<>]+)>$`)
	// DcoRule is the rule being registered
	DCORule = validate.Rule{
		Name:        "DCO",
		Description: "makes sure the commits are signed",
		Run:         ValidateDCO,
		Default:     true,
	}
)

// ValidateDCO checks that the commit has been signed off, per the DCO process
func ValidateDCO(_ *git.Repository, c *object.Commit) (vr validate.Result, err error) {
	vr.Commit = c
	if c.NumParents() > 1 {
		vr.Pass = true
		vr.Msg = "merge commits do not require DCO"
		return vr, nil
	}

	hasValid := false
	for _, line := range strings.Split(c.Message, "\n") {
		if ValidDCO.MatchString(line) {
			hasValid = true
		}
	}
	if !hasValid {
		vr.Pass = false
		vr.Msg = "does not have a valid DCO"
	} else {
		vr.Pass = true
		vr.Msg = "has a valid DCO"
	}

	return vr, nil
}
