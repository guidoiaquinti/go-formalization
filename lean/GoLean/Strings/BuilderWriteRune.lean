import GoLean.Strings.BuilderWriteString
import GoLean.Strings.IndexRune

namespace GoLean.Strings

/-- Go's append encodes invalid rune arguments as U+FFFD. -/
def builderRuneBytes (r : Int) : List (BitVec 8) :=
  (encodeRune (if validRune r then r.toNat else 0xfffd)).map (BitVec.ofNat 8)

/-- The source checks receiver identity before UTF-8 append and returns
    the byte length added to the buffer.
    https://github.com/golang/go/blob/go1.27.0/src/strings/builder.go#L102-L109
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L410-L437 -/
def builderWriteRune (receiver : Option BuilderState) (copied : Bool)
    (r : Int) : Option BuilderWriteResult :=
  builderWriteString receiver copied (builderRuneBytes r)

theorem builderWriteRune_success (b : BuilderState) (copied : Bool) (r : Int)
    (canWrite : (b.identitySet && copied) = false) :
    builderWriteRune (some b) copied r =
      some ⟨some (builderRuneBytes r).length, (b.buffer.getD ([], 0)).1 ++ builderRuneBytes r, true⟩ :=
  builderWriteString_success _ _ _ canWrite

theorem builderWriteRune_invalid (b : BuilderState) (copied : Bool) (r : Int)
    (canWrite : (b.identitySet && copied) = false) (invalid : validRune r = false) :
    builderWriteRune (some b) copied r =
      some ⟨some 3, (b.buffer.getD ([], 0)).1 ++ [0xef, 0xbf, 0xbd], true⟩ := by
  rw [builderWriteRune_success b copied r canWrite]
  have replacement : builderRuneBytes r = [0xef, 0xbf, 0xbd] := by
    unfold builderRuneBytes
    simp only [invalid, Bool.false_eq_true, ite_false]
    decide
  rw [replacement]
  rfl

theorem builderWriteRune_copied (b : BuilderState) (r : Int) (identified : b.identitySet = true) :
    builderWriteRune (some b) true r =
      some ⟨none, (b.buffer.getD ([], 0)).1, true⟩ :=
  builderWriteString_copied _ _ identified

theorem builderWriteRune_nil (copied : Bool) (r : Int) :
    builderWriteRune none copied r = none := rfl

end GoLean.Strings
