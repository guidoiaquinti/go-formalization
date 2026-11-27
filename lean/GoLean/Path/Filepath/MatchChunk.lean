import GoLean.Path.Filepath.GetEsc

namespace GoLean.Path.Filepath

private inductive Atom where
  | literal (byte : Nat)
  | anyRune
  | characterClass (negated : Bool) (intervals : List (Nat × Nat))

private def parseClass (windows : Bool) :
    Nat → List (Nat × Nat) → List Nat → Option (List (Nat × Nat) × List Nat)
  | 0, _, _ => none
  | fuel + 1, acc, input =>
    if input.head? == some 93 && !acc.isEmpty then some (acc, input.drop 1)
    else
      let (lo, tail, bad) := getEsc windows input
      if bad then none
      else if tail.head? == some 45 then
        let (hi, rest, bad) := getEsc windows (tail.drop 1)
        if bad then none else parseClass windows fuel ((lo, hi) :: acc) rest
      else parseClass windows fuel ((lo, lo) :: acc) tail

private def parseChunk (windows : Bool) : Nat → List Nat → Option (List Atom)
  | 0, _ => none
  | _ + 1, [] => some []
  | fuel + 1, byte :: tail => do
    if byte = 91 then
      let negated := tail.head? == some 94
      let body := if negated then tail.drop 1 else tail
      let (intervals, rest) ← parseClass windows (body.length + 1) [] body
      return Atom.characterClass negated intervals :: (← parseChunk windows fuel rest)
    else if byte = 63 then
      return Atom.anyRune :: (← parseChunk windows fuel tail)
    else if byte = 92 && !windows then
      match tail with
      | [] => none
      | literal :: rest => return Atom.literal literal :: (← parseChunk windows fuel rest)
    else return Atom.literal byte :: (← parseChunk windows fuel tail)

private def consume (windows : Bool) (atom : Atom) (input : List Nat) : Option (List Nat) :=
  match input with
  | [] => none
  | byte :: tail =>
    match atom with
    | .literal lit => if byte = lit then some tail else none
    | .anyRune =>
      if byte = (if windows then 92 else 47) then none
      else some (input.drop (GoLean.Unicode.Utf8.decodeNonempty byte tail).2)
    | .characterClass negated intervals =>
      let (rune, width) := GoLean.Unicode.Utf8.decodeNonempty byte tail
      let member := intervals.any fun (lo, hi) => lo ≤ rune && rune ≤ hi
      if member == negated then none else some (input.drop width)

def matchChunk (windows : Bool) (chunk input : List Nat) : List Nat × Bool × Bool :=
  match parseChunk windows (chunk.length + 1) chunk with
  | none => ([], false, true)
  | some atoms =>
    match atoms.foldlM (fun rest atom => consume windows atom rest) input with
    | none => ([], false, false)
    | some rest => (rest, true, false)

theorem matchChunk_empty (windows : Bool) (input : List Nat) :
    matchChunk windows [] input = (input, true, false) := rfl

theorem matchChunk_invalidAfterMismatch :
    matchChunk false [97, 91] [98] = ([], false, true) := by decide

theorem matchChunk_questionSeparator :
    matchChunk false [63] [47] = ([], false, false) ∧
    matchChunk true [63] [47] = ([], true, false) := by decide

theorem matchChunk_classRange :
    matchChunk false [91, 97, 45, 99, 93] [98, 120] = ([120], true, false) := by decide

end GoLean.Path.Filepath
