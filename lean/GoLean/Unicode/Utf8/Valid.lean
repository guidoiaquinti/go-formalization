import GoLean.Strings.Count

namespace GoLean.Unicode.Utf8

/-- The pinned word reads pack bytes into disjoint lanes. Testing their
    high-bit mask is equivalent to checking every byte for ASCII. -/
def asciiBlock (width : Nat) (bytes : List Nat) : Bool :=
  width < bytes.length && (bytes.take width).all (· < 128)

private def skipLargeBlocks (width : Nat) : Nat → List Nat → List Nat
  | 0, bytes => bytes
  | fuel + 1, bytes =>
    if asciiBlock width bytes then skipLargeBlocks width fuel (bytes.drop width) else bytes

private def skipASCII (wordWidth : Nat) (bytes : List Nat) : List Nat :=
  if asciiBlock wordWidth bytes then
    let rest := bytes.drop wordWidth
    if asciiBlock (2 * wordWidth) rest then
      let rest := rest.drop (2 * wordWidth)
      skipLargeBlocks (4 * wordWidth) rest.length rest
    else rest
  else bytes

private def validate (wordWidth : Nat) : Nat → List Nat → Bool
  | _, [] => true
  | 0, _ :: _ => false
  | fuel + 1, lead :: rest =>
    if lead < 128 then validate wordWidth fuel (skipASCII wordWidth rest)
    else
      let width := GoLean.Strings.runeWidth lead rest
      if width = 1 then false else validate wordWidth fuel (rest.drop (width - 1))

private theorem skipLargeBlocks_sublist (width fuel : Nat) (bytes : List Nat) :
    (skipLargeBlocks width fuel bytes).Sublist bytes := by
  induction fuel generalizing bytes with
  | zero => exact List.Sublist.refl _
  | succ fuel ih =>
    simp only [skipLargeBlocks]
    split
    · exact (ih _).trans (List.drop_sublist _ _)
    · exact List.Sublist.refl _

private theorem skipASCII_sublist (width : Nat) (bytes : List Nat) :
    (skipASCII width bytes).Sublist bytes := by
  unfold skipASCII
  split
  · dsimp only
    split
    · exact (skipLargeBlocks_sublist _ _ _).trans
        ((List.drop_sublist _ _).trans (List.drop_sublist _ _))
    · exact List.drop_sublist _ _
  · exact List.Sublist.refl _

private theorem validate_ascii_only (width fuel : Nat) (bytes : List Nat)
    (enough : bytes.length ≤ fuel) (ascii : ∀ b ∈ bytes, b < 128) :
    validate width fuel bytes = true := by
  induction fuel generalizing bytes with
  | zero =>
    have empty : bytes = [] := List.length_eq_zero_iff.mp (by omega)
    simp [empty, validate]
  | succ fuel ih =>
    cases bytes with
    | nil => rfl
    | cons b bs =>
      have head := ascii b (by simp)
      simp only [validate, head, ↓reduceIte]
      have sub := skipASCII_sublist width bs
      apply ih
      · have shorter := sub.length_le
        simp only [List.length_cons] at enough
        omega
      · intro c mem
        exact ascii c (List.mem_cons_of_mem b (sub.subset mem))

/-- Go 1.27.0 src/unicode/utf8/utf8.go:460–519. Word width is in bytes. -/
def valid (wordWidth : Nat) (bytes : List Nat) : Bool :=
  validate wordWidth bytes.length bytes

theorem valid_ascii_only (wordWidth : Nat) (bytes : List Nat)
    (ascii : ∀ b ∈ bytes, b < 128) : valid wordWidth bytes = true :=
  validate_ascii_only wordWidth bytes.length bytes (Nat.le_refl _) ascii

theorem valid_empty (wordWidth : Nat) : valid wordWidth [] = true := rfl

theorem valid_rejects_malformed_lead (wordWidth lead : Nat) (rest : List Nat)
    (nonASCII : 128 ≤ lead) (malformed : GoLean.Strings.runeWidth lead rest = 1) :
    valid wordWidth (lead :: rest) = false := by
  simp [valid, validate, Nat.not_lt_of_ge nonASCII, malformed]

theorem valid_multibyte_step (wordWidth lead : Nat) (rest : List Nat)
    (nonASCII : 128 ≤ lead) (legal : GoLean.Strings.runeWidth lead rest ≠ 1) :
    valid wordWidth (lead :: rest) =
      validate wordWidth rest.length (rest.drop (GoLean.Strings.runeWidth lead rest - 1)) := by
  simp [valid, validate, Nat.not_lt_of_ge nonASCII, legal]

end GoLean.Unicode.Utf8
