import GoLean.Strings.SplitSeq

namespace GoLean.Bytes

/-- Go 1.27.0 src/bytes/iter.go:36–69. With fixed backing bytes, the
    fragments and captured suffix follow the string scanner. Each yielded
    slice uses a full slice expression to cap capacity at its length. -/
def splitSeq (input separator : List Nat) (budget : Nat) :
    List (List Nat × Nat) × List Nat :=
  let (fragments, remaining) := GoLean.Strings.splitSeq input separator budget
  (fragments.map (fun fragment => (fragment, fragment.length)), remaining)

theorem splitSeq_caps (input separator : List Nat) (budget : Nat) :
    ∀ fragment ∈ (splitSeq input separator budget).1, fragment.2 = fragment.1.length := by
  intro fragment member
  simp only [splitSeq, List.mem_map] at member
  obtain ⟨bytes, _, rfl⟩ := member
  rfl

theorem splitSeq_stop_first (input separator : List Nat) :
    (splitSeq input separator 1).2 = input := by
  exact GoLean.Strings.splitSeq_stop_first input separator

end GoLean.Bytes
