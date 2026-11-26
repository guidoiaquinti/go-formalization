namespace GoLean.Net.Netip

/-- The three reachable address-detail classes tested by BitLen.
IPv6 includes both scoped and unscoped addresses. -/
inductive AddrDetailClass where
  | invalid
  | ipv4
  | ipv6
  deriving DecidableEq

/-- Go 1.27.0 src/net/netip/netip.go:65-71,397-405. -/
def addrBitLen : AddrDetailClass → Nat
  | .invalid => 0
  | .ipv4 => 32
  | .ipv6 => 128

theorem addrBitLen_spec :
    addrBitLen .invalid = 0 ∧ addrBitLen .ipv4 = 32 ∧ addrBitLen .ipv6 = 128 :=
  ⟨rfl, rfl, rfl⟩

end GoLean.Net.Netip
