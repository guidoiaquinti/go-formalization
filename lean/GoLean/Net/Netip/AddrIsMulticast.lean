import GoLean.Net.Netip.AddrUnmap
import GoLean.Net.Netip.AddrAs16
import GoLean.Net.Netip.AddrIs4

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:557-573.
The v4(0) shift and truncation at 360-362 selects byte 12 of the stored words. -/
def addrIsMulticast (input : ZonedAddr) : Bool :=
  let ip := if addrIs4In6 input.detail input.hi input.lo then addrUnmap input else input
  if addrIs4 ip.detail then (addrAs16 ip.hi ip.lo 12 &&& 0xf0) == 0xe0
  else if addrIs6 ip.detail then (ip.hi >>> 56) == 0xff
  else false

theorem addrIsMulticast_spec (ip : ZonedAddr) :
    addrIsMulticast ip = true ↔
      ((ip.detail = .ipv4 ∨
          (ip.detail = .ipv6 ∧ ip.hi = 0 ∧ ip.lo.toNat / 2^32 = 0xffff)) ∧
        (addrAs16 ip.hi ip.lo 12 &&& 0xf0) = 0xe0) ∨
      (ip.detail = .ipv6 ∧ ip.hi.toNat / 2^56 = 255) := by
  have top : ip.hi >>> 56 = (255 : UInt64) ↔ ip.hi.toNat / 2^56 = 255 := by
    rw [← UInt64.toNat_inj]
    simp [Nat.shiftRight_eq_div_pow]
  cases h : ip.detail <;>
    simp [addrIsMulticast, addrUnmap, addrIs4In6_spec, addrIs4, addrIs6, h]
  by_cases high : ip.hi = 0 <;> by_cases low : ip.lo.toNat / 2^32 = 0xffff <;>
    simp_all

end GoLean.Net.Netip
