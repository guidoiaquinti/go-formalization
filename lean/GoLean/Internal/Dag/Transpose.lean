namespace GoLean.Internal.Dag

/-- Go rebuilds one row per node and writes a true entry for each reversed key.
Grouping the writes by destination models the independent rows. Old true and
false values both contribute keys. Every old destination must be a node.
https://github.com/golang/go/blob/go1.27.0/src/internal/dag/alg.go#L7-L21 -/
def transpose (nodes : List Nat) (old : List (Nat × Nat × Bool)) :
    List (Nat × List (Nat × Bool)) :=
  nodes.map fun to =>
    (to, (old.filter fun edge => edge.2.1 == to).map fun edge => (edge.1, true))

theorem transpose_rows (nodes : List Nat) (old : List (Nat × Nat × Bool)) :
    (transpose nodes old).map Prod.fst = nodes := by
  simp [transpose, List.map_map, Function.comp_def]

theorem transpose_spec (nodes : List Nat) (old : List (Nat × Nat × Bool))
    (to src : Nat) (value : Bool) :
    (∃ row, (to, row) ∈ transpose nodes old ∧ (src, value) ∈ row) ↔
      to ∈ nodes ∧ value = true ∧ ∃ prior, (src, to, prior) ∈ old := by
  constructor
  · rintro ⟨row, present, entry⟩
    rcases List.mem_map.mp present with ⟨key, hk, eq⟩
    cases eq
    rcases List.mem_map.mp entry with ⟨⟨src', dest, prior⟩, hem, heq⟩
    simp only [List.mem_filter, beq_iff_eq] at hem
    simp only [Prod.mk.injEq] at heq
    obtain ⟨mem, rfl⟩ := hem
    obtain ⟨rfl, rfl⟩ := heq
    exact ⟨hk, rfl, prior, mem⟩
  · rintro ⟨hn, rfl, prior, hem⟩
    refine ⟨_, List.mem_map.mpr ⟨to, hn, rfl⟩, List.mem_map.mpr ⟨(src, to, prior), ?_, rfl⟩⟩
    simp [hem]

end GoLean.Internal.Dag
