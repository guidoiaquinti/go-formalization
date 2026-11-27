import GoLean.Internal.Dag.HasEdge

namespace GoLean.Internal.Dag

/-- Go looks up the source row and deletes the destination from that row.
    Nil receivers panic. Missing rows and nil maps make delete a no-op.
    Rows model distinct Go maps and association lists represent unique keys.
    https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L51-L55
    https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L75-L77 -/
def delEdge (receiver : Option (Option EdgeMap)) (source to : String) : Option (Option EdgeMap) :=
  receiver.map fun outer =>
    outer.map fun rows =>
      rows.map fun (key, row) =>
        (key, if key == source then row.map (·.filter fun (destination, _) => destination != to)
              else row)

private theorem lookupFiltered (entries : List (String × Bool)) (to query : String) :
    (entries.filter fun (destination, _) => destination != to).lookup query =
      if query = to then none else entries.lookup query := by
  induction entries with
  | nil => simp
  | cons entry entries ih =>
    rcases entry with ⟨key, value⟩
    by_cases hk : key = to
    · subst key
      rw [List.filter_cons_of_neg (by simp), ih, List.lookup_cons]
      by_cases hq : query = to
      · simp [hq]
      · simp [hq, beq_eq_false_iff_ne.mpr hq]
    · rw [List.filter_cons_of_pos (by simpa using hk), List.lookup_cons, List.lookup_cons]
      by_cases hq : query = to
      · subst query
        rw [show (to == key) = false from beq_eq_false_iff_ne.mpr (Ne.symm hk), ih]
      · rw [ite_eq_right hq]
        cases query == key <;> simp [ih, hq]

private theorem lookupMapped (rows : EdgeMap) (source to query : String) :
    (rows.map fun (key, row) =>
      (key, if key == source then row.map (·.filter fun (destination, _) => destination != to)
            else row)).lookup query =
      (rows.lookup query).map (fun row =>
        if query == source then row.map (·.filter fun (destination, _) => destination != to)
        else row) := by
  induction rows with
  | nil => rfl
  | cons entry rows ih =>
    rcases entry with ⟨key, row⟩
    by_cases h : query = key
    · subst key
      simp
    · simp only [List.map_cons, List.lookup_cons, beq_eq_false_iff_ne.mpr h]
      exact ih

/-- Lookup observes panic, outer and inner nilness, and key presence independently.
    Only the destination at the specified source becomes absent. -/
theorem delEdge_spec (receiver : Option (Option EdgeMap)) (source to querySrc queryDst : String) :
    ((delEdge receiver source to).map fun outer =>
      outer.map fun rows => (rows.lookup querySrc).map fun row =>
        row.map (·.lookup queryDst)) =
    (receiver.map fun outer => outer.map fun rows => (rows.lookup querySrc).map fun row =>
      row.map fun entries =>
        if querySrc = source ∧ queryDst = to then none else entries.lookup queryDst) := by
  cases receiver with
  | none => rfl
  | some outer =>
    cases outer with
    | none => rfl
    | some rows =>
      simp only [delEdge, Option.map_some, lookupMapped]
      cases h : rows.lookup querySrc with
      | none => rfl
      | some row =>
        cases row with
        | none => by_cases hs : querySrc = source <;> simp [hs]
        | some entries =>
          by_cases hs : querySrc = source
          · simp [hs, lookupFiltered]
          · simp [hs]

end GoLean.Internal.Dag
