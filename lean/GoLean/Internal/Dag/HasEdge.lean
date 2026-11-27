namespace GoLean.Internal.Dag

/-- An edge map with a nullable map of destinations for each source. -/
abbrev EdgeMap := List (String × Option (List (String × Bool)))

/-- Go reads the outer map, checks that the inner map is non-nil, then
    reads the destination with the Boolean zero default. A nil receiver
    panics, represented by a missing result.
    https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L51-L55
    https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L79-L81 -/
def hasEdge (receiver : Option (Option EdgeMap)) (source to : String) : Option Bool := do
  let outer ← receiver
  let inner := (outer.bind fun rows => rows.lookup source).join
  return inner.isSome && ((inner.bind fun row => row.lookup to).getD false)

theorem hasEdge_spec (receiver : Option (Option EdgeMap)) (source to : String) :
    hasEdge receiver source to =
      match receiver with
      | none => none
      | some none => some false
      | some (some rows) =>
        some (match rows.lookup source with
          | some (some row) => row.lookup to == some true
          | _ => false) := by
  cases receiver with
  | none => rfl
  | some outer =>
    cases outer with
    | none => rfl
    | some rows =>
      cases h : rows.lookup source with
      | none => simp [hasEdge, h]
      | some row =>
        cases row with
        | none => simp [hasEdge, h]
        | some entries =>
          cases k : entries.lookup to with
          | none => simp [hasEdge, h, k]
          | some value => cases value <;> simp [hasEdge, h, k]

end GoLean.Internal.Dag
