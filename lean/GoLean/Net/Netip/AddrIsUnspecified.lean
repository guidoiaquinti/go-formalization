import GoLean.Net.Netip.AddrUnmap
import GoLean.Net.Netip.IPv4Unspecified
import GoLean.Net.Netip.IPv6Unspecified

namespace GoLean.Net.Netip

def matchesAddr (ip : ZonedAddr) (constant : Addr) : Bool :=
  let detail := if constant.isV6 then AddrDetailClass.ipv6 else .ipv4
  (ip.detail == detail) && (ip.hi == constant.hi) &&
    (ip.lo == constant.lo) && (ip.zone == constant.zone)

/-- Go 1.27.0 src/net/netip/netip.go:670-672 compares with the two unzoned constants.
Compare the observable detail contents, which the Go handles canonicalize. -/
def addrIsUnspecified (ip : ZonedAddr) : Bool :=
  matchesAddr ip ipv4Unspecified || matchesAddr ip ipv6Unspecified

theorem addrIsUnspecified_spec (ip : ZonedAddr) :
    addrIsUnspecified ip = true ↔
      (ip.detail = .ipv4 ∧ ip.hi = 0 ∧ ip.lo = 0xffff00000000 ∧ ip.zone = "") ∨
      (ip.detail = .ipv6 ∧ ip.hi = 0 ∧ ip.lo = 0 ∧ ip.zone = "") := by
  simp [addrIsUnspecified, matchesAddr, ipv4Unspecified_spec, ipv6Unspecified_spec, and_assoc]

end GoLean.Net.Netip
