import GoLean.Net.Netip.AddrZone

namespace GoLean.Net.Netip

/-- The comparison called by Less in Go 1.27.0 src/net/netip/netip.go:415-449. -/
def addrCompare (a b : ZonedAddr) : Int :=
  let af := addrBitLen a.detail
  let bf := addrBitLen b.detail
  if af < bf then -1 else if af > bf then 1
  else if a.hi < b.hi then -1 else if a.hi > b.hi then 1
  else if a.lo < b.lo then -1 else if a.lo > b.lo then 1
  else if addrIs6 a.detail then
    if addrZone a < addrZone b then -1 else if addrZone a > addrZone b then 1 else 0
  else 0

def addrOrderedBefore (a b : ZonedAddr) (compareZones : Prop) : Prop :=
  addrBitLen a.detail < addrBitLen b.detail ∨
    (addrBitLen a.detail = addrBitLen b.detail ∧
      (a.hi < b.hi ∨ (a.hi = b.hi ∧
        (a.lo < b.lo ∨ (a.lo = b.lo ∧ compareZones ∧ addrZone a < addrZone b)))))

theorem addrCompare_spec (a b : ZonedAddr) :
    (addrCompare a b = -1 ↔ addrOrderedBefore a b (a.detail = .ipv6)) ∧
    (addrCompare a b = 1 ↔ addrOrderedBefore b a (a.detail = .ipv6)) ∧
    (addrCompare a b = 0 ↔
      ¬addrOrderedBefore a b (a.detail = .ipv6) ∧ ¬addrOrderedBefore b a (a.detail = .ipv6)) := by
  have is6 := addrIs6_spec a.detail
  unfold addrCompare addrOrderedBefore
  repeat' first | split
  all_goals simp_all
  all_goals first | omega | grind

end GoLean.Net.Netip
