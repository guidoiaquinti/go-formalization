import Init.Data.BitVec.Lemmas

namespace GoLean.Math.Bits

/- The input carry and borrow must be 0 or 1 in Go. Bool enforces that
   precondition and represents the two possible output flags. Results are
   ordered exactly as in Go. All word inputs are unsigned. -/

/-- Add32 widens before splitting the sum. Add64's high-bit logic detects the
    same unsigned overflow. The full sum here has room for the carry bit. -/
private def addWord (x y : BitVec w) (c : Bool) : BitVec w × Bool :=
  let total := x.toNat + y.toNat + c.toNat
  (BitVec.ofNat w total, decide (2 ^ w ≤ total))

def IsSum (x y : BitVec w) (c : Bool) (result : BitVec w × Bool) : Prop :=
  result.1.toNat + (if result.2 then 2 ^ w else 0) = x.toNat + y.toNat + c.toNat

private theorem addWord_spec (x y : BitVec w) (c : Bool) : IsSum x y c (addWord x y c) := by
  have hx := x.isLt
  have hy := y.isLt
  have hc : c.toNat ≤ 1 := by cases c <;> decide
  unfold IsSum addWord
  simp only [BitVec.toNat_ofNat, decide_eq_true_eq]
  split
  · rw [Nat.mod_eq_sub_mod (by omega), Nat.mod_eq_of_lt (by omega)]
    omega
  · rw [Nat.mod_eq_of_lt (by omega), Nat.add_zero]

/-- A borrow supplies one radix when the unsigned subtrahend exceeds x. -/
private def subWord (x y : BitVec w) (b : Bool) : BitVec w × Bool :=
  let debt := y.toNat + b.toNat
  if x.toNat < debt then (BitVec.ofNat w (2 ^ w + x.toNat - debt), true)
  else (BitVec.ofNat w (x.toNat - debt), false)

def IsDifference (x y : BitVec w) (b : Bool) (result : BitVec w × Bool) : Prop :=
  result.1.toNat + y.toNat + b.toNat = x.toNat + (if result.2 then 2 ^ w else 0)

private theorem subWord_spec (x y : BitVec w) (b : Bool) :
    IsDifference x y b (subWord x y b) := by
  have hx := x.isLt
  have hy := y.isLt
  have hb : b.toNat ≤ 1 := by cases b <;> decide
  by_cases h : x.toNat < y.toNat + b.toNat <;> simp [IsDifference, subWord, h] <;>
    rw [Nat.mod_eq_of_lt (by omega)] <;> omega

/-- The halves of a widened product, used by Go's Mul32. -/
private def mulWord (x y : BitVec w) : BitVec w × BitVec w :=
  let product := x.toNat * y.toNat
  (BitVec.ofNat w (product / 2 ^ w), BitVec.ofNat w product)

def IsProduct (x y : BitVec w) (result : BitVec w × BitVec w) : Prop :=
  2 ^ w * result.1.toNat + result.2.toNat = x.toNat * y.toNat

private theorem mulWord_spec (x y : BitVec w) : IsProduct x y (mulWord x y) := by
  have bound := Nat.mul_lt_mul_of_lt_of_le x.isLt (Nat.le_of_lt y.isLt) (Nat.two_pow_pos w)
  have hi := (Nat.div_lt_iff_lt_mul (Nat.two_pow_pos w)).2 bound
  unfold IsProduct mulWord
  simp only [BitVec.toNat_ofNat, Nat.mod_eq_of_lt hi]
  exact Nat.div_add_mod _ _

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L374 -/
def add32 (x y : BitVec 32) (c : Bool) := addWord x y c

theorem add32_spec (x y : BitVec 32) (c : Bool) : IsSum x y c (add32 x y c) :=
  addWord_spec x y c

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L386 -/
def add64 (x y : BitVec 64) (c : Bool) := addWord x y c

