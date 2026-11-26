import GoLean.Net.Netip.AddrBitLen

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:457-462 compares the detail handle with z4. -/
def addrIs4 : AddrDetailClass → Bool
  | .ipv4 => true
  | .invalid | .ipv6 => false

theorem addrIs4_spec (detail : AddrDetailClass) :
    addrIs4 detail = true ↔ detail = .ipv4 := by
  cases detail <;> simp [addrIs4]


end GoLean.Net.Netip
