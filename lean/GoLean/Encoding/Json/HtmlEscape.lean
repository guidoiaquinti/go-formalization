namespace GoLean.Encoding.Json

private def escapedByte : Nat → Option (List Nat)
  | 60 => some [92,117,48,48,51,99]
  | 62 => some [92,117,48,48,51,101]
  | 38 => some [92,117,48,48,50,54]
  | _ => none

private def escapeChunks : List Nat → List Nat → List (List Nat)
  | 226 :: 128 :: 168 :: rest, pending =>
    pending.reverse :: [92,117,50,48,50,56] :: escapeChunks rest []
  | 226 :: 128 :: 169 :: rest, pending =>
    pending.reverse :: [92,117,50,48,50,57] :: escapeChunks rest []
  | b :: rest, pending =>
    match escapedByte b with
    | some encoded => pending.reverse :: encoded :: escapeChunks rest []
    | none => escapeChunks rest (b :: pending)
  | [], pending => [pending.reverse]

private def fitting (spare : Nat) : List (List Nat) → List Nat
  | [] => []
  | chunk :: rest =>
    if chunk.length ≤ spare then chunk ++ fitting (spare - chunk.length) rest else []

structure HtmlAppend where
  output : Option (List Nat)
  original : Option (List Nat)
  reused : Bool

/-- Go 1.27.0 src/encoding/json/indent.go and v2_indent.go.
    A growth detaches only at the first append chunk that does not fit. -/
def appendHTMLEscape (backing : Option (List Nat)) (length : Nat) (input : List Nat) :
    HtmlAppend :=
  let original := backing.getD []
  let chunks := escapeChunks input []
  let encoded := chunks.flatten
  let written := fitting (original.length - length) chunks
  let output := original.take length ++ encoded
  ⟨if output.isEmpty then backing.map (fun _ => []) else some output,
   backing.map (fun _ => original.take length ++ written ++ original.drop (length + written.length)),
   !original.isEmpty && output.length ≤ original.length⟩

/-- Go 1.27.0 src/encoding/json/indent.go:16–19 and v2_indent.go:21–24.
    Observe the visible buffer after appending, with no capacity claim. -/
def htmlEscape (receiver : Option (List Nat)) (input : List Nat) : Option (List Nat) :=
  receiver.map fun initial => initial ++ (escapeChunks input []).flatten

theorem appendHTMLEscape_partial :
    let r := appendHTMLEscape (some [80, 99, 99]) 1 [97, 60, 98]
    (r.output, r.original, r.reused) =
      (some [80,97,92,117,48,48,51,99,98], some [80,97,99], false) := rfl

theorem htmlEscape_separator :
    htmlEscape (some [80]) [226,128,168,226,128,169] =
      some [80,92,117,50,48,50,56,92,117,50,48,50,57] := rfl

theorem htmlEscape_nil (input : List Nat) : htmlEscape none input = none := rfl

end GoLean.Encoding.Json
