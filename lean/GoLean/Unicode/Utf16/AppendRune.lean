import GoLean.Unicode.Utf16.EncodeRune

namespace GoLean.Unicode.Utf16

/-- Go appends a single BMP word, delegates supplementary pairs to EncodeRune,
    and appends U+FFFD for invalid runes. Each appended word fits in uint16.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf16/utf16.go#L97-L112 -/
def appendRune (dst : List Nat) (r : Int) : List Nat :=
  if (0 ≤ r ∧ r < 0xd800) ∨ (0xe000 ≤ r ∧ r < 0x10000) then
    dst ++ [r.toNat]
  else if 0x10000 ≤ r ∧ r ≤ 0x10ffff then
    let (first, second) := encodeRune r
    dst ++ [first, second]
  else
    dst ++ [0xfffd]

theorem appendRune_spec (dst : List Nat) (r : Int) :
    appendRune dst r =
      dst ++ (if r / 65536 = 0 ∧ r / 1024 ≠ 54 ∧ r / 1024 ≠ 55 then [r.toNat]
        else if 1 ≤ r / 65536 ∧ r / 65536 ≤ 16 then
          let payload := (r - 65536).toNat
          [55296 + payload / 1024, 56320 + payload % 1024]
        else [65533]) := by
  have bmp : ((0 ≤ r ∧ r < 0xd800) ∨ (0xe000 ≤ r ∧ r < 0x10000)) ↔
      (r / 65536 = 0 ∧ r / 1024 ≠ 54 ∧ r / 1024 ≠ 55) := by omega
  have supplementary : (0x10000 ≤ r ∧ r ≤ 0x10ffff) ↔
      (1 ≤ r / 65536 ∧ r / 65536 ≤ 16) := by omega
  unfold appendRune
  split
  next inside => rw [ite_eq_left (bmp.mp inside)]
  next outside =>
    rw [ite_eq_right (mt bmp.mpr outside)]
    split
    next inside =>
      have hp := supplementary.mp inside
      rw [ite_eq_left hp, encodeRune_spec, ite_eq_left hp]
    next invalid =>
      rw [ite_eq_right (mt supplementary.mpr invalid)]

end GoLean.Unicode.Utf16
