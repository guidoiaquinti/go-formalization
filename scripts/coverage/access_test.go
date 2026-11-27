package main

import "testing"

func TestExternalCallRestriction(t *testing.T) {
	for _, tc := range []struct {
		pkg, symbol, want string
	}{
		{"strings", "Clone", ""},
		{"strings", "private.String", ""},
		{"strings", "Builder.copyCheck", "unexported declaration"},
		{"unicode/utf16", "decode", "unexported declaration"},
		{"internal/dag", "Graph.Edges", "Go internal import restriction"},
		{"encoding/json/internal/jsonwire", "appendHex", "unexported declaration, Go internal import restriction"},
		{"example/internalized", "Read", ""},
	} {
		if got := externalCallRestriction(declaration{Package: tc.pkg, Symbol: tc.symbol}); got != tc.want {
			t.Errorf("%s.%s got %q, want %q", tc.pkg, tc.symbol, got, tc.want)
		}
	}
}
