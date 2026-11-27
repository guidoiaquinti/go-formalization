import Init.Data.BitVec.Basic

namespace GoLean.Regexp.Syntax

private def instOpNames : Array String := #[
  "InstAlt", "InstAltMatch", "InstCapture", "InstEmptyWidth", "InstMatch", "InstFail",
  "InstNop", "InstRune", "InstRune1", "InstRuneAny", "InstRuneAnyNotNL"]

/-- Pinned Go 1.27.0 src/regexp/syntax/prog.go:27–60. -/
def instOpString (op : BitVec 8) : String :=
  if h : op.toNat < instOpNames.size then instOpNames[op.toNat] else ""

private def opNames : Array String := #[
  "NoMatch", "EmptyMatch", "Literal", "CharClass", "AnyCharNotNL", "AnyChar",
  "BeginLine", "EndLine", "BeginText", "EndText", "WordBoundary", "NoWordBoundary",
  "Capture", "Star", "Plus", "Quest", "Repeat", "Concat", "Alternate"]

/-- Pinned Go 1.27.0 src/regexp/syntax/op_string.go:42–52. -/
def opString (op : BitVec 8) : String :=
  let i := op.toNat
  if 1 ≤ i && i ≤ 19 then opNames[i - 1]?.getD ""
  else if i = 128 then "opPseudo"
  else s!"Op({i})"

theorem instOpString_outside (op : BitVec 8) (outside : 11 ≤ op.toNat) :
    instOpString op = "" := by
  have noIndex : ¬ op.toNat < instOpNames.size := Nat.not_lt.mpr outside
  simp only [instOpString, dite_eq_right noIndex]

theorem instOpString_last : instOpString 10 = "InstRuneAnyNotNL" := rfl

theorem opString_pseudo : opString 128 = "opPseudo" := rfl

theorem opString_last : opString 19 = "Alternate" := rfl

theorem opString_outside (op : BitVec 8) (outside : ¬ (1 ≤ op.toNat ∧ op.toNat ≤ 19))
    (notPseudo : op.toNat ≠ 128) : opString op = s!"Op({op.toNat})" := by
  simp [opString, outside, notPseudo, Bool.and_eq_true]

end GoLean.Regexp.Syntax
