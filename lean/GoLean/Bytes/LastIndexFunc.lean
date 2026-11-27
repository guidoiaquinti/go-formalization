import GoLean.Strings.LastIndexFunc

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:892–894 and 913–925.
    For fixed bytes, the direct ASCII tail path and DecodeLastRune
    give the same width and scalar as the shared reverse scan. -/
def lastIndexFunc (input : List (BitVec 8)) (predicate : Option (Nat → Bool)) :
    Option (Int × List Nat) :=
  GoLean.Strings.lastIndexFunc input predicate

theorem lastIndexFunc_empty (predicate : Option (Nat → Bool)) :
    lastIndexFunc [] predicate = some (-1, []) := rfl

theorem lastIndexFunc_nil (head : BitVec 8) (tail : List (BitVec 8)) :
    lastIndexFunc (head :: tail) none = none := rfl

theorem lastIndexFunc_firstCall (head : BitVec 8) (tail : List (BitVec 8)) (f : Nat → Bool)
    (matched : f (GoLean.Unicode.Utf8.decodeLastRuneInString ((head :: tail).map BitVec.toNat)).1 = true) :
    lastIndexFunc (head :: tail) (some f) =
      let decoded := GoLean.Unicode.Utf8.decodeLastRuneInString ((head :: tail).map BitVec.toNat)
      some ((Int.ofNat ((head :: tail).length - decoded.2)), [decoded.1]) :=
  GoLean.Strings.lastIndexFunc_firstCall head tail f matched

end GoLean.Bytes
