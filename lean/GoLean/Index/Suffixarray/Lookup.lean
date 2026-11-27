import GoLean.Index.Suffixarray.Accessors
import GoLean.Sort.Search
import Init.Data.List.Lex

namespace GoLean.Index.Suffixarray

private def suffixAt (s : IndexState) (i : Nat) : List (BitVec 8) :=
  match indexAt (some s) (Int.ofNat i) with
  | .ok v => (v.backing.getD []).take v.length
  | .error _ => []

private def rangeStart (s : IndexState) (needle : List (BitVec 8)) (length : Nat) : Nat :=
  GoLean.Sort.search (Int.ofNat length)
    (fun i => !(suffixAt s i).lex needle (fun a b => decide (a.toNat < b.toNat)))

private def rangeCount (s : IndexState) (needle : List (BitVec 8)) (lo length : Nat) : Nat :=
  GoLean.Sort.search (Int.ofNat (length - lo))
    (fun j => (suffixAt s (lo + j)).take needle.length != needle)

/-- Go 1.27.0 src/index/suffixarray/suffixarray.go:251–258. The model assumes a
    complete suffix index in lexicographic order. The two binary searches
    locate the first matching suffix and the first suffix beyond the group. -/
def lookupAll (receiver : Option IndexState) (needle : List (BitVec 8)) :
    Except AccessFailure Ints := do
  let .some s := receiver | throw .nilReceiver
  let length ← intsLen (some s.sa)
  let lo := rangeStart s needle length
  let hi := lo + rangeCount s needle lo length
  intsSlice (some s.sa) (Int.ofNat lo) (Int.ofNat hi)

/-- Go 1.27.0 src/index/suffixarray/suffixarray.go:265–285. Guards precede
    receiver access. A successful nonempty output is a fresh native-int list. -/
def indexLookup (width : Nat) (receiver : Option IndexState) (needle : List (BitVec 8))
    (limit : Int) : Except AccessFailure (Option (List Int)) := do
  if needle.isEmpty || limit = 0 then return none
  let located ← lookupAll receiver needle
  let length ← intsLen (some located)
  let count := if limit < 0 then length else min length limit.toNat
  if count = 0 then return none
  let found ← (List.range count).mapM fun i => intsGet (some located) (Int.ofNat i)
  return some (found.map fun i => (BitVec.ofInt width i).toInt)

theorem indexLookup_empty (width : Nat) (receiver : Option IndexState) (limit : Int) :
    indexLookup width receiver [] limit = .ok none := by simp [indexLookup]; rfl

theorem indexLookup_zero (width : Nat) (receiver : Option IndexState)
    (needle : List (BitVec 8)) :
    indexLookup width receiver needle 0 = .ok none := by simp [indexLookup]; rfl

private def bananaIndex : IndexState :=
  ⟨⟨some [98, 97, 110, 97, 110, 97], 6⟩,
    ⟨⟨some [5, 3, 1, 0, 4, 2, 99], 6⟩, ⟨none, 0⟩⟩⟩

private theorem lookupAll_bananaExact :
    lookupAll (some bananaIndex) [97, 110, 97] =
      .ok ⟨⟨some [3, 1, 0, 4, 2, 99], 2⟩, ⟨none, 0⟩⟩ := by
  have first : rangeStart bananaIndex [97, 110, 97] 6 = 1 := by
    unfold rangeStart
    apply GoLean.Sort.search_spec
    · decide
    · intro i h
      change i < 6 at h
      have cases : i = 0 ∨ i = 1 ∨ i = 2 ∨ i = 3 ∨ i = 4 ∨ i = 5 := by omega
      rcases cases with rfl | rfl | rfl | rfl | rfl | rfl <;> decide
  have second : rangeCount bananaIndex [97, 110, 97] 1 6 = 2 := by
    unfold rangeCount
    apply GoLean.Sort.search_spec
    · decide
    · intro i h
      change i < 5 at h
      have cases : i = 0 ∨ i = 1 ∨ i = 2 ∨ i = 3 ∨ i = 4 := by omega
      rcases cases with rfl | rfl | rfl | rfl | rfl <;> decide
  let firstSlot := rangeStart bananaIndex [97, 110, 97] 6
  change intsSlice (some bananaIndex.sa) (Int.ofNat firstSlot)
    (Int.ofNat (firstSlot + rangeCount bananaIndex [97, 110, 97] firstSlot 6)) = _
  have slot : firstSlot = 1 := first
  rw [slot, second]
  rfl

theorem lookupAll_banana :
    ((lookupAll (some bananaIndex) [97, 110, 97]).toOption.map
      (fun located => (located.words32.length, located.words32.backing))) =
      some (2, some [3, 1, 0, 4, 2, 99]) := by
  rw [lookupAll_bananaExact]
  rfl

theorem indexLookup_bananaLimited :
    (indexLookup 64 (some bananaIndex) [97, 110, 97] 1).toOption = some (some [3]) := by
  simp only [indexLookup, lookupAll_bananaExact]
  rfl

end GoLean.Index.Suffixarray
