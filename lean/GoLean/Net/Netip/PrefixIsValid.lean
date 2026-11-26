import GoLean.Net.Netip.PrefixBits

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:1312-1315 tests the length tag. -/
def prefixIsValid (p : Prefix) : Bool := decide (p.bitsPlusOne > 0)

theorem prefixIsValid_spec (p : Prefix) :
    prefixIsValid p = true ↔ prefixBits p ≥ 0 := by
  simp [prefixIsValid, prefixBits, UInt8.lt_iff_toNat_lt]
  omega


end GoLean.Net.Netip
