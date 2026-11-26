import Init.Data.UInt.Bitwise
import GoLean.Net.Netip.AddrIs6

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:464-469. -/
def addrIs4In6 (detail : AddrDetailClass) (hi lo : UInt64) : Bool :=
  addrIs6 detail && hi == 0 && (lo >>> 32) == 0xffff

theorem addrIs4In6_spec (detail : AddrDetailClass) (hi lo : UInt64) :
    addrIs4In6 detail hi lo = true ↔
      detail = .ipv6 ∧ hi = 0 ∧ lo.toNat / 2^32 = 0xffff := by
  have mappedPrefix : lo >>> 32 = (0xffff : UInt64) ↔ lo.toNat / 2^32 = 0xffff := by
    constructor
    · intro h
      simpa [Nat.shiftRight_eq_div_pow] using congrArg UInt64.toNat h
    · intro h
      apply UInt64.toNat_inj.1
      simpa [Nat.shiftRight_eq_div_pow] using h
  simp [addrIs4In6, addrIs6_spec, mappedPrefix, and_assoc]

end GoLean.Net.Netip
