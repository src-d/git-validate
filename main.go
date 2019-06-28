package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/src-d/git-compliance/rules/dco"
	_ "github.com/src-d/git-compliance/rules/dockerfile"
	_ "github.com/src-d/git-compliance/rules/file"
	_ "github.com/src-d/git-compliance/rules/shortsubject"
)

var (
	flRun          = flag.String("run", "", "comma delimited list of rules to run. Defaults to all.")
	flVerbose      = flag.Bool("v", false, "verbose")
	flDebug        = flag.Bool("D", false, "debug output")
	flQuiet        = flag.Bool("q", false, "less output")
	flDir          = flag.String("d", ".", "git directory to compliance from")
	flNoTravis     = flag.Bool("no-travis", false, "disables travis environment checks (when env TRAVIS=true is set)")
	flTravisPROnly = flag.Bool("travis-pr-only", true, "when on travis, only run validations if the CI-Build is checking pull-request build")
)

func main() {
	flag.Parse()

	if *flDebug {
		os.Setenv("DEBUG", "1")
	}
	if *flQuiet {
		os.Setenv("QUIET", "1")
	}

	if *flTravisPROnly && strings.ToLower(os.Getenv("TRAVIS_PULL_REQUEST")) == "false" {
		fmt.Printf("only to check travis PR builds and this not a PR build. yielding.\n")
		return
	}

	runner, err := NewRunner(*flDir, "compliance.yml", *flVerbose)
	if err != nil {
		log.Fatal(err)
	}

	results, err := runner.Run()
	if err != nil {
		log.Fatal(err)
	}
	_, fail := results.PassFail()
	if fail > 0 {
		fmt.Printf("%d commits to fix\n", fail)
		os.Exit(1)
	}

}
