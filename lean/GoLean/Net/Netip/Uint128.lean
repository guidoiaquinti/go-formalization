import Init.Data.BitVec.Basic

namespace GoLean.Net.Netip

structure Uint128 where
  hi : BitVec 64
  lo : BitVec 64
  deriving DecidableEq

/-- Go 1.27.0 src/net/netip/uint128.go:20. Signed negative shift counts panic. -/
def mask6 (n : Int) : Option Uint128 :=
  if n < 0 ∨ 128 < n then none
  else some ⟨~~~((~~~(0 : BitVec 64)) >>> n.toNat),
    (~~~(0 : BitVec 64)) <<< (128 - n).toNat⟩

def uint128IsZero (u : Uint128) : Bool := (u.hi ||| u.lo) == 0
def uint128And (u v : Uint128) : Uint128 := ⟨u.hi &&& v.hi, u.lo &&& v.lo⟩
def uint128Xor (u v : Uint128) : Uint128 := ⟨u.hi ^^^ v.hi, u.lo ^^^ v.lo⟩
def uint128Or (u v : Uint128) : Uint128 := ⟨u.hi ||| v.hi, u.lo ||| v.lo⟩
def uint128Not (u : Uint128) : Uint128 := ⟨~~~u.hi, ~~~u.lo⟩

def uint128SubOne (u : Uint128) : Uint128 :=
  ⟨u.hi - (if u.lo == 0 then 1 else 0), u.lo - 1⟩

def uint128AddOne (u : Uint128) : Uint128 :=
  let low := u.lo + 1
  ⟨u.hi + (if low == 0 then 1 else 0), low⟩

/-- The two returned references designate distinct mutable fields in order.
This model observes a read through each reference and a write back through each. -/
def uint128Halves (u : Option Uint128) (write : Uint128) : Option (Uint128 × Uint128) :=
  u.map fun original => (original, write)

def uint128BitsSetFrom (u : Uint128) (bit : BitVec 8) : Option Uint128 :=
  (mask6 (Int.ofNat bit.toNat)).map fun mask => uint128Or u (uint128Not mask)

def uint128BitsClearedFrom (u : Uint128) (bit : BitVec 8) : Option Uint128 :=
  (mask6 (Int.ofNat bit.toNat)).map fun mask => uint128And u mask

theorem mask6_boundary : mask6 64 = some ⟨0xffffffffffffffff, 0⟩ := by decide
theorem uint128IsZero_zero : uint128IsZero ⟨0, 0⟩ = true := rfl
theorem uint128And_halves (u v : Uint128) :
    uint128And u v = ⟨u.hi &&& v.hi, u.lo &&& v.lo⟩ := rfl
theorem uint128Xor_halves (u v : Uint128) :
    uint128Xor u v = ⟨u.hi ^^^ v.hi, u.lo ^^^ v.lo⟩ := rfl
theorem uint128Or_halves (u v : Uint128) :
    uint128Or u v = ⟨u.hi ||| v.hi, u.lo ||| v.lo⟩ := rfl
theorem uint128Not_halves (u : Uint128) : uint128Not u = ⟨~~~u.hi, ~~~u.lo⟩ := rfl
theorem uint128SubOne_borrow (hi : BitVec 64) :
    uint128SubOne ⟨hi, 0⟩ = ⟨hi - 1, 0xffffffffffffffff⟩ := rfl
theorem uint128AddOne_carry (hi : BitVec 64) :
    uint128AddOne ⟨hi, 0xffffffffffffffff⟩ = ⟨hi + 1, 0⟩ := by
  simp [uint128AddOne]
theorem uint128Halves_reads_and_writes (u v : Uint128) :
    uint128Halves (some u) v = some (u, v) := rfl
theorem uint128BitsSetFrom_all (u : Uint128) :
    uint128BitsSetFrom u 0 = some (uint128Or u ⟨0xffffffffffffffff, 0xffffffffffffffff⟩) := rfl
theorem uint128BitsClearedFrom_all (u : Uint128) :
    uint128BitsClearedFrom u 0 = some (uint128And u ⟨0, 0⟩) := rfl

end GoLean.Net.Netip
