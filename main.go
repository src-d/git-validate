package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/src-d/git-validate/validate"

	_ "github.com/src-d/git-validate/rule/dco"
	_ "github.com/src-d/git-validate/rule/dockerfile"
	_ "github.com/src-d/git-validate/rule/file"
	_ "github.com/src-d/git-validate/rule/shortsubject"
	_ "github.com/src-d/git-validate/rule/stalebranch"
	"gopkg.in/src-d/go-git.v4"
)

var (
	flDir = flag.String("d", ".", "git directory to compliance from")
)

func main() {
	flag.Parse()

	var cfg validate.Config
	runner, err := validate.NewRunner(&cfg)
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
