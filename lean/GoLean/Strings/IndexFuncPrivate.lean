import GoLean.Strings.IndexFunc

namespace GoLean.Strings

/-- Register the shared forward scan as the model of the private entry
    point. Go 1.27.0 strings/strings.go:910–917. -/
abbrev indexFuncPrivate := indexFuncWithTruth

theorem indexFuncPrivate_empty (predicate : Option (Nat → Bool)) (truth : Bool) :
    indexFuncPrivate [] predicate truth = some (-1, []) := rfl

theorem indexFuncPrivate_nil (lead : BitVec 8) (rest : List (BitVec 8)) (truth : Bool) :
    indexFuncPrivate (lead :: rest) none truth = none := rfl

theorem indexFuncPrivate_true (input : List (BitVec 8)) (predicate : Option (Nat → Bool)) :
    indexFuncPrivate input predicate true = indexFunc input predicate := rfl

end GoLean.Strings
