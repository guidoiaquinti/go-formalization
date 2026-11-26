import GoLean.Net.Netip.PrefixBits

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:1293-1302 and withoutZone at 505-510. -/
def prefixFrom (ip : ZonedAddr) (bits : Int) : Prefix :=
  let width := addrBitLen ip.detail
  let tag := if ip.detail ≠ .invalid ∧ 0 ≤ bits ∧ bits ≤ Int.ofNat width
    then bits.toNat.toUInt8 + 1 else 0
  let stored := if addrIs6 ip.detail then { ip with zone := "" } else ip
  ⟨stored, tag⟩

theorem prefixFrom_spec (ip : ZonedAddr) (bits : Int) :
    prefixAddr (prefixFrom ip bits) =
      (if ip.detail = .ipv6 then { ip with zone := "" } else ip) ∧
    prefixBits (prefixFrom ip bits) =
      (if ip.detail ≠ .invalid ∧ 0 ≤ bits ∧ bits ≤ Int.ofNat (addrBitLen ip.detail)
        then bits else -1) := by
  constructor
  · simp [prefixAddr, prefixFrom, addrIs6_spec]
  · unfold prefixFrom prefixBits
    dsimp only
    split
    · rename_i h
      have bound : Int.ofNat (addrBitLen ip.detail) ≤ 128 := by
        cases ip.detail <;> decide
      simp [UInt8.toNat_add]
      omega
    · rfl

end GoLean.Net.Netip
