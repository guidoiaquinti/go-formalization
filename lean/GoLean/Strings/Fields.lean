import GoLean.Strings.FieldsSeq

namespace GoLean.Strings

private def asciiByteSpace (b : BitVec 8) : Bool :=
  (9 ≤ b.toNat && b.toNat ≤ 13) || b.toNat == 32

/-- The source counts transitions from space and ORs byte bits before selecting a path. -/
private def fieldPrescan : List (BitVec 8) → Bool → Nat × Bool
  | [], _ => (0, false)
  | byte :: rest, wasSpace =>
    let space := asciiByteSpace byte
    let (count, nonASCII) := fieldPrescan rest space
    (count + (if wasSpace && !space then 1 else 0), nonASCII || byte.msb)

/-- Skip a separator run, then consume the next nonempty field if one remains. -/
private def asciiFields (input : List (BitVec 8)) : Nat → List (List Nat)
  | 0 => []
  | fuel + 1 =>
    let remaining := input.dropWhile asciiByteSpace
    let field := remaining.takeWhile (!asciiByteSpace ·)
    if field.isEmpty then []
    else field.map BitVec.toNat :: asciiFields (remaining.drop (field.length + 1)) fuel

/-- Go 1.27.0 strings/strings.go:365–412. The prescan supplies the ASCII
    result count. The Unicode branch represents the exact output of
    FieldsFunc with unicode.IsSpace via the shared complete field scan. -/
def fields (input : List (BitVec 8)) : List (List Nat) :=
  let (count, nonASCII) := fieldPrescan input true
  if nonASCII then fieldsSeq input none else asciiFields input count

/-- If the ASCII prescan finds a field, its first result ends at the next space. -/
theorem fields_ascii_first (input : List (BitVec 8)) (count : Nat)
    (ascii : fieldPrescan input true = (count + 1, false))
    (nonempty : (input.dropWhile asciiByteSpace |>.takeWhile (!asciiByteSpace ·)) ≠ []) :
    (fields input).head? =
      some (((input.dropWhile asciiByteSpace).takeWhile (!asciiByteSpace ·)).map BitVec.toNat) := by
  simp [fields, ascii, asciiFields, nonempty]

end GoLean.Strings
