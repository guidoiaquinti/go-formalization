namespace GoLean.Slices

/-- Go 1.27.0 slices/sort.go:181–185. The returned value is the updated uint64 state. -/
def xorshiftNext (state : Option (BitVec 64)) : Option (BitVec 64 × BitVec 64) :=
  state.map fun (x : BitVec 64) =>
    let x := x ^^^ (x <<< 13)
    let x := x ^^^ (x >>> 7)
    let x := x ^^^ (x <<< 17)
    (x, x)

/-- Go 1.27.0 slices/sort.go:188–190. The shift and result use native uint width. -/
def nextPowerOfTwo (width : Nat) (length : Int) : BitVec width :=
  let word := (BitVec.ofInt width length).toNat
  let bitLength := if word = 0 then 0 else word.log2 + 1
  (1 : BitVec width) <<< bitLength

/-- Go 1.27.0 slices/sort.go:194–196. Self-inequality is represented by the
    IEEE exponent and payload predicate for binary32 and binary64. -/
def sortIsNaN (fractionWidth exponentWidth : Nat) (floating : Bool) (word : Nat) : Bool :=
  floating && (word / 2 ^ fractionWidth % 2 ^ exponentWidth == 2 ^ exponentWidth - 1) &&
    (word % 2 ^ fractionWidth != 0)

theorem xorshiftNext_zero : xorshiftNext (some 0) = some (0, 0) := rfl

theorem xorshiftNext_result_eq_state (state : Option (BitVec 64)) :
    (xorshiftNext state).map Prod.fst = (xorshiftNext state).map Prod.snd := by
  cases state <;> rfl

theorem nextPowerOfTwo_boundaries :
    nextPowerOfTwo 64 0 = 1 ∧ nextPowerOfTwo 64 32 = 64 ∧
    nextPowerOfTwo 64 (-1) = 0 ∧ nextPowerOfTwo 32 (-1) = 0 := by decide

theorem sortIsNaN_nonfloat (fractionWidth exponentWidth word : Nat) :
    sortIsNaN fractionWidth exponentWidth false word = false := rfl

theorem sortIsNaN_binary_boundaries :
    sortIsNaN 23 8 true 0x7f800001 = true ∧
    sortIsNaN 23 8 true 0xff800000 = false ∧
    sortIsNaN 52 11 true 0xfff0000000000001 = true ∧
    sortIsNaN 52 11 true 0x7ff0000000000000 = false := by decide

end GoLean.Slices
