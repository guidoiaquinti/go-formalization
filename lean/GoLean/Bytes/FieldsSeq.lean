import GoLean.Strings.FieldsSeq
import Init.Data.List.TakeDrop

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/iter.go:83–109. The same whitespace scan runs
    for each call on fixed backing bytes. Full slice expressions cap fields. -/
def fieldsSeq (input : List (BitVec 8)) (budget : Option Nat) : List (List Nat × Nat) :=
  (GoLean.Strings.fieldsSeq input budget).map fun field => (field, field.length)

theorem fieldsSeq_take (input : List (BitVec 8)) (limit : Nat) :
    fieldsSeq input (some limit) = (fieldsSeq input none).take limit := by
  simp [fieldsSeq, GoLean.Strings.fieldsSeq_take, List.map_take]

theorem fieldsSeq_caps (input : List (BitVec 8)) (budget : Option Nat) :
    ∀ field ∈ fieldsSeq input budget, field.2 = field.1.length := by
  intro field member
  simp only [fieldsSeq, List.mem_map] at member
  obtain ⟨_, _, rfl⟩ := member
  rfl

end GoLean.Bytes
