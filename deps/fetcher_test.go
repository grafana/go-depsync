package deps_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/go-depsync/deps"
)

func Test_FromGomod(t *testing.T) {
	t.Parallel()

	gomod := `
module github.com/grafana/go-depsync

go 1.21.4

require golang.org/x/mod v0.14.0
`

	dependencies, err := deps.FromGomod([]byte(gomod))
	if err != nil {
		t.Fatalf("parsing dependencies from go.mod: %v", err)
	}

	expected := deps.Dependencies{
		"golang.org/x/mod": "v0.14.0",
	}

	if diff := cmp.Diff(dependencies, expected); diff != "" {
		t.Fatalf("dependencies do not match expected:\n%s", diff)
	}
}

func Test_FromPackage(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		target string
		version string
	}{
		{
			target: "go.k6.io/k6",
			version: "v0.47.0",
		},
		{
			target: "github.com/prometheus/prometheus",
			version: "v0.35.0",
		},
	} {
		tc := tc
		t.Run(tc.target, func(t *testing.T) {
			t.Parallel()

			dependencies, err := deps.FromModule(tc.target, tc.version)
			if err != nil {
				t.Fatalf("fetching dependencies: %v", err)
			}

			if len(dependencies) == 0 {
				t.Fatalf("expected to have at least one dependency, returned %d", len(dependencies))
			}
		})
	}
}
