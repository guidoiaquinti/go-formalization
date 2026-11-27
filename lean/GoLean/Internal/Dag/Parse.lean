import GoLean.Internal.Dag.ParseRules

namespace GoLean.Internal.Dag

abbrev Label := List UInt8
abbrev Arc := Label × Label

inductive ParseOutcome where
  | graph (nodes : List Label) (edges : List Arc)
  | syntax (reason : SyntaxErrorOutcome)
  | semantic (message : List UInt8)
  | exhausted

private def text (s : String) : List UInt8 := s.toUTF8.data.toList
private def sentinel : Label := text "NONE"

private structure ParsedGraph where
  nodes : List Label := []
  arcs : List Arc := []
  errors : List (List UInt8) := []
  denied : List Rule := []

private def reportError (g : ParsedGraph) (msg : List UInt8) : ParsedGraph :=
  { g with errors := g.errors ++ [msg] }

private def insertArc (g : ParsedGraph) (source target : Label) : ParsedGraph :=
  if g.arcs.contains (source, target) then g else { g with arcs := (source, target) :: g.arcs }

private def addDefinition (g : ParsedGraph) (less : List Label) (name : Label) : ParsedGraph :=
  if name == sentinel then
    reportError g (text "NONE cannot be a predecessor")
  else
    let g := if g.nodes.contains name then reportError g (text "multiple definitions for " ++ name)
      else { g with nodes := g.nodes ++ [name] }
    less.foldl (fun g target =>
      if target == sentinel then g
      else if g.nodes.contains target then insertArc g name target
      else reportError g (text "use of " ++ target ++ text " before its definition")) g

private def addRule (g : ParsedGraph) (rule : Rule) : ParsedGraph :=
  if rule.op == [33, 60] then {g with denied := g.denied ++ [rule]}
  else rule.defn.foldl (fun g name => addDefinition g rule.less name) g

private def closeGraph (g : ParsedGraph) : ParsedGraph :=
  let nodes := g.nodes
  nodes.foldl (fun g k =>
    nodes.foldl (fun g i =>
      nodes.foldl (fun g j =>
        if i != k && k != j && g.arcs.contains (i, k) && g.arcs.contains (k, j) then
          let g := if i == j then
            reportError g (text "graph cycle: " ++ j ++ text " < " ++ k ++ text " < " ++ i)
            else g
          insertArc g i j
        else g) g) g) g

private def checkAssertions (g : ParsedGraph) : ParsedGraph :=
  g.denied.foldl (fun g rule =>
    rule.less.foldl (fun g less =>
      rule.defn.foldl (fun g name =>
        if g.arcs.contains (name, less) then
          reportError g (text "graph edge assertion failed: " ++ less ++ text " !< " ++ name)
        else g) g) g) g

/-- Build rows in definition order, run the same in-place closure traversal,
    then check negative assertions. Every inserted destination has a row,
    so the missing-definition scan cannot add an error on constructed graphs.
    https://github.com/golang/go/blob/go1.27.0/src/internal/dag/parse.go#L94-L179 -/
def parse (width : Nat) (input : List UInt8) : ParseOutcome :=
  match parseRules width input with
  | .failed _ err => .syntax err
  | .exhausted => .exhausted
  | .returned rules =>
    let graph := checkAssertions (closeGraph (rules.foldl addRule {}))
    if graph.errors.isEmpty then .graph graph.nodes graph.arcs
    else .semantic (List.intercalate [10] graph.errors)

theorem parse_syntax (width : Nat) (input : List UInt8) (partialRules : List Rule)
    (reason : SyntaxErrorOutcome) (h : parseRules width input = .failed partialRules reason) :
    parse width input = .syntax reason := by
  simp [parse, h]

theorem parse_empty (width : Nat) (input : List UInt8)
    (h : parseRules width input = .returned []) :
    parse width input = .graph [] [] := by
  simp [parse, h, closeGraph, checkAssertions]

theorem parse_single_definition (width : Nat) (input : List UInt8) (label : Label)
    (ordinary : label ≠ sentinel)
    (h : parseRules width input = .returned [⟨[sentinel], [60], [label]⟩]) :
    parse width input = .graph [label] [] := by
  simp [parse, h, addRule, addDefinition, ordinary, closeGraph, checkAssertions]

end GoLean.Internal.Dag
