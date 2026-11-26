import GoLean.Net.Netip.AddrCompare

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:452-455. -/
def addrLess (a b : ZonedAddr) : Bool := addrCompare a b == -1

theorem addrLess_spec (a b : ZonedAddr) :
    addrLess a b = true ↔
      addrBitLen a.detail < addrBitLen b.detail ∨
      (addrBitLen a.detail = addrBitLen b.detail ∧
        (a.hi < b.hi ∨ (a.hi = b.hi ∧
          (a.lo < b.lo ∨ (a.lo = b.lo ∧ a.detail = .ipv6 ∧ addrZone a < addrZone b))))) := by
  simpa only [addrLess, beq_iff_eq, addrOrderedBefore] using (addrCompare_spec a b).1

end GoLean.Net.Netip
