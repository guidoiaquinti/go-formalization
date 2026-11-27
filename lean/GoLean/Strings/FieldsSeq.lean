import GoLean.Unicode.Utf8.DecodeRuneInString

namespace GoLean.Strings

/-- Pinned asciiSpace entries in strings/strings.go:358. -/
private def asciiSpace (c : Nat) : Bool :=
  (9 ≤ c && c ≤ 13) || c == 32

/-- Latin-1 cases in unicode/graphic.go and the White_Space table in unicode/tables.go. -/
private def unicodeSpace (c : Nat) : Bool :=
  asciiSpace c || c == 0x85 || c == 0xa0 || c == 0x1680 ||
  (0x2000 ≤ c && c ≤ 0x200a) || c == 0x2028 || c == 0x2029 ||
  c == 0x202f || c == 0x205f || c == 0x3000

/-- ASCII classification precedes UTF-8 decoding, as in strings/iter.go:89–95. -/
private def fieldStep (lead : Nat) (rest : List Nat) : Bool × Nat :=
  if lead < 128 then (asciiSpace lead, 1)
  else
    let (rune, width) := GoLean.Unicode.Utf8.decodeNonempty lead rest
    (unicodeSpace rune, width)

private def scanFields (input : List Nat) (fuel : Nat) (i : Nat)
    (start : Option Nat) (budget : Option Nat) : List (List Nat) :=
  match fuel, budget with
  | _, some 0 => []
  | 0, _ => (start.map fun first => input.drop first).toList
  | fuel + 1, _ =>
    match input.drop i with
    | [] => (start.map fun first => input.drop first).toList
    | lead :: rest =>
      let (isSpace, width) := fieldStep lead rest
      let next := i + width
      if isSpace then
        match start with
        | none => scanFields input fuel next none budget
        | some first =>
          ((input.drop first).take (i - first)) ::
            scanFields input fuel next none (budget.map Nat.pred)
      else
        scanFields input fuel next (some (start.getD i)) budget

/--
Go 1.27.0 strings/iter.go:83–109. The scan uses byte offsets and an optional
open field. Each invocation starts from byte zero. A finite budget counts
the callback that returns false. A zero budget represents an omitted call.
-/
def fieldsSeq (input : List (BitVec 8)) (budget : Option Nat) : List (List Nat) :=
  scanFields (input.map BitVec.toNat) input.length 0 none budget

private theorem scanFields_take (input : List Nat) (fuel i : Nat) (start : Option Nat)
    (limit : Nat) :
    scanFields input fuel i start (some limit) = (scanFields input fuel i start none).take limit := by
  induction fuel generalizing i start limit with
  | zero =>
    cases limit with
    | zero => simp [scanFields]
    | succ limit => cases start <;> simp [scanFields]
  | succ fuel ih =>
    cases limit with
    | zero => simp [scanFields]
    | succ limit =>
      simp only [scanFields]
      split
      · cases start <;> simp
      · split
        · split
          · exact ih _ _ _
          · simp [ih]
        · exact ih _ _ _

/-- Stopping after a bounded number of callbacks yields the corresponding prefix. -/
theorem fieldsSeq_take (input : List (BitVec 8)) (limit : Nat) :
    fieldsSeq input (some limit) = (fieldsSeq input none).take limit := by
  exact scanFields_take _ _ _ _ _

end GoLean.Strings
