import GoLean.Strings.FieldsFunc

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/bytes.go:506–544. The shared scanner collects
    spans before constructing capacity-capped slices from fixed backing bytes. -/
def fieldsFunc (input : List (BitVec 8)) (predicate : Nat → Bool) :
    List Nat × List (List Nat × Nat) :=
  let (calls, fields) := GoLean.Strings.fieldsFunc input predicate
  (calls, fields.map fun field => (field, field.length))

theorem fieldsFunc_allSeparators (input : List (BitVec 8)) :
    (fieldsFunc input (fun _ => true)).2 = [] := by
  simp [fieldsFunc, GoLean.Strings.fieldsFunc_allSeparators]

theorem fieldsFunc_noSeparators (lead : BitVec 8) (rest : List (BitVec 8)) :
    (fieldsFunc (lead :: rest) (fun _ => false)).2 =
      [((lead :: rest).map BitVec.toNat, (lead :: rest).length)] := by
  simp [fieldsFunc, GoLean.Strings.fieldsFunc_noSeparators]

theorem fieldsFunc_caps (input : List (BitVec 8)) (predicate : Nat → Bool) :
    ∀ field ∈ (fieldsFunc input predicate).2, field.2 = field.1.length := by
  intro field member
  simp only [fieldsFunc, List.mem_map] at member
  obtain ⟨_, _, rfl⟩ := member
  rfl

end GoLean.Bytes
