import GoLean.Unicode.Utf8.FullRune

namespace GoLean.Unicode.Utf8

/-- Go 1.27.0 src/unicode/utf8/utf8.go:129–146 uses the same table,
    bounds, and byte tests as FullRune. The input is raw bytes and may
    encode an invalid string. -/
def fullRuneInString (bytes : List UInt8) : Bool :=
  fullRune (bytes.map UInt8.toNat)

theorem fullRuneInString_empty : fullRuneInString [] = false := rfl

theorem fullRuneInString_short (lead : UInt8) (rest : List UInt8) :
    fullRuneInString (lead :: rest) = false ↔
      rest.length + 1 < (fullLead lead.toNat).1 ∧
      (∀ c, rest[0]? = some c →
        (fullLead lead.toNat).2.1 ≤ c.toNat ∧ c.toNat ≤ (fullLead lead.toNat).2.2) ∧
      (∀ c, rest[1]? = some c → 0x80 ≤ c.toNat ∧ c.toNat ≤ 0xbf) := by
  unfold fullRuneInString
  rw [List.map_cons, fullRune_short]
  cases h₀ : rest[0]? <;> cases h₁ : rest[1]? <;>
    simp [List.getElem?_map, h₀, h₁]

end GoLean.Unicode.Utf8
