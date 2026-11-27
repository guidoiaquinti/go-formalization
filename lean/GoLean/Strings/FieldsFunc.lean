import GoLean.Unicode.Utf8.DecodeRuneInString

namespace GoLean.Strings

/-- Scan byte positions before forming any result substrings. A missing start
    represents the negative sentinel, whose exact value is immaterial. -/
private def collectFieldSpans (predicate : Nat → Bool) (size : Nat)
    (fuel : Nat) (suffix : List Nat) (offset : Nat) (start : Option Nat) :
    List Nat × List (Nat × Nat) :=
  match fuel, suffix with
  | 0, _ | _, [] => ([], (start.map fun first => (first, size)).toList)
  | fuel + 1, lead :: rest =>
    let (rune, width) := GoLean.Unicode.Utf8.decodeNonempty lead rest
    let next := offset + width
    let remaining := rest.drop (width - 1)
    if predicate rune then
      let (calls, spans) := collectFieldSpans predicate size fuel remaining next none
      (rune :: calls, (start.map fun first => (first, offset)).toList ++ spans)
    else
      let (calls, spans) := collectFieldSpans predicate size fuel remaining next
        (some (start.getD offset))
      (rune :: calls, spans)

/-- Go 1.27.0 strings/strings.go:422–465. The first component records the
    pinned predicate calls. The second materializes the collected spans.
    The predicate is a total deterministic function of the rune. -/
def fieldsFunc (input : List (BitVec 8)) (predicate : Nat → Bool) :
    List Nat × List (List Nat) :=
  let bytes := input.map BitVec.toNat
  let (calls, spans) := collectFieldSpans predicate bytes.length bytes.length bytes 0 none
  (calls, spans.map fun (first, last) => (bytes.drop first).take (last - first))

private theorem collectFieldSpans_allSeparators (size fuel offset : Nat) (suffix : List Nat) :
    (collectFieldSpans (fun _ => true) size fuel suffix offset none).2 = [] := by
  induction fuel generalizing suffix offset with
  | zero => simp [collectFieldSpans]
  | succ fuel ih =>
    cases suffix with
    | nil => simp [collectFieldSpans]
    | cons lead rest =>
      simpa [collectFieldSpans] using (ih
        (offset + (GoLean.Unicode.Utf8.decodeNonempty lead rest).2)
        (rest.drop ((GoLean.Unicode.Utf8.decodeNonempty lead rest).2 - 1)))

private theorem collectFieldSpans_keep (size fuel offset first : Nat) (suffix : List Nat) :
    (collectFieldSpans (fun _ => false) size fuel suffix offset (some first)).2 =
      [(first, size)] := by
  induction fuel generalizing suffix offset with
  | zero => simp [collectFieldSpans]
  | succ fuel ih =>
    cases suffix with
    | nil => simp [collectFieldSpans]
    | cons lead rest =>
      simpa [collectFieldSpans] using (ih
        (offset + (GoLean.Unicode.Utf8.decodeNonempty lead rest).2)
        (rest.drop ((GoLean.Unicode.Utf8.decodeNonempty lead rest).2 - 1)))

theorem fieldsFunc_allSeparators (input : List (BitVec 8)) :
    (fieldsFunc input (fun _ => true)).2 = [] := by
  simp [fieldsFunc, collectFieldSpans_allSeparators]

theorem fieldsFunc_noSeparators (lead : BitVec 8) (rest : List (BitVec 8)) :
    (fieldsFunc (lead :: rest) (fun _ => false)).2 = [(lead :: rest).map BitVec.toNat] := by
  simp [fieldsFunc, collectFieldSpans, collectFieldSpans_keep]
  rw [List.take_of_length_le (by simp)]

end GoLean.Strings
