import GoLean.Strings.SplitSeq

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/iter.go:36–77. Retain the full separator and
    cap each yielded subslice at its length. The backing bytes are fixed. -/
def splitAfterSeq (input separator : List Nat) (budget : Nat) :
    List (List Nat × Nat) × List Nat :=
  let (fragments, remaining) := GoLean.Strings.splitAfterSeq input separator budget
  (fragments.map (fun fragment => (fragment, fragment.length)), remaining)

theorem splitAfterSeq_caps (input separator : List Nat) (budget : Nat) :
    ∀ fragment ∈ (splitAfterSeq input separator budget).1, fragment.2 = fragment.1.length := by
  intro fragment member
  simp only [splitAfterSeq, List.mem_map] at member
  obtain ⟨bytes, _, rfl⟩ := member
  rfl

theorem splitAfterSeq_first (input separator : List Nat) (i : Nat)
    (nonempty : separator ≠ []) (found : GoLean.Strings.index input separator = some i) :
    splitAfterSeq input separator 1 =
      ([(input.take i ++ separator, (input.take i ++ separator).length)], input) := by
  simp [splitAfterSeq, GoLean.Strings.splitAfterSeq_first input separator i nonempty found]

end GoLean.Bytes
