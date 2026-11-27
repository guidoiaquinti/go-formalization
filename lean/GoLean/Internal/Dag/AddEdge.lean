import GoLean.Internal.Dag.HasEdge

namespace GoLean.Internal.Dag

private def setDestination (entries : List (String × Bool)) (to : String) :
    List (String × Bool) :=
  (to, true) :: entries.filter (fun (destination, _) => destination != to)

/-- Go looks up the source row and assigns true to the destination in that row.
    A nil receiver or nil destination map panics.
    https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L51-L55
    https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L71-L73 -/
def addEdge (receiver : Option (Option EdgeMap)) (source to : String) : Option EdgeMap :=
  match receiver with
  | some (some rows) =>
    match (rows.lookup source).join with
    | none => none
    | some _ => some (rows.map fun (key, entries) =>
      (key, if key == source then entries.map (setDestination · to) else entries))
  | _ => none

private theorem lookupOther (entries : List (String × Bool)) (to query : String)
    (h : query ≠ to) :
    (entries.filter fun (key, _) => key != to).lookup query = entries.lookup query := by
  induction entries with
  | nil => rfl
  | cons entry entries ih =>
    rcases entry with ⟨key, value⟩
    by_cases hk : key = to
    · subst key
      simp [List.lookup_cons, beq_eq_false_iff_ne.mpr h, ih]
    · simp [List.lookup_cons, hk, ih]

private theorem lookupSet (entries : List (String × Bool)) (to query : String) :
    (setDestination entries to).lookup query =
      if query = to then some true else entries.lookup query := by
  by_cases h : query = to
  · simp [setDestination, h]
  · simp [setDestination, List.lookup_cons, h, beq_eq_false_iff_ne.mpr h, lookupOther]

private theorem lookupMapped (rows : EdgeMap) (source to query : String) :
    (rows.map fun (key, row) =>
      (key, if key == source then row.map (setDestination · to) else row)).lookup query =
    (rows.lookup query).map (fun row =>
      if query == source then row.map (setDestination · to) else row) := by
  induction rows with
  | nil => rfl
  | cons entry rows ih =>
    rcases entry with ⟨key, row⟩
    by_cases h : query = key
    · subst key
      simp
    · simp only [List.map_cons, List.lookup_cons, beq_eq_false_iff_ne.mpr h]
      exact ih

/-- The result is absent on panic. On success the selected entry is present
    and true, and every other observation of the maps is preserved. -/
theorem addEdge_spec (receiver : Option (Option EdgeMap)) (source to querySrc queryDst : String) :
    (addEdge receiver source to).map (fun rows =>
      (rows.lookup querySrc).map fun row => row.map (·.lookup queryDst)) =
    match receiver with
    | some (some rows) =>
      match rows.lookup source with
      | some (some _) => some ((rows.lookup querySrc).map fun row =>
        row.map fun entries =>
          if querySrc = source ∧ queryDst = to then some true else entries.lookup queryDst)
      | _ => none
    | _ => none := by
  cases receiver with
  | none => rfl
  | some outer =>
    cases outer with
    | none => rfl
    | some rows =>
      cases h : rows.lookup source with
      | none => simp [addEdge, h]
      | some row =>
        cases row with
        | none => simp [addEdge, h]
        | some entries =>
          simp only [addEdge, h, Option.join_some, Option.map_some, lookupMapped]
          by_cases hs : querySrc = source
          · subst querySrc
            simp [h, lookupSet]
          · simp [hs, beq_eq_false_iff_ne.mpr hs]

end GoLean.Internal.Dag
