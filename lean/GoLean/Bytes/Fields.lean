import GoLean.Strings.Fields

namespace GoLean.Bytes

/-- Go 1.27.0 bytes/bytes.go:447–489. The native ASCII prescan and
    delegated Unicode scanner have the same byte boundaries as strings.Fields.
    Both paths produce a nonnil result with capacity-capped fields. -/
def fields (input : List (BitVec 8)) : Option (List (List Nat × Nat)) :=
  some ((GoLean.Strings.fields input).map fun field => (field, field.length))

theorem fields_empty : fields [] = some [] := rfl

theorem fields_nonnil (input : List (BitVec 8)) : (fields input).isSome = true := rfl

theorem fields_caps (input : List (BitVec 8)) :
    ∀ field ∈ (fields input).getD [], field.2 = field.1.length := by
  intro field member
  simp only [fields, Option.getD_some, List.mem_map] at member
  obtain ⟨_, _, rfl⟩ := member
  rfl

end GoLean.Bytes
