import GoLean.Net.Netip.PrefixIsValid
import GoLean.Net.Netip.PrefixFrom

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:1319-1320. -/
def prefixIsSingleIP (p : Prefix) : Bool :=
  prefixIsValid p && (prefixBits p == Int.ofNat (addrBitLen p.ip.detail))

theorem prefixIsSingleIP_spec (ip : ZonedAddr) (bits : Int) :
    prefixIsSingleIP (prefixFrom ip bits) = true ↔
      ip.detail ≠ .invalid ∧ bits = Int.ofNat (addrBitLen ip.detail) := by
  rcases prefixFrom_spec ip bits with ⟨stored, length⟩
  have detail : (prefixFrom ip bits).ip.detail = ip.detail := by
    change (prefixAddr (prefixFrom ip bits)).detail = _
    rw [stored]
    split <;> rfl
  simp only [prefixIsSingleIP, Bool.and_eq_true, beq_iff_eq, prefixIsValid_spec,
    detail, length]
  split
  · rename_i h
    exact ⟨fun hr => ⟨h.1, hr.2⟩, fun hr => ⟨h.2.1, hr.2⟩⟩
  · rename_i h
    constructor
    · intro hr
      omega
    · rintro ⟨valid, rfl⟩
      exact (h ⟨valid, Int.natCast_nonneg _, Int.le_refl _⟩).elim

end GoLean.Net.Netip
