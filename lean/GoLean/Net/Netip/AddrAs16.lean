import Init.Data.UInt.Bitwise

namespace GoLean.Net.Netip

/-- Go 1.27.0 src/net/netip/netip.go:701-710 and
src/internal/byteorder/byteorder.go:126-136. Each output byte is a shifted word
truncated to eight bits. -/
def addrAs16 (hi lo : UInt64) (i : Fin 16) : UInt8 :=
  let word := if i.val < 8 then hi else lo
  let shift := 8 * (7 - i.val % 8)
  (word >>> UInt64.ofNat shift).toUInt8

theorem addrAs16_spec (hi lo : UInt64) (i : Fin 16) :
    (addrAs16 hi lo i).toNat =
      (if i.val < 8 then hi.toNat else lo.toNat) / 256 ^ (7 - i.val % 8) % 256 := by
  have shiftBound : 8 * (7 - i.val % 8) < 64 := by omega
  have powEq : 2 ^ (8 * (7 - i.val % 8)) = 256 ^ (7 - i.val % 8) := by
    rw [Nat.pow_mul]
  simp [addrAs16, Nat.shiftRight_eq_div_pow, Nat.mod_eq_of_lt shiftBound, powEq, apply_ite]

end GoLean.Net.Netip
