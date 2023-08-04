package parser

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
)

type Parser struct {
	BaseDir string
	TFFiles *regexp.Regexp
}

func (p Parser) Deps() (map[string]map[string]struct{}, error) {
	parser := hclparse.NewParser()
	depTree := map[string]map[string]struct{}{}
	err := filepath.WalkDir(p.BaseDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !p.TFFiles.MatchString(path) {
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
		return nil, err
	}
	return depTree, nil
}
