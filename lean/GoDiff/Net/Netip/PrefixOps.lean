import GoLean.Net.Netip.PrefixOps
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Netip

private def inputAddr (family hi lo zone : Int) : ZonedAddr :=
  ⟨if family = 4 then .ipv4 else if family = 6 then .ipv6 else .invalid,
    UInt64.ofNat hi.toNat, UInt64.ofNat lo.toNat, if zone = 0 then "" else "zone"⟩

private def wire (p : Prefix) : String :=
  let family := match p.ip.detail with | .invalid => 0 | .ipv4 => 4 | .ipv6 => 6
  s!"{family},{p.ip.hi.toNat},{p.ip.lo.toNat},{prefixBits p}"

def main : IO Unit := runTuple fun row => do
  let [family, hi, lo, zone, n, qfamily, qhi, qlo, qzone, m] := row |
    throw (IO.userError "expected two addresses and prefix lengths")
  let a := inputAddr family hi lo zone
  let b := inputAddr qfamily qhi qlo qzone
  let p := prefixFrom a n
  let q := prefixFrom b m
  let address := match addrPrefix a n with
    | .ok ip => wire ip
    | .error err => "error:" ++ err
  return s!"{address}|{wire (prefixMasked p)}|{prefixIsZero p}|{prefixContains p b}|{prefixOverlaps p q}|{prefixCompare p q}"