theorem add64_spec (x y : BitVec 64) (c : Bool) : IsSum x y c (add64 x y c) :=
  addWord_spec x y c

/-- Add on the pinned 64-bit Go profile.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L360 -/
def add (x y : BitVec 64) (c : Bool) := add64 x y c

theorem add_spec (x y : BitVec 64) (c : Bool) : IsSum x y c (add x y c) :=
  addWord_spec x y c

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L416 -/
def sub32 (x y : BitVec 32) (b : Bool) := subWord x y b

theorem sub32_spec (x y : BitVec 32) (b : Bool) : IsDifference x y b (sub32 x y b) :=
  subWord_spec x y b

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L431 -/
def sub64 (x y : BitVec 64) (b : Bool) := subWord x y b

theorem sub64_spec (x y : BitVec 64) (b : Bool) : IsDifference x y b (sub64 x y b) :=
  subWord_spec x y b

/-- Sub on the pinned 64-bit Go profile.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L402 -/
def sub (x y : BitVec 64) (b : Bool) := sub64 x y b

theorem sub_spec (x y : BitVec 64) (b : Bool) : IsDifference x y b (sub x y b) :=
  subWord_spec x y b

/-- https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L459 -/
def mul32 (x y : BitVec 32) := mulWord x y

theorem mul32_spec (x y : BitVec 32) : IsProduct x y (mul32 x y) :=
  mulWord_spec x y

private abbrev limbBase : Nat := 0x100000000

/-- Go's Mul64 uses 32-bit limbs. Remainder by limbBase is its mask32,
    division by limbBase its unsigned right shift, and reduction modulo 2^64
    gives the uint64 intermediates. The low output is the wrapping product.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L470-L484 -/
def mul64 (x y : BitVec 64) : BitVec 64 × BitVec 64 :=
  let x0 := x.toNat % limbBase
  let x1 := x.toNat / limbBase
  let y0 := y.toNat % limbBase
  let y1 := y.toNat / limbBase
  let w0 := x0 * y0
  let t := (x1 * y0 + w0 / limbBase) % 2^64
  let w1 := (t % limbBase + x0 * y1) % 2^64
  let w2 := t / limbBase
  (BitVec.ofNat 64 (x1 * y1 + w2 + w1 / limbBase), x * y)

theorem mul64_spec (x y : BitVec 64) : IsProduct x y (mul64 x y) := by
  let x0 := x.toNat % limbBase
  let x1 := x.toNat / limbBase
  let y0 := y.toNat % limbBase
  let y1 := y.toNat / limbBase
  have hx := x.isLt
  have hy := y.isLt
  have bound (a b : Nat) (ha : a < limbBase) (hb : b < limbBase) :
      a * b ≤ (limbBase - 1) * (limbBase - 1) :=
    Nat.mul_le_mul (by omega) (by omega)
  have hx0 : x0 < limbBase := Nat.mod_lt _ (by decide)
  have hy0 : y0 < limbBase := Nat.mod_lt _ (by decide)
  have hx1 : x1 < limbBase := by dsimp [x1, limbBase]; omega
  have hy1 : y1 < limbBase := by dsimp [y1, limbBase]; omega
  let w0 := x0 * y0
  let t := x1 * y0 + w0 / limbBase
  let w1 := t % limbBase + x0 * y1
  let high := x1 * y1 + t / limbBase + w1 / limbBase
  have ht : t < 2^64 := by
    have := bound x1 y0 hx1 hy0
    have := bound x0 y0 hx0 hy0
    dsimp [t, w0, limbBase] at *
    omega
  have hw1 : w1 < 2^64 := by
    have := bound x0 y1 hx0 hy1
    dsimp [w1, limbBase] at *
    omega
  have expand : x.toNat * y.toNat =
      w0 + limbBase * (x1 * y0 + x0 * y1) + limbBase^2 * (x1 * y1) := by
    calc
      _ = (x0 + limbBase * x1) * (y0 + limbBase * y1) :=
        by rw [show x0 + limbBase * x1 = _ from Nat.mod_add_div ..,
               show y0 + limbBase * y1 = _ from Nat.mod_add_div ..]
      _ = _ := by dsimp [w0]; grind only
  have productBound := Nat.mul_lt_mul_of_lt_of_le hx (Nat.le_of_lt hy) (Nat.two_pow_pos 64)
  have quotient : high = (x.toNat * y.toNat) / 2^64 := by
    rw [expand]
    dsimp [high, w1, t, limbBase]
    omega
  have highFits : high < 2^64 := quotient ▸ (Nat.div_lt_iff_lt_mul (Nat.two_pow_pos 64)).2 productBound
  change 2^64 * _ + _ = _
  simp only [mul64, BitVec.toNat_ofNat, BitVec.toNat_mul]
  change 2^64 * ((x1*y1 + t % 2^64 / limbBase +
      ((t % 2^64 % limbBase + x0*y1) % 2^64) / limbBase) % 2^64) + _ = _
  rw [Nat.mod_eq_of_lt ht, Nat.mod_eq_of_lt hw1, Nat.mod_eq_of_lt highFits,
    quotient, Nat.div_add_mod]

