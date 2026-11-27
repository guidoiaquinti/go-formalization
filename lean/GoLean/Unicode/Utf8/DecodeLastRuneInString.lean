import GoLean.Unicode.Utf8.DecodeRuneInString
import GoLean.Unicode.Utf8.RuneStart

namespace GoLean.Unicode.Utf8

/-- Inspect exactly the source loop's candidate indices in descending order.
    Upon exhaustion, Go's loop index is one below the limit, clamped to zero. -/
def lastRuneStart (bytes : List Nat) (limit : Nat) : Nat → Nat
  | 0 => limit - 1
  | fuel + 1 =>
    let i := limit + fuel
    if runeStart (BitVec.ofNat 8 bytes[i]!) then i else lastRuneStart bytes limit fuel

private def decodeSuffix (bytes : List Nat) (start : Nat) : Nat × Nat :=
  let decoded := decodeRuneInString (bytes.drop start)
  if start + decoded.2 = bytes.length then decoded else (0xfffd, 1)

/-- Go 1.27.0 src/unicode/utf8/utf8.go:316–343. -/
def decodeLastRuneInString (bytes : List Nat) : Nat × Nat :=
  if bytes.isEmpty then (0xfffd, 0)
  else
    let lead := bytes[bytes.length - 1]!
    if lead < 128 then (lead, 1)
    else
      let limit := bytes.length - 4
      let start := lastRuneStart bytes limit (bytes.length - 1 - limit)
      decodeSuffix bytes start

theorem decodeLastRuneInString_empty : decodeLastRuneInString [] = (0xfffd, 0) := rfl

/-- A non-ASCII tail is decoded only at the selected candidate, and a
    decoded prefix that ends short of the final byte is a one-byte error. -/
theorem decodeLastRuneInString_spec (bytes : List Nat) (nonempty : bytes ≠ []) :
    decodeLastRuneInString bytes =
      let last := bytes[bytes.length - 1]!
      if last < 128 then (last, 1)
      else
        let limit := bytes.length - 4
        let start := lastRuneStart bytes limit (min (bytes.length - 1) 3)
        let decoded := decodeRuneInString (bytes.drop start)
        if decoded.2 = bytes.length - start then decoded else (0xfffd, 1) := by
  have budget : bytes.length - 1 - (bytes.length - 4) = min (bytes.length - 1) 3 := by omega
  have bound (l n : Nat) (h : l + n ≤ bytes.length) :
      lastRuneStart bytes l n ≤ bytes.length := by
    induction n with
    | zero => simp [lastRuneStart]; omega
    | succ n ih =>
      simp only [lastRuneStart]
      split
      · omega
      · exact ih (by omega)
  cases bytes with
  | nil => exact False.elim (nonempty rfl)
  | cons b bs =>
    simp only [decodeLastRuneInString, List.isEmpty_cons, Bool.false_eq_true, ↓reduceIte, budget,
      decodeSuffix]
    split
    · rfl
    · have h := bound ((b :: bs).length - 4) (min ((b :: bs).length - 1) 3) (by omega)
      let start := lastRuneStart (b :: bs) ((b :: bs).length - 4) (min ((b :: bs).length - 1) 3)
      let decoded := decodeRuneInString ((b :: bs).drop start)
      change (if start + decoded.2 = (b :: bs).length then decoded else (0xfffd, 1)) =
        (if decoded.2 = (b :: bs).length - start then decoded else (0xfffd, 1))
      have le : start ≤ (b :: bs).length := h
      have guard : (start + decoded.2 = (b :: bs).length) =
          (decoded.2 = (b :: bs).length - start) := propext (by omega)
      simp only [guard]

end GoLean.Unicode.Utf8
