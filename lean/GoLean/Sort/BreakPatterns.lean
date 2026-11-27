import Init.Data.BitVec.Basic

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/sort.go:70. The returned word is also the updated state. -/
def xorshiftNext (state : BitVec 64) : BitVec 64 :=
  let left := state ^^^ (state <<< 13)
  let right := left ^^^ (left >>> 7)
  right ^^^ (right <<< 17)

theorem xorshiftNext_zero : xorshiftNext 0 = 0 := by decide

theorem xorshiftNext_one : xorshiftNext 1 = 1082269761 := by decide

/-- Go 1.27.0 src/sort/sort.go:77, operating on the unsigned image of the int. -/
def nextPowerOfTwo (word : BitVec width) : BitVec width :=
  let n := word.toNat
  let shift := if n = 0 then 0 else n.log2 + 1
  (1 : BitVec width) <<< shift

theorem nextPowerOfTwo_zero : nextPowerOfTwo (0 : BitVec 64) = 1 := by decide

theorem nextPowerOfTwo_strictAtEight : nextPowerOfTwo (8 : BitVec 64) = 16 := by decide

theorem nextPowerOfTwo_signedMin : nextPowerOfTwo (9223372036854775808 : BitVec 64) = 0 := by decide

/-- Go 1.27.0 src/sort/zsortinterface.go:240. Valid nonnegative bounds give
three Swap arguments around the midpoint for ranges at least eight wide. -/
def breakPatterns (width a b : Nat) : List (Nat × Nat) :=
  let size := b - a
  if size < 8 then [] else
    let mask := nextPowerOfTwo (BitVec.ofNat width size) - 1
    let first := a + (size / 4) * 2 - 1
    let state := BitVec.ofNat 64 size
    let random1 := xorshiftNext state
    let random2 := xorshiftNext random1
    let random3 := xorshiftNext random2
    let destination (random : BitVec 64) :=
      let candidate := (BitVec.ofNat width random.toNat &&& mask).toNat
      a + (if candidate ≥ size then candidate - size else candidate)
    [(first, destination random1), (first + 1, destination random2), (first + 2, destination random3)]

def breakPatternsFunc := breakPatterns

theorem breakPatterns_short (width a b : Nat) (h : b - a < 8) :
    breakPatterns width a b = [] := by simp [breakPatterns, h]

theorem breakPatterns_count (width a b : Nat) (h : 8 ≤ b - a) :
    (breakPatterns width a b).length = 3 := by simp [breakPatterns, Nat.not_lt.mpr h]

theorem breakPatternsFunc_same (width a b : Nat) :
    breakPatternsFunc width a b = breakPatterns width a b := rfl

end GoLean.Sort
