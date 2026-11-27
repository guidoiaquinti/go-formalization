import GoLean.Sort.Pdqsort
import GoLean.Math.IsNaN
import Init.Data.List.Lex

namespace GoLean.Sort

/-- Go 1.27.0 src/sort/sort.go:148 and src/math/bits.go use the same
floating self-inequality test. -/
def sortIsNaN := GoLean.Math.isNaN

theorem sortIsNaN_encoding (x : BitVec 64) : sortIsNaN x = true ↔
    (0x7ff0000000000000 < x.toNat ∧ x.toNat < 0x8000000000000000) ∨
      0xfff0000000000000 < x.toNat := GoLean.Math.isNaN_spec x

private def floatKey (word : BitVec 64) : Int :=
  if sortIsNaN word then -1 else
    let magnitude := word.toNat % 2^63
    if word.toNat < 2^63 then Int.ofNat (2^63 + magnitude)
    else Int.ofNat (2^63 - magnitude)

private def stringKeys (words : Array (List UInt8)) : Array Int :=
  words.map fun word =>
    Int.ofNat (words.foldl (fun count other => if other.lex word then count + 1 else count) 0)

private def sortedPayload [Inhabited α] (width : Nat) (keys : Array Int) (items : Array α) :
    Array α :=
  if items.size ≤ 1 then items else
    let result := pdqsort width ⟨keys, []⟩ 0 items.size (items.size.log2 + 1)
    result.reversedCalls.reverse.foldl (fun xs (swap, i, j) =>
      if swap then (xs.set! i xs[j]!).set! j xs[i]! else xs) items

/-- Go 1.27.0 src/sort/sort.go:170 delegates to slices.Sort. Rank sorting
uses the ordered PDQ path; replaying its swaps preserves payload identity. -/
def sortInts (width : Nat) (xs : Array Int) : Array Int :=
  sortedPayload width xs xs

/-- The generic slices.Sort and Interface form compare binary64 values
with all NaNs first and both zero encodings equivalent. -/
def sortFloat64s (width : Nat) (xs : Array (BitVec 64)) : Array (BitVec 64) :=
  sortedPayload width (xs.map floatKey) xs

def float64SliceSort := sortFloat64s

/-- Raw Go strings compare in unsigned byte lexicographic order. The
count of items less than each item realizes that ordering as a signed rank. -/
def sortStrings (width : Nat) (xs : Array (List UInt8)) : Array (List UInt8) :=
  sortedPayload width (stringKeys xs) xs

def stringSliceSort := sortStrings

theorem sortInts_empty (w : Nat) : sortInts w #[] = #[] := rfl

theorem sortFloat64s_single (w : Nat) (x : BitVec 64) :
    sortFloat64s w #[x] = #[x] := rfl

theorem float64SliceSort_same (w : Nat) (xs : Array (BitVec 64)) :
    float64SliceSort w xs = sortFloat64s w xs := rfl

theorem sortStrings_single (w : Nat) (x : List UInt8) :
    sortStrings w #[x] = #[x] := rfl

theorem stringSliceSort_same (w : Nat) (xs : Array (List UInt8)) :
    stringSliceSort w xs = sortStrings w xs := rfl

end GoLean.Sort
