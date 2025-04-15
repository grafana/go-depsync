package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/grafana/go-depsync/deps"
)

const long = `
synchronizes dependencies between two go modules

shows the diff between the current module and an target module and the go command(s)
required to synchronize them.

If the version of the target module is not explicitly defined, it is obtained from the
the module's dependencies.

This tools is useful when a module is used as a dependency in another module (the target),
to ensure it won't introduce unwanted dependency upgrades.
`

// print usage message
func printHelp() {
	fmt.Printf("Usage of %s:\n", os.Args[0])
	fmt.Print(long)
	flag.PrintDefaults()
}

func main() {
	var (
		gomod    string
		target string
		version  string
		usage    bool
	)

	flag.StringVar(&gomod, "gomod", "./go.mod", "Path to the local go.mod file")
	flag.StringVar(&target, "parent", "", "Deprecated. Use '-target' instead")
	flag.StringVar(&target, "target", "", "Name of the package to sync dependencies with")
	flag.StringVar(&version, "version", "", "Version of target package. If not defined it is obtained from the local go.mod")
	flag.BoolVar(&usage, "usage", false, "display long help")
	flag.Parse()

	if usage {
		printHelp()
		return
	}

	if target == "" {
		flag.Usage()
		log.Fatalf("You must specify the name of the target package")
	}

	ownDeps, err := deps.FromGomodFile(gomod)
	if err != nil {
		log.Fatalf("Couldn't parse own dependencies: %v", err)
	}

	if version == "" {
		hasParent := false
		version, hasParent = ownDeps[target]
		if !hasParent {
			flag.Usage()
			log.Fatalf("target version not specified and target module not found as a dependency")
		}

		log.Printf("Found target %s@%s", target, version)
	}

	log.Printf("synchronizing with target %s@%s", target, version)

	parentDeps, err := deps.FromModule(target, version)
	if err != nil {
		log.Fatalf("Cannot parse target dependencies %v", err)
	}

	mismatched := deps.Mismatched(ownDeps, parentDeps)
	if len(mismatched) == 0 {
		return
	}

	deps.WriteVersionTable(os.Stderr, ownDeps, mismatched)
	fmt.Fprintln(os.Stderr)

	fmt.Println(deps.GoGetCommand(mismatched))
}
