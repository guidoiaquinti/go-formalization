import GoLean.Strings.LastIndexFunc
import GoLean.Bytes.LastIndexFunc

namespace GoLean.Bytes

/-- Register the shared reverse scan as the model of the private entry
    point. Go 1.27.0 bytes/bytes.go:914–926. -/
abbrev lastIndexFuncPrivate := GoLean.Strings.lastIndexFuncWithTruth

theorem lastIndexFuncPrivate_empty (predicate : Option (Nat → Bool)) (truth : Bool) :
    lastIndexFuncPrivate [] predicate truth = some (-1, []) := rfl

theorem lastIndexFuncPrivate_nil (lead : BitVec 8) (rest : List (BitVec 8)) (truth : Bool) :
    lastIndexFuncPrivate (lead :: rest) none truth = none := rfl

theorem lastIndexFuncPrivate_true (input : List (BitVec 8)) (predicate : Option (Nat → Bool)) :
    lastIndexFuncPrivate input predicate true = lastIndexFunc input predicate := rfl

end GoLean.Bytes
