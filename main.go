package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

func main() {
	parser := hclparse.NewParser()
	depTree := map[string]map[string]struct{}{}
	err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) != ".tf" {
			return nil
		}
		modDir := filepath.Dir(path)
		f, diags := parser.ParseHCLFile(path)
		if diags.HasErrors() {
			return diags
		}
		schema := &hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{
				{Type: "module", LabelNames: []string{"name"}},
			},
		}
		moduleSchema := &hcl.BodySchema{
			Attributes: []hcl.AttributeSchema{
				{Name: "source", Required: true},
			},
		}
		content, _, diags := f.Body.PartialContent(schema)
		if diags.HasErrors() {
			fmt.Println(diags)
		}
		for _, block := range content.Blocks {

			content, _, _ := block.Body.PartialContent(moduleSchema)
			source := content.Attributes["source"]
			if source == nil {
				continue
			}
			if depTree[modDir] == nil {
				depTree[modDir] = map[string]struct{}{}
			}
			val, err := source.Expr.Value(nil)
			if err.HasErrors() {
				return err
			}
			src := val.AsString()
			if !strings.HasPrefix(src, "./") && !strings.HasPrefix(src, "../") {
				continue
			}
			src = filepath.Clean(filepath.Join(modDir, src))
			depTree[modDir][src] = struct{}{}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	//dump as graph
	fmt.Println("graph TD")
	lines := []string{}
	for modDir, deps := range depTree {
		for dep := range deps {
			lines = append(lines, fmt.Sprintf("\t%s --> %s", modDir, dep))
		}
	}
	sort.Strings(lines)
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println("end")
}
