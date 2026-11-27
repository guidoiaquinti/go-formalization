import GoLean.Unicode.Utf16.AppendRune
import GoLean.Unicode.Utf16.RuneLen

namespace GoLean.Unicode.Utf16

private def writeRune (r : Int) : List Nat :=
  match runeLen r with
  | 1 => [r.toNat]
  | 2 =>
    let (hi, lo) := encodeRune r
    [hi, lo]
  | _ => [0xfffd]

/-- Go counts potential two-word runes, allocates that many words, writes the
    encoding, and slices the result to the written length. The model tracks
    returned words and capacity. The writes are represented by concatenation.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf16/utf16.go#L68-L95 -/
def encode (runes : List Int) : List Nat × Nat :=
  let capacity := runes.foldl (fun n r => if 65536 ≤ r then n + 1 else n) runes.length
  (runes.flatMap writeRune, capacity)

private theorem writeRune_eq_appendRune (r : Int) :
    writeRune r = appendRune [] r := by
  by_cases bmp : (0 ≤ r ∧ r < 0xd800) ∨ (0xe000 ≤ r ∧ r < 0x10000)
  · simp [writeRune, runeLen, appendRune, bmp]
  · by_cases supplementary : 0x10000 ≤ r ∧ r ≤ 0x10ffff
    · simp [writeRune, runeLen, appendRune, bmp, supplementary]
    · simp [writeRune, runeLen, appendRune, bmp, supplementary]

private theorem capacity_eq_count (runes : List Int) (n : Nat) :
    runes.foldl (fun acc r => if 65536 ≤ r then acc + 1 else acc) n =
      n + runes.countP (fun r => decide (0 < r / 65536)) := by
  induction runes generalizing n with
  | nil => rfl
  | cons r tail ih =>
    have plane : (65536 ≤ r) ↔ 0 < r / 65536 := by omega
    rw [List.foldl_cons, ih, List.countP_cons]
    by_cases h : 65536 ≤ r
    · simp [h, plane.mp h]
      omega
    · simp [h, mt plane.mpr h]

theorem encode_spec (runes : List Int) :
    encode runes =
      (runes.flatMap (appendRune []),
        runes.length + runes.countP (fun r => decide (0 < r / 65536))) := by
  rw [encode, capacity_eq_count, show writeRune = appendRune [] from funext writeRune_eq_appendRune]

end GoLean.Unicode.Utf16
