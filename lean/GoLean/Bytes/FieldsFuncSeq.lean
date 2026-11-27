import GoLean.Strings.FieldsFuncSeq

namespace GoLean.Bytes

inductive FieldEvent where
  | predicate (rune : Nat)
  | yield (bytes : List Nat) (capacity : Nat)

private def withCapacity : GoLean.Strings.FieldEvent → FieldEvent
  | .predicate r => .predicate r
  | .yield field => .yield field field.length

/-- Go 1.27.0 src/bytes/iter.go:116–136. With fixed backing bytes, the
    predicate and yield trace follows the string scanner. Each yielded
    fragment receives capacity equal to its length. -/
def fieldsFuncSeq (input : List (BitVec 8)) (predicate : Nat → Bool)
    (budget : Option Nat) : List FieldEvent :=
  (GoLean.Strings.fieldsFuncSeq input predicate budget).map withCapacity

def FieldEvent.stopAfter : Nat → List FieldEvent → List FieldEvent
  | 0, _ => []
  | _ + 1, [] => []
  | n + 1, .predicate r :: tail => .predicate r :: stopAfter (n + 1) tail
  | n + 1, .yield bytes cap :: tail => .yield bytes cap :: stopAfter n tail

private theorem stopAfter_map (events : List GoLean.Strings.FieldEvent) (n : Nat) :
    (GoLean.Strings.FieldEvent.stopAfter n events).map withCapacity =
      FieldEvent.stopAfter n (events.map withCapacity) := by
  induction events generalizing n with
  | nil => cases n <;> rfl
  | cons head tail ih =>
    cases n with
    | zero => rfl
    | succ n =>
      cases head <;> simp [GoLean.Strings.FieldEvent.stopAfter, FieldEvent.stopAfter, withCapacity, ih]

theorem fieldsFuncSeq_stop (input : List (BitVec 8)) (predicate : Nat → Bool) (n : Nat) :
    fieldsFuncSeq input predicate (some n) =
      FieldEvent.stopAfter n (fieldsFuncSeq input predicate none) := by
  simp [fieldsFuncSeq, GoLean.Strings.fieldsFuncSeq_stop, stopAfter_map]

theorem fieldsFuncSeq_caps (input : List (BitVec 8)) (predicate : Nat → Bool)
    (budget : Option Nat) (field : List Nat) (capacity : Nat)
    (member : FieldEvent.yield field capacity ∈ fieldsFuncSeq input predicate budget) :
    capacity = field.length := by
  simp only [fieldsFuncSeq, List.mem_map] at member
  obtain ⟨event, _, mapped⟩ := member
  cases event <;> simp [withCapacity] at mapped
  rcases mapped with ⟨rfl, rfl⟩
  rfl

end GoLean.Bytes
