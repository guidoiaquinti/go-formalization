import GoLean.Net.Netip.AddrBitLen

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:385-390 compares the address detail with z0. -/
def addrIsValid : AddrDetailClass → Bool
  | .invalid => false
  | .ipv4 | .ipv6 => true

theorem addrIsValid_spec (detail : AddrDetailClass) :
    addrIsValid detail = true ↔ detail ≠ .invalid := by
  cases detail <;> simp [addrIsValid]


end GoLean.Net.Netip
