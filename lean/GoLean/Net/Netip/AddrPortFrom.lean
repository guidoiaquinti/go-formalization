import GoLean.Net.Netip.AddrPortPort

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:1058-1065. -/
def addrPortFrom (ip : ZonedAddr) (port : UInt16) : AddrPort :=
  { ip := ip, port := port }

theorem addrPortFrom_spec (ip : ZonedAddr) (port : UInt16) :
    (addrPortFrom ip port).ip = ip ∧ (addrPortFrom ip port).port = port :=
  ⟨rfl, rfl⟩


end GoLean.Net.Netip
