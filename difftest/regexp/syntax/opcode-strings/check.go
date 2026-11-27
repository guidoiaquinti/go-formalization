package main

import (
	"bytes"
	"fmt"
	"log"
	"regexp/syntax"
	"strings"

	"github.com/guidoiaquinti/go-formalization/difftest/internal/check"
)

func main() {
	expected := map[syntax.InstOp]string{
		syntax.InstAlt:          "InstAlt",
		syntax.InstAltMatch:     "InstAltMatch",
		syntax.InstCapture:      "InstCapture",
		syntax.InstEmptyWidth:   "InstEmptyWidth",
		syntax.InstMatch:        "InstMatch",
		syntax.InstFail:         "InstFail",
		syntax.InstNop:          "InstNop",
		syntax.InstRune:         "InstRune",
		syntax.InstRune1:        "InstRune1",
		syntax.InstRuneAny:      "InstRuneAny",
		syntax.InstRuneAnyNotNL: "InstRuneAnyNotNL",
	}
	opNames := map[syntax.Op]string{
		syntax.OpNoMatch:        "NoMatch",
		syntax.OpEmptyMatch:     "EmptyMatch",
		syntax.OpLiteral:        "Literal",
		syntax.OpCharClass:      "CharClass",
		syntax.OpAnyCharNotNL:   "AnyCharNotNL",
		syntax.OpAnyChar:        "AnyChar",
		syntax.OpBeginLine:      "BeginLine",
		syntax.OpEndLine:        "EndLine",
		syntax.OpBeginText:      "BeginText",
		syntax.OpEndText:        "EndText",
		syntax.OpWordBoundary:   "WordBoundary",
		syntax.OpNoWordBoundary: "NoWordBoundary",
		syntax.OpCapture:        "Capture",
		syntax.OpStar:           "Star",
		syntax.OpPlus:           "Plus",
		syntax.OpQuest:          "Quest",
		syntax.OpRepeat:         "Repeat",
		syntax.OpConcat:         "Concat",
		syntax.OpAlternate:      "Alternate",
		128:                     "opPseudo",
	}
	var input bytes.Buffer
	for raw := range 256 {
		fmt.Fprintln(&input, raw)
	}
	for raw, row := range check.LeanRows("godiff-regexp-syntax-opcode-strings", 256, &input) {
		inst, op := syntax.InstOp(raw), syntax.Op(raw)
		wantInst := expected[inst]
		wantOp, known := opNames[op]
		if !known {
			wantOp = fmt.Sprintf("Op(%d)", raw)
		}
		instRow, opRow, separated := strings.Cut(row, "|")
		if !separated || inst.String() != wantInst || instRow != wantInst ||
			op.String() != wantOp || opRow != wantOp {
			log.Fatalf("opcode %d: Go (%q, %q), Lean %q, want (%q, %q)",
				raw, inst.String(), op.String(), row, wantInst, wantOp)
		}
	}
	fmt.Println("all 256 opcode bytes agree with pinned Go and the independent name oracle")
}
