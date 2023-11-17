package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/grafana/go-depsync/deps"
)

func main() {
	gomod := flag.String("gomod", "./go.mod", "Path to the local go.mod file")
	parent := flag.String("parent", "", "Name of the parent package to sync dependencies with")
	flag.Parse()

	if *parent == "" {
		flag.Usage()
		log.Fatalf("You must specify the name of the parent package")
	}

	ownDeps, err := deps.FromGomodFile(*gomod)
	if err != nil {
		log.Fatalf("Couldn't parse own dependencies: %v", err)
	}

	parentVer, hasParent := ownDeps[*parent]
	if !hasParent {
		log.Fatalf("Parent package %q not found in local go.mod %q", *parent, *gomod)
	}

	log.Printf("Found parent %s@%s", *parent, parentVer)

	parentDeps, err := deps.FromModule(*parent, parentVer)
	if err != nil {
		log.Fatalf("Cannot parse dependencies of %q: %v", *parent, err)
	}

	mismatched := deps.Mismatched(ownDeps, parentDeps)
	if len(mismatched) == 0 {
		return
	}

	deps.WriteVersionTable(os.Stderr, ownDeps, mismatched)
	fmt.Fprintln(os.Stderr)

	fmt.Println(deps.GoGetCommand(mismatched))
}
