namespace GoLean.Unicode.Utf8

/-- The signed rune is represented by its mathematical integer value. The
    cases follow Go's order, including rejection of surrogate code points.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L14-L43
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L345-L363 -/
def runeLen (r : Int) : Int :=
  if r < 0 then -1
  else if r ≤ 127 then 1
  else if r ≤ 2047 then 2
  else if 0xd800 ≤ r ∧ r ≤ 0xdfff then -1
  else if r ≤ 65535 then 3
  else if r ≤ 0x10ffff then 4
  else -1

/-- For a valid scalar value, select the first encoding with enough payload
    bits. Planes outside 0–16 and surrogate blocks 54 and 55 are invalid. -/
theorem runeLen_spec (r : Int) :
    runeLen r =
      if r / 65536 < 0 ∨ 17 ≤ r / 65536 ∨ r / 1024 = 54 ∨ r / 1024 = 55 then -1
      else if r / 128 = 0 then 1
      else if r / 2048 = 0 then 2
      else if r / 65536 = 0 then 3
      else 4 := by
  unfold runeLen
  repeat first | split | omega

end GoLean.Unicode.Utf8
