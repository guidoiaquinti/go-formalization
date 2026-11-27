package dag

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

func TestGoLeanParseRules(t *testing.T) {
	data, err := os.ReadFile(os.Getenv("GOLEAN_DAG_PARSE_RULES_INPUT"))
	if err != nil {
		t.Fatal(err)
	}
	var cases []struct {
		Text     []byte
		Expected string
	}
	if err := json.Unmarshal(data, &cases); err != nil {
		t.Fatal(err)
	}
	for i, c := range cases {
		rules, err := parseRules(string(c.Text))
		if len(rules) == 0 && rules != nil {
			t.Fatalf("case %d allocated an empty rules slice", i)
		}
		kind := "ok"
		if err != nil {
			if _, ok := err.(syntaxError); !ok {
				t.Fatalf("case %d returned unexpected error type %T", i, err)
			}
			kind = "error"
		}
		row := []string{kind, strconv.Itoa(len(rules))}
		for _, r := range rules {
			row = append(row, goLeanEncodeRuleNames(r.less), goLeanEncodeRuleBytes(r.op), goLeanEncodeRuleNames(r.def))
		}
		if err != nil {
			row = append(row, goLeanEncodeRuleBytes(err.Error()))
		}
		if got := strings.Join(row, " "); got != c.Expected {
			t.Fatalf("case %d text %q: Go=%q expected=%q", i, c.Text, got, c.Expected)
		}
	}
}

func goLeanEncodeRuleNames(names []string) string {
	row := []string{strconv.Itoa(len(names))}
	for _, name := range names {
		row = append(row, goLeanEncodeRuleBytes(name))
	}
	return strings.Join(row, " ")
}

func goLeanEncodeRuleBytes(s string) string {
	var out strings.Builder
	fmt.Fprint(&out, len(s))
	for i := range len(s) {
		fmt.Fprint(&out, " ", s[i])
	}
	return out.String()
}
