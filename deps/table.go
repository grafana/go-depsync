package deps

import (
	"io"

	"github.com/olekukonko/tablewriter"
)

// WriteVersionTable writes to w an ascii table showing the dependency versions that have changed between old and new.
// Dependencies present in old but not present in new are not printed out, and all dependencies in new are assumed to
// be also in old.
// This function is meant to be called with the output of Match as the value of new.
func WriteVersionTable(w io.Writer, old, new Dependencies) {
	tw := tablewriter.NewWriter(w)

	tw.SetHeader([]string{"Dependency", "Current version", "New version"})
	for dep, v := range new {
		tw.Append([]string{dep, old[dep], v})
	}

	tw.Render()
}
