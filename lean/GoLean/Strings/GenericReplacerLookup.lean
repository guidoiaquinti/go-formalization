import Init.Data.BitVec.Basic

namespace GoLean.Strings

/-- Suffix of a constructor key reachable at the current trie position. -/
private structure LookupCandidate where
  suffix : List (BitVec 8)
  value : List (BitVec 8)
  priority : Nat

structure GenericLookupMatch where
  value : List (BitVec 8)
  keyLength : Nat
  deriving DecidableEq

private def advanceCandidates (byte : BitVec 8) (candidates : List LookupCandidate) :
    List LookupCandidate :=
  candidates.filterMap fun candidate =>
    match candidate.suffix with
    | first :: rest =>
      if first == byte then some {candidate with suffix := rest} else none
    | [] => none

private def visitLookup (candidates : List LookupCandidate) (keyLength : Nat) (skipTerminal : Bool)
    (best : Option (Nat × GenericLookupMatch)) (input : List (BitVec 8)) :
    Option (Nat × GenericLookupMatch) :=
  let best := if skipTerminal then best else
    match candidates.find? (fun candidate => candidate.suffix.isEmpty) with
    | some terminal =>
      if terminal.priority > (best.map Prod.fst).getD 0 then
        some (terminal.priority, ⟨terminal.value, keyLength⟩)
      else best
    | none => best
  match input with
  | [] => best
  | byte :: rest =>
    match advanceCandidates byte candidates with
    | [] => best
    | next => visitLookup next (keyLength + 1) false best rest

/-- Go 1.27.0 src/strings/replace.go:221–250 and 268–298.
    Follow the single reachable byte path and retain the best terminal.
    Constructor priorities decrease in pair order. Candidate suffixes summarize
    compressed edges and byte tables while retaining terminal behavior. -/
def genericReplacerLookup (pairs : List (List (BitVec 8) × List (BitVec 8)))
    (input : List (BitVec 8)) (ignoreRoot : Bool) : Option GenericLookupMatch :=
  let candidates := pairs.zipIdx.map fun ((key, value), index) =>
    ⟨key, value, 2 * (pairs.length - index)⟩
  (visitLookup candidates 0 ignoreRoot none input).map Prod.snd

theorem genericReplacerLookup_noPairs (input : List (BitVec 8)) (ignoreRoot : Bool) :
    genericReplacerLookup [] input ignoreRoot = none := by
  cases input <;> cases ignoreRoot <;> rfl

theorem genericReplacerLookup_firstDuplicate (key value later : List (BitVec 8)) :
    genericReplacerLookup [([], value), ([], later)] key false = some ⟨value, 0⟩ := by
  cases key <;> rfl

theorem genericReplacerLookup_ignoreEmpty (input value : List (BitVec 8)) :
    genericReplacerLookup [([], value)] input true = none := by
  cases input <;> rfl

theorem genericReplacerLookup_shorterPriority :
    genericReplacerLookup [([97], [1]), ([97, 98], [2])] [97, 98] false = some ⟨[1], 1⟩ := rfl

theorem genericReplacerLookup_longerPriority :
    genericReplacerLookup [([97, 98], [1]), ([97], [2])] [97, 98] false = some ⟨[1], 2⟩ := rfl

end GoLean.Strings
