package deps_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/go-depsync/deps"
)

func Test_Mismatched(t *testing.T) {
	own := deps.Dependencies{
		"something.tld/1": "v1.0.0",
		"something.tld/2": "v2.0.0",
		"something.tld/3": "v3.0.0",
		"something.tld/4": "v4.0.0",
		"something.tld/5": "v5.0.0",
		"something.tld/6": "v6.0.0",
	}

	parent := deps.Dependencies{
		"something.tld/2": "v2.0.0",
		"something.tld/3": "v9.9.9",
		"something.tld/4": "v9.9.9",
		"something.tld/8": "v8.0.0",
		"something.tld/9": "v9.0.0",
	}

	expected := deps.Dependencies{
		"something.tld/3": "v9.9.9",
		"something.tld/4": "v9.9.9",
	}

	mismatched := deps.Mismatched(own, parent)
	if diff := cmp.Diff(mismatched, expected); diff != "" {
		t.Fatalf("mismatched dependencies do not match expected:\n%s", diff)
	}
}

func Test_GoGetCommand(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name     string
		deps     deps.Dependencies
		expected string
	}{
		{
			name:     "No dependencies",
			deps:     deps.Dependencies{},
			expected: "",
		},
		{
			name: "One dependency",
			deps: deps.Dependencies{
				"foo/bar": "v1.2.3",
			},
			expected: "go get foo/bar@v1.2.3",
		},
		{
			name: "Several dependencies",
			deps: deps.Dependencies{
				"foo/bar": "v1.2.3",
				"foo/baz": "v4.5.6",
				"foo/boo": "v7.8.101",
			},
			expected: "go get foo/bar@v1.2.3 foo/baz@v4.5.6 foo/boo@v7.8.101",
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cmd := deps.GoGetCommand(tc.deps)
			if diff := cmp.Diff(cmd, tc.expected); diff != "" {
				t.Fatalf("command does not match expected:\n%s", diff)
			}
		})
	}
}
