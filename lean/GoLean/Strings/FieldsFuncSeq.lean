import GoLean.Unicode.Utf8.DecodeRuneInString

namespace GoLean.Strings

/-- Observable calls in one invocation of the iterator. -/
inductive FieldEvent where
  | predicate (rune : Nat)
  | yield (field : List Nat)

private def scanFuncFields (input : List Nat) (predicate : Nat → Bool)
    (fuel i : Nat) (start : Option Nat) (budget : Option Nat) : List FieldEvent :=
  match fuel, budget with
  | _, some 0 => []
  | 0, _ => (start.map fun first => FieldEvent.yield (input.drop first)).toList
  | fuel + 1, _ =>
    match input.drop i with
    | [] => (start.map fun first => FieldEvent.yield (input.drop first)).toList
    | lead :: rest =>
      let (rune, width) := GoLean.Unicode.Utf8.decodeNonempty lead rest
      let next := i + width
      FieldEvent.predicate rune ::
        if predicate rune then
          match start with
          | none => scanFuncFields input predicate fuel next none budget
          | some first =>
            FieldEvent.yield ((input.drop first).take (i - first)) ::
              scanFuncFields input predicate fuel next none (budget.map Nat.pred)
        else scanFuncFields input predicate fuel next (some (start.getD i)) budget

/--
Go 1.27.0 strings/iter.go:116–136. Each invocation starts at byte zero.
A finite budget counts the yield call that returns false. Zero represents
skipping the invocation. The predicate is deterministic and total.
-/
def fieldsFuncSeq (input : List (BitVec 8)) (predicate : Nat → Bool)
    (budget : Option Nat) : List FieldEvent :=
  scanFuncFields (input.map BitVec.toNat) predicate input.length 0 none budget

/-- Keep events through the yield call that exhausts the budget. -/
def FieldEvent.stopAfter : Nat → List FieldEvent → List FieldEvent
  | 0, _ => []
  | _ + 1, [] => []
  | limit + 1, predicate r :: rest => predicate r :: stopAfter (limit + 1) rest
  | limit + 1, yield field :: rest => yield field :: stopAfter limit rest

@[simp] private theorem stopAfter_nil (limit : Nat) : FieldEvent.stopAfter limit [] = [] := by
  cases limit <;> rfl

private theorem scanFuncFields_stop (input : List Nat) (predicate : Nat → Bool)
    (fuel i : Nat) (start : Option Nat) (limit : Nat) :
    scanFuncFields input predicate fuel i start (some limit) =
      FieldEvent.stopAfter limit (scanFuncFields input predicate fuel i start none) := by
  induction fuel generalizing i start limit with
  | zero =>
    cases limit with
    | zero => simp [scanFuncFields, FieldEvent.stopAfter]
    | succ limit => cases start <;> simp [scanFuncFields, FieldEvent.stopAfter]
  | succ fuel ih =>
    cases limit with
    | zero => simp [scanFuncFields, FieldEvent.stopAfter]
    | succ limit =>
      simp only [scanFuncFields]
      split
      · cases start <;> simp [FieldEvent.stopAfter]
      · simp only [FieldEvent.stopAfter, List.cons.injEq, true_and]
        split
        · split
          · exact ih _ _ _
          · simp only [FieldEvent.stopAfter, List.cons.injEq, true_and, Option.map_some, Nat.pred_succ]
            exact ih _ _ _
        · exact ih _ _ _

/-- A stopped invocation makes exactly the calls through its final yield. -/
theorem fieldsFuncSeq_stop (input : List (BitVec 8)) (predicate : Nat → Bool) (limit : Nat) :
    fieldsFuncSeq input predicate (some limit) =
      FieldEvent.stopAfter limit (fieldsFuncSeq input predicate none) :=
  scanFuncFields_stop _ _ _ _ _ _

end GoLean.Strings
