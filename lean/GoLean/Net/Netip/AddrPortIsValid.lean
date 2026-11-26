import GoLean.Net.Netip.AddrPortPort
import GoLean.Net.Netip.AddrIsValid

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:1141 delegates to Addr.IsValid. -/
def addrPortIsValid (p : AddrPort) : Bool := addrIsValid p.ip.detail

theorem addrPortIsValid_spec (p : AddrPort) :
    addrPortIsValid p = true ↔ p.ip.detail ≠ .invalid :=
  addrIsValid_spec p.ip.detail


end GoLean.Net.Netip
