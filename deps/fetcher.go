package deps

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/mod/modfile"
)

type Dependencies map[string]string

// FromGomod returns a map of dependencies given the contents of a go.mod file.
func FromGomod(gomod []byte) (Dependencies, error) {
	parsed, err := modfile.Parse("go.mod", gomod, nil)
	if err != nil {
		return nil, fmt.Errorf("parsing contents: %w", err)
	}

	deps := make(map[string]string)
	for _, req := range parsed.Require {
		deps[req.Mod.Path] = req.Mod.Version
	}

	return deps, nil
}

// FromModule returns a map of dependencies for a given go module, given its name and version.
// It retrieves the go.mod file from proxy.golang.org.
func FromModule(pkg, version string) (Dependencies, error) {
	url := fmt.Sprintf("https://proxy.golang.org/%s/@v/%s.mod", pkg, version)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("requesting %q: %w", url, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%q returned unexpected status %d", url, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response from %q: %w", url, err)
	}

	return FromGomod(data)
}

// FromGomodFile returns a map of dependencies given the path to a local go.mod file.
func FromGomodFile(path string) (Dependencies, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening %q: %w", path, err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("reading %q: %w", path, err)
	}

	return FromGomod(data)
}
