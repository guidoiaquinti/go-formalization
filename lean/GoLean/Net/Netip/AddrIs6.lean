import GoLean.Net.Netip.AddrBitLen

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:471-475 tests z0 before z4.
The detail classes encode which sentinel comparisons succeed. -/
def addrIs6 (detail : AddrDetailClass) : Bool :=
  let isZero := match detail with
    | .invalid => true
    | _ => false
  let isIPv4 := match detail with
    | .ipv4 => true
    | _ => false
  !isZero && !isIPv4

theorem addrIs6_spec (detail : AddrDetailClass) :
    addrIs6 detail = true ↔ detail = .ipv6 := by
  cases detail <;> simp [addrIs6]


end GoLean.Net.Netip
