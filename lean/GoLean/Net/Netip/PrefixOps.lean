import GoLean.Net.Netip.PrefixFrom
import GoLean.Net.Netip.AddrCompare
import GoLean.Net.Netip.Uint128

namespace GoLean.Net.Netip

private def zeroAddr : ZonedAddr := ⟨.invalid, 0, 0, ""⟩
private def zeroPrefix : Prefix := ⟨zeroAddr, 0⟩

private def addressWord (ip : ZonedAddr) : Uint128 := ⟨ip.hi.toBitVec, ip.lo.toBitVec⟩

/-- Go 1.27.0 src/net/netip/netip.go:679. IPv4 applies the prefix after
its stored 96-bit mapped marker. Negative lengths are rejected first. -/
def addrPrefix (ip : ZonedAddr) (bits : Int) : Except String Prefix :=
  if bits < 0 then .error "negative Prefix bits"
  else if ip.detail = .invalid then .ok zeroPrefix
  else
    let width := addrBitLen ip.detail
    if bits > width then
      let family := if ip.detail = .ipv4 then "IPv4" else "IPv6"
      .error s!"prefix length {bits} too large for {family}"
    else
      let maskBits := if ip.detail = .ipv4 then bits + 96 else bits
      let mask := (mask6 maskBits).getD ⟨0, 0⟩
      let masked := uint128And (addressWord ip) mask
      .ok (prefixFrom {ip with hi := UInt64.ofBitVec masked.hi, lo := UInt64.ofBitVec masked.lo} bits)

def prefixMasked (p : Prefix) : Prefix :=
  match addrPrefix p.ip (prefixBits p) with
  | .ok masked => masked
  | .error _ => zeroPrefix

def prefixIsZero (p : Prefix) : Bool :=
  p.bitsPlusOne == 0 && p.ip.detail == .invalid && p.ip.hi == 0 && p.ip.lo == 0 && p.ip.zone.isEmpty

def prefixContains (p : Prefix) (ip : ZonedAddr) : Bool :=
  if p.bitsPlusOne == 0 || !ip.zone.isEmpty || ip.detail == .invalid || p.ip.detail != ip.detail then false
  else
    let mismatch := uint128Xor (addressWord p.ip) (addressWord ip)
    let bits := prefixBits p
    if ip.detail = .ipv4 then
      ((mismatch.lo >>> (32 - bits).toNat).setWidth 32) == 0
    else
      match mask6 bits with
      | none => false
      | some mask => uint128IsZero (uint128And mismatch mask)

def prefixOverlaps (p q : Prefix) : Bool :=
  if p.bitsPlusOne == 0 || q.bitsPlusOne == 0 || p.ip.detail != q.ip.detail then false
  else
    let n := min (prefixBits p) (prefixBits q)
    if n = 0 then true
    else
      match addrPrefix p.ip n, addrPrefix q.ip n with
      | .ok a, .ok b => addrCompare a.ip b.ip == 0
      | _, _ => false

def prefixCompare (p q : Prefix) : Int :=
  let masked := addrCompare (prefixMasked p).ip (prefixMasked q).ip
  if masked != 0 then masked
  else if prefixBits p < prefixBits q then -1
  else if prefixBits q < prefixBits p then 1
  else addrCompare p.ip q.ip

theorem addrPrefix_negative (ip : ZonedAddr) (n : Int) (h : n < 0) :
    addrPrefix ip n = .error "negative Prefix bits" := by simp [addrPrefix, h]

theorem prefixMasked_zero : prefixMasked zeroPrefix = zeroPrefix := rfl
theorem prefixIsZero_zero : prefixIsZero zeroPrefix = true := rfl

theorem prefixContains_invalid (p : Prefix) (h : p.bitsPlusOne = 0) (ip : ZonedAddr) :
    prefixContains p ip = false := by simp [prefixContains, h]

theorem prefixOverlaps_invalid (p q : Prefix) (h : p.bitsPlusOne = 0) :
    prefixOverlaps p q = false := by simp [prefixOverlaps, h]

theorem prefixCompare_zero : prefixCompare zeroPrefix zeroPrefix = 0 := rfl

end GoLean.Net.Netip
