import Init.Data.Nat.Bitwise.Lemmas

namespace GoLean.Unicode.Utf16

/-- Go subtracts the supplementary-plane base, then shifts and masks the
    nonnegative payload. The payload and both results fit in signed runes.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf16/utf16.go#L12-L25
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf16/utf16.go#L44-L53 -/
def encodeRune (r : Int) : Nat × Nat :=
  if r < 0x10000 ∨ 0x10ffff < r then
    (0xfffd, 0xfffd)
  else
    let payload := (r - 0x10000).toNat
    (0xd800 + ((payload >>> 10) &&& 0x3ff), 0xdc00 + (payload &&& 0x3ff))

theorem encodeRune_spec (r : Int) :
    encodeRune r =
      if 1 ≤ r / 65536 ∧ r / 65536 ≤ 16 then
        let payload := (r - 65536).toNat
        (55296 + payload / 1024, 56320 + payload % 1024)
      else (65533, 65533) := by
  unfold encodeRune
  split
  next invalid =>
    rw [ite_eq_right (by omega)]
  next valid =>
    rw [ite_eq_left (by omega)]
    dsimp only
    have high : (r - 65536).toNat / 1024 < 1024 := by omega
    rw [Nat.shiftRight_eq_div_pow]
    change
      (55296 + ((r - 65536).toNat / 1024 &&& (2 ^ 10 - 1)),
        56320 + ((r - 65536).toNat &&& (2 ^ 10 - 1))) = _
    rw [Nat.and_two_pow_sub_one_of_lt_two_pow high, Nat.and_two_pow_sub_one_eq_mod]

end GoLean.Unicode.Utf16
