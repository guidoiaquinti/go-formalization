import GoLean.Net.Netip.AddrCompare
import GoLean.Net.Netip.AddrPortAddr

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:1146-1151, with cmp.Compare specialized
to uint16. Ports are compared only when the address comparison ties. -/
def addrPortCompare (a b : AddrPort) : Int :=
  let order := addrCompare (addrPortAddr a) (addrPortAddr b)
  if order ≠ 0 then order
  else if addrPortPort a < addrPortPort b then -1
  else if addrPortPort a > addrPortPort b then 1
  else 0

theorem addrPortCompare_spec (a b : AddrPort) :
    let order := addrCompare a.ip b.ip
    (addrPortCompare a b = -1 ↔ order = -1 ∨ (order = 0 ∧ a.port < b.port)) ∧
    (addrPortCompare a b = 1 ↔ order = 1 ∨ (order = 0 ∧ b.port < a.port)) ∧
    (addrPortCompare a b = 0 ↔ order = 0 ∧ a.port = b.port) := by
  simp only [addrPortCompare, addrPortAddr, addrPortPort]
  repeat' first | split
  all_goals simp_all
  all_goals first | omega | grind

end GoLean.Net.Netip
