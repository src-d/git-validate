package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/src-d/git-compliance/compliance"

	_ "github.com/src-d/git-compliance/rule/dco"
	_ "github.com/src-d/git-compliance/rule/dockerfile"
	_ "github.com/src-d/git-compliance/rule/file"
	_ "github.com/src-d/git-compliance/rule/shortsubject"
	_ "github.com/src-d/git-compliance/rule/stalebranch"
	"gopkg.in/src-d/go-git.v4"
)

var (
	flDir = flag.String("d", ".", "git directory to compliance from")
)

func main() {
	flag.Parse()

	var cfg compliance.Config
	runner, err := compliance.NewRunner(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	r, err := git.PlainOpen(*flDir)
	if err != nil {
		log.Fatal(err)
	}

	results, err := runner.Run(r)
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		fmt.Println(result)
	}

	if len(results) > 0 {
		fmt.Printf("%d commits to fix\n", len(results))
		os.Exit(1)
	}

}
