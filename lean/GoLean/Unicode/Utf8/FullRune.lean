namespace GoLean.Unicode.Utf8

/-- Width and second-byte limits from the lead-byte and accept tables.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L48-L106 -/
def fullLead (b : Nat) : Nat × Nat × Nat :=
  if 0xc2 ≤ b ∧ b ≤ 0xdf then (2, 0x80, 0xbf)
  else if 0xe0 ≤ b ∧ b ≤ 0xef then
    if b = 0xe0 then (3, 0xa0, 0xbf)
    else if b = 0xed then (3, 0x80, 0x9f)
    else (3, 0x80, 0xbf)
  else if 0xf0 ≤ b ∧ b ≤ 0xf4 then
    if b = 0xf0 then (4, 0x90, 0xbf)
    else if b = 0xf4 then (4, 0x80, 0x8f)
    else (4, 0x80, 0xbf)
  else (1, 0x80, 0xbf)

/-- A malformed sequence is full as soon as it is known to decode as a
    one-byte error. The length check precedes the continuation checks.
    https://github.com/golang/go/blob/go1.27.0/src/unicode/utf8/utf8.go#L108-L127 -/
def fullRune : List Nat → Bool
  | [] => false
  | b :: rest =>
    let (size, lo, hi) := fullLead b
    if rest.length + 1 ≥ size then true
    else if (rest[0]?).any (fun c => c < lo || hi < c) then true
    else (rest[1]?).any (fun c => c < 0x80 || 0xbf < c)

theorem fullRune_empty : fullRune [] = false := rfl

theorem fullRune_short (b : Nat) (rest : List Nat) :
    fullRune (b :: rest) = false ↔
      rest.length + 1 < (fullLead b).1 ∧
      (∀ c, rest[0]? = some c → (fullLead b).2.1 ≤ c ∧ c ≤ (fullLead b).2.2) ∧
      (∀ c, rest[1]? = some c → 0x80 ≤ c ∧ c ≤ 0xbf) := by
  rcases h : fullLead b with ⟨size, lo, hi⟩
  cases h₀ : rest[0]? <;> cases h₁ : rest[1]? <;>
    simp [fullRune, h, h₀, h₁, Bool.or_eq_false_iff] <;> omega

end GoLean.Unicode.Utf8
