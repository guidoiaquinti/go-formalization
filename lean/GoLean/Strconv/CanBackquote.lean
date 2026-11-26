namespace GoLean.Strconv

/-- Go 1.27.0 src/strconv/quote.go:208-226. Decode the entire byte string
before testing the runes. Any malformed sequence makes the result false. -/
def canBackquote (input : ByteArray) : Bool :=
  match String.fromUTF8? input with
  | none => false
  | some text => text.toList.all fun c =>
      if c.toNat ≥ 128 then c.toNat != 0xfeff
      else (c.toNat ≥ 32 || c.toNat == 9) && c.toNat != 96 && c.toNat != 127

theorem canBackquote_spec (input : ByteArray) :
    canBackquote input = true ↔
      ∃ text, String.fromUTF8? input = some text ∧
        ∀ c ∈ text.toList,
          (32 ≤ c.toNat ∨ c.toNat = 9) ∧
          c.toNat ≠ 96 ∧ c.toNat ≠ 127 ∧ c.toNat ≠ 0xfeff := by
  cases h : String.fromUTF8? input with
  | none => simp [canBackquote, h]
  | some text =>
    simp only [canBackquote, h, Option.some.injEq, exists_eq_left', List.all_eq_true]
    apply forall_congr'
    intro c
    apply forall_congr'
    intro _
    split <;> simp <;> omega

end GoLean.Strconv
