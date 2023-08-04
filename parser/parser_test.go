package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"gotest.tools/golden"
)

func dumpDeps(deps map[string]map[string]struct{}) string {
	lines := []string{}
	for modDir, deps := range deps {
		for dep := range deps {
			lines = append(lines, fmt.Sprintf("%s --> %s", modDir, dep))
		}
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

func TestCases(t *testing.T) {
	dir, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range dir {
		if !d.IsDir() {
			continue
		}
		name := d.Name()
		t.Run(name, func(t *testing.T) {
			p := Parser{
				BaseDir: filepath.Join("testdata", name),
				TFFiles: regexp.MustCompile(`\.tf$`),
			}
			deps, err := p.Deps()
			if err != nil {
				t.Fatal(err)
			}
			dump := dumpDeps(deps)
			golden.Assert(t, dump, name+".golden")
		})
	}
}
