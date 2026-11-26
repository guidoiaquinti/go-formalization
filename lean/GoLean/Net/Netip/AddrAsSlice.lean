import GoLean.Net.Netip.AddrAs16
import GoLean.Net.Netip.AddrBitLen

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:726-741.
None represents nil. The populated branches return the full fixed-size buffer. -/
def addrAsSlice (detail : AddrDetailClass) (hi lo : UInt64) : Option (List UInt8) :=
  match detail with
  | .invalid => none
  | .ipv4 => some (List.ofFn fun i : Fin 4 => addrAs16 hi lo ⟨12 + i.val, by omega⟩)
  | .ipv6 => some (List.ofFn (addrAs16 hi lo))

theorem addrAsSlice_spec (detail : AddrDetailClass) (hi lo : UInt64) :
    (addrAsSlice detail hi lo).map (List.map UInt8.toNat) =
      match detail with
      | .invalid => none
      | .ipv4 => some (List.ofFn fun i : Fin 4 => lo.toNat / 256 ^ (3 - i.val) % 256)
      | .ipv6 => some (List.ofFn fun i : Fin 16 =>
          (if i.val < 8 then hi.toNat else lo.toNat) / 256 ^ (7 - i.val % 8) % 256) := by
  cases detail
  · rfl
  · simp only [addrAsSlice, Option.map_some, List.map_ofFn, Option.some.injEq]
    congr 1
    funext i
    rw [Function.comp_apply, addrAs16_spec]
    have offset : (12 + i.val) % 8 = 4 + i.val := by omega
    have exponent : 7 - (4 + i.val) = 3 - i.val := by omega
    simp [offset, exponent, show ¬12 + i.val < 8 by omega]
  · simp only [addrAsSlice, Option.map_some, List.map_ofFn, Option.some.injEq]
    congr 1
    funext i
    exact addrAs16_spec hi lo i

end GoLean.Net.Netip
