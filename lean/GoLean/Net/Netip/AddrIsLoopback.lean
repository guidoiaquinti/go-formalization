import GoLean.Net.Netip.AddrUnmap
import GoLean.Net.Netip.AddrAs16
import GoLean.Net.Netip.AddrIs4

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:538-554.
The v4(0) shift and truncation at 360-362 selects byte 12 of the stored words. -/
def addrIsLoopback (input : ZonedAddr) : Bool :=
  let ip := if addrIs4In6 input.detail input.hi input.lo then addrUnmap input else input
  if addrIs4 ip.detail then addrAs16 ip.hi ip.lo 12 == 127
  else if addrIs6 ip.detail then ip.hi == 0 && ip.lo == 1
  else false

theorem addrIsLoopback_spec (ip : ZonedAddr) :
    addrIsLoopback ip = true ↔
      ((ip.detail = .ipv4 ∨
          (ip.detail = .ipv6 ∧ ip.hi = 0 ∧ ip.lo.toNat / 2^32 = 0xffff)) ∧
        ip.lo.toNat / 256^3 % 256 = 127) ∨
      (ip.detail = .ipv6 ∧ ip.hi = 0 ∧ ip.lo = 1) := by
  have octet : addrAs16 ip.hi ip.lo 12 = 127 ↔ ip.lo.toNat / 256^3 % 256 = 127 := by
    rw [← UInt8.toNat_inj, addrAs16_spec]
    rfl
  cases h : ip.detail <;>
    simp [addrIsLoopback, addrUnmap, addrIs4In6_spec, addrIs4, addrIs6, h, octet]
  by_cases high : ip.hi = 0 <;> by_cases low : ip.lo.toNat / 2^32 = 0xffff <;>
    simp_all
  intro one
  simp [one] at low

end GoLean.Net.Netip
