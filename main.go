package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"

	"github.com/mixpanel/terraform-deps/parser"
)

func main() {
	p := parser.Parser{
		BaseDir: ".",
		TFFiles: regexp.MustCompile(`\.tf$`),
	}
	deps, err := p.Deps()
	if err != nil {
		log.Fatal(err)
	}
	//dump as graph
	fmt.Println("graph TD")
	lines := []string{}
	for modDir, deps := range deps {
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
