import Init.Data.Nat.Bitwise.Lemmas

namespace GoLean.Unicode.Utf16

/-- Go tests the two rune ranges and combines the ten-bit payloads.
    The payloads and their sum with 0x10000 fit in a signed Go rune.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf16/utf16.go#L12-L42 -/
def decodeRune (first second : Int) : Nat :=
  if 0xd800 ≤ first ∧ first < 0xdc00 ∧ 0xdc00 ≤ second ∧ second < 0xe000 then
    (((first - 0xd800).toNat <<< 10) ||| (second - 0xdc00).toNat) + 0x10000
  else
    0xfffd

private theorem combinePayload (high low : Nat) (h : low < 1024) :
    (high <<< 10 ||| low) = high * 1024 + low := by
  rw [Nat.shiftLeft_eq]
  have quotient : (high * 2 ^ 10 ||| low) / 2 ^ 10 = high := by
    rw [Nat.or_div_two_pow]
    simp [Nat.div_eq_of_lt h]
  have remainder : (high * 2 ^ 10 ||| low) % 2 ^ 10 = low := by
    rw [Nat.or_mod_two_pow]
    simp [Nat.mod_eq_of_lt h]
  omega

theorem decodeRune_spec (first second : Int) :
    decodeRune first second =
      if first / 1024 = 54 ∧ second / 1024 = 55 then
        65536 + 1024 * (first - 0xd800).toNat + (second - 0xdc00).toNat
      else 65533 := by
  have valid : (0xd800 ≤ first ∧ first < 0xdc00 ∧ 0xdc00 ≤ second ∧ second < 0xe000) ↔
      (first / 1024 = 54 ∧ second / 1024 = 55) := by omega
  unfold decodeRune
  split
  next h =>
    rw [ite_eq_left (valid.mp h), combinePayload _ _ (by omega)]
    omega
  next h =>
    rw [ite_eq_right (mt valid.mpr h)]

end GoLean.Unicode.Utf16
