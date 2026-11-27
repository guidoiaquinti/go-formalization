namespace GoLean.Strings

private def fillASCIISet (table : BitVec 8 → Bool) :
    List (BitVec 8) → (BitVec 8 → Bool) × Bool
  | [] => (table, true)
  | b :: rest =>
    if b.toNat ≥ 128 then (table, false)
    else fillASCIISet (fun c => c == b || table c) rest

/-- Go 1.27.0 strings/strings.go:952–961. The function represents the
    256 table entries. Non-ASCII stops insertion and returns the table so far. -/
def makeASCIISet (input : List (BitVec 8)) : (BitVec 8 → Bool) × Bool :=
  fillASCIISet (fun _ => false) input

theorem makeASCIISet_empty : makeASCIISet [] = (fun _ => false, true) := rfl

private theorem fillASCIISet_stop (table : BitVec 8 → Bool)
    (front suffix : List (BitVec 8)) (b : BitVec 8)
    (ascii : ∀ c ∈ front, c.toNat < 128) (high : 128 ≤ b.toNat) :
    fillASCIISet table (front ++ b :: suffix) = ((fillASCIISet table front).1, false) := by
  induction front generalizing table with
  | nil => simp [fillASCIISet, high]
  | cons c rest ih =>
    have low : ¬128 ≤ c.toNat := Nat.not_le.mpr (ascii c (by simp))
    simp only [List.cons_append, fillASCIISet, low, ite_false]
    exact ih _ (fun d hd => ascii d (by simp [hd]))

theorem makeASCIISet_stop (front suffix : List (BitVec 8)) (b : BitVec 8)
    (ascii : ∀ c ∈ front, c.toNat < 128) (high : 128 ≤ b.toNat) :
    makeASCIISet (front ++ b :: suffix) = ((makeASCIISet front).1, false) :=
  fillASCIISet_stop _ _ _ _ ascii high

private theorem fillASCIISet_high (table : BitVec 8 → Bool)
    (input : List (BitVec 8)) (b : BitVec 8) (high : 128 ≤ b.toNat) :
    (fillASCIISet table input).1 b = table b := by
  induction input generalizing table with
  | nil => rfl
  | cons c rest ih =>
    simp only [fillASCIISet]
    split
    · rfl
    · rename_i low
      have distinct : b ≠ c := by
        intro eq
        exact low (eq ▸ high)
      simp [ih, distinct]

theorem makeASCIISet_high (input : List (BitVec 8)) (b : BitVec 8)
    (high : 128 ≤ b.toNat) : (makeASCIISet input).1 b = false :=
  fillASCIISet_high _ input b high

end GoLean.Strings
