namespace GoLean.Regexp.Syntax

abbrev RuneRange := Nat × Nat

structure ClassStore where
  visible : List RuneRange
  backing : List RuneRange
  shared : Bool

def rangesLen (pairs : List RuneRange) : Nat := pairs.length

def rangesLess (pairs : List RuneRange) (i j : Nat) : Bool :=
  let (a, b) := (pairs[i]!, pairs[j]!)
  a.1 < b.1 || (a.1 == b.1 && a.2 > b.2)

def rangesSwap (pairs : List RuneRange) (i j : Nat) : List RuneRange :=
  (pairs.set i pairs[j]!).set j pairs[i]!

private def overwrite (t : ClassStore) (i : Nat) (r : RuneRange) : ClassStore :=
  ⟨t.visible.set i r, if t.shared then t.backing.set i r else t.backing, t.shared⟩

private def appendPair (t : ClassStore) (r : RuneRange) : ClassStore :=
  let shared := t.shared && t.visible.length < t.backing.length
  ⟨t.visible ++ [r], if shared then t.backing.set t.visible.length r else t.backing, shared⟩

/-- Go 1.27.0 parse.go:1987. Only the two latest ranges are candidates. -/
def appendRange (t : ClassStore) (r : RuneRange) : ClassStore :=
  let n := t.visible.length
  let merge i := do
    let p ← t.visible[i]?
    if r.1 ≤ p.2 + 1 && p.1 ≤ r.2 + 1 then
      some (overwrite t i (min p.1 r.1, max p.2 r.2))
    else none
  let tail := if n ≥ 1 then merge (n-1) else none
  let prev := if n ≥ 2 then merge (n-2) else none
  (tail.orElse (fun _ => prev)).getD (appendPair t r)

def appendClass (t : ClassStore) (clean : List RuneRange) : ClassStore :=
  clean.foldl appendRange t

private def maxRune : Nat := 0x10ffff

private def initialGaps (clean : List RuneRange) : Nat × List RuneRange :=
  clean.foldl (fun (start, gaps) (lo, hi) =>
    (hi + 1, if start < lo then gaps ++ [(start, lo - 1)] else gaps)) (0, [])

def appendNegatedClass (t : ClassStore) (clean : List RuneRange) : ClassStore :=
  let (start, gaps) := initialGaps clean
  let result := appendClass t gaps
  if start ≤ maxRune then appendRange result (start, maxRune) else result

def negateClass (t : ClassStore) : ClassStore :=
  let (start, gaps) := initialGaps t.visible
  let backing := if t.shared then gaps ++ t.backing.drop gaps.length else t.backing
  let narrowed : ClassStore := ⟨gaps, backing, t.shared⟩
  if start ≤ maxRune then appendPair narrowed (start, maxRune) else narrowed

/-- Clean ranges are sorted and separated. The model summarizes sort.Find as membership. -/
def inCharClass (rune : Nat) (clean : List RuneRange) : Bool :=
  clean.any fun r => r.1 ≤ rune && rune ≤ r.2

theorem rangesLen_empty : rangesLen [] = 0 := rfl

theorem rangesLess_highFirst : rangesLess [(7, 11), (7, 9)] 0 1 = true := by decide

theorem rangesSwap_first : rangesSwap [(1, 2), (3, 4)] 0 1 = [(3, 4), (1, 2)] := rfl

theorem appendRange_checksSecond :
    (appendRange ⟨[(1, 2), (9, 10)], [(1, 2), (9, 10)], true⟩ (3, 4)).visible =
      [(1, 4), (9, 10)] := by decide

theorem appendClass_empty (t : ClassStore) : appendClass t [] = t := rfl

theorem appendNegatedClass_full (t : ClassStore) :
    appendNegatedClass t [(0, 0x10ffff)] = t := rfl

theorem negateClass_empty :
    (negateClass ⟨[], [], true⟩).visible = [(0, 0x10ffff)] := rfl

theorem inCharClass_inclusive :
    inCharClass 1 [(1, 3)] = true ∧ inCharClass 3 [(1, 3)] = true := by decide

end GoLean.Regexp.Syntax