/-- Straight-line model of the pinned Go source. Each variable and
    operation has the unsigned 64-bit semantics of its Go counterpart.
    The untyped Go constant `1<<32 - 1` is exact before conversion. -/
def mul64Source (x y : BitVec 64) : BitVec 64 × BitVec 64 :=
  let mask32 : BitVec 64 := 0xffffffff
  let x0 := x &&& mask32
  let x1 := x >>> (32 : Nat)
  let y0 := y &&& mask32
  let y1 := y >>> (32 : Nat)
  let w0 := x0 * y0
  let t := x1 * y0 + (w0 >>> (32 : Nat))
  let w1 := t &&& mask32
  let w2 := t >>> (32 : Nat)
  let w1 := w1 + x0 * y1
  let hi := x1 * y1 + w2 + (w1 >>> (32 : Nat))
  let lo := x * y
  (hi, lo)

theorem mul64Source_eq_mul64 (x y : BitVec 64) :
    mul64Source x y = mul64 x y := by
  have mask : (0xffffffff : BitVec 64).toNat = 2^32 - 1 := by decide
  have limbProductFits {a b : Nat} (ha : a < 2^32) (hb : b < 2^32) :
      a * b < 2^64 := by
    have fit : a * b < (2^32) * (2^32) :=
      Nat.mul_lt_mul_of_lt_of_le ha (Nat.le_of_lt hb) (by decide)
    exact fit
  have lowProductFits : x.toNat % 2^32 * (y.toNat % 2^32) < 2^64 :=
    limbProductFits (Nat.mod_lt _ (by decide)) (Nat.mod_lt _ (by decide))
  simp only [mul64Source, mul64, Prod.mk.injEq, and_true]
  apply BitVec.eq_of_toNat_eq
  simp only [BitVec.toNat_and, BitVec.toNat_ushiftRight, BitVec.toNat_mul,
    BitVec.toNat_add, BitVec.toNat_ofNat, mask, Nat.and_two_pow_sub_one_eq_mod,
    Nat.shiftRight_eq_div_pow, Nat.mod_eq_of_lt lowProductFits, Nat.add_mod_mod,
    Nat.mod_add_mod]

/-- The exact word operations of the source model return the full product. -/
theorem mul64Source_spec (x y : BitVec 64) : IsProduct x y (mul64Source x y) :=
  mul64Source_eq_mul64 x y ▸ mul64_spec x y

/-- Mul on the pinned 64-bit Go profile.
    https://github.com/golang/go/blob/go1.27.0/src/math/bits/bits.go#L445 -/
def mul (x y : BitVec 64) := mul64 x y

theorem mul_spec (x y : BitVec 64) : IsProduct x y (mul x y) :=
  mul64_spec x y

end GoLean.Math.Bits
