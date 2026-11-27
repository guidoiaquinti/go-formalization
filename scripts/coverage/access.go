package main

import (
	"go/ast"
	"strings"
)

func externalCallRestriction(row declaration) string {
	// An exported method on a private type may still be callable on a value
	// returned by a public function. Check the method name itself.
	name := row.Symbol[strings.LastIndexByte(row.Symbol, '.')+1:]
	var reasons []string
	if !ast.IsExported(name) {
		reasons = append(reasons, "unexported declaration")
	}
	if strings.Contains("/"+row.Package+"/", "/internal/") {
		reasons = append(reasons, "Go internal import restriction")
	}
	return strings.Join(reasons, ", ")
}
