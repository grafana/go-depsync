package deps

import (
	"fmt"
	"log"
)

// Mismatched looks at the dependencies of a package and its parent, and returns the dependencies that are present in
// both packages with the version they have on the parent package.
func Mismatched(own, parent Dependencies) Dependencies {
	mismatched := Dependencies{}

	for dep, version := range own {
		parentVersion, inParent := parent[dep]
		if !inParent {
			continue
		}

		if version == parentVersion {
			continue
		}

		log.Printf("Mismatched versions for %s: %s (this package) -> %s (parent)", dep, version, parentVersion)
		mismatched[dep] = parentVersion
	}

	return mismatched
}

// GoGetCommand returns a go get command that installs the provided dependencies and versions.
func GoGetCommand(deps Dependencies) string {
	if len(deps) == 0 {
		return ""
	}

	cmd := "go get"

	for dep, version := range deps {
		cmd += fmt.Sprintf(" %s@%s", dep, version)
	}

	return cmd
}
