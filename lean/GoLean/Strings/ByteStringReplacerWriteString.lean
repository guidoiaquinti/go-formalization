import GoLean.Strings.SingleStringReplacerWriteString

namespace GoLean.Strings

private def unchangedEvent (stringWriter : Bool) (reversed : List (BitVec 8)) :
    List ReplacerWriteEvent :=
  if reversed.isEmpty then [] else [⟨stringWriter, reversed.reverse⟩]

private def byteStringWriteEvents (lookup : BitVec 8 → Option (List (BitVec 8))) (stringWriter : Bool)
    (reversed : List (BitVec 8)) : List (BitVec 8) → List ReplacerWriteEvent
  | [] => unchangedEvent stringWriter reversed
  | byte :: tail =>
    match lookup byte with
    | none => byteStringWriteEvents lookup stringWriter (byte :: reversed) tail
    | some replacement => unchangedEvent stringWriter reversed ++
        ⟨false, replacement⟩ :: byteStringWriteEvents lookup stringWriter [] tail

private def applyByteStringWriteReplies (script : List ReplacerWriteReply) :
    List ReplacerWriteEvent → ReplacerWriteResult
  | [] => ⟨0, 0, []⟩
  | event :: events =>
    let reply := script.headD ⟨none, 0⟩
    let size := event.bytes.length
    let written := min size (reply.limit.getD size)
    if reply.error != 0 then ⟨written, reply.error, [event]⟩
    else
      let rest := applyByteStringWriteReplies script.tail events
      ⟨written + rest.written, rest.error, event :: rest.calls⟩

/-- Go 1.27.0 src/strings/replace.go:550–577. Send only nonempty unchanged
    spans through the selected string writer. Send each present replacement through
    the underlying Write method directly, even when it is empty. Stop at the first error. -/
def byteStringReplacerWriteString (mapping : Option (BitVec 8 → Option (List (BitVec 8))))
    (input : List (BitVec 8)) (stringWriter : Bool) (script : List ReplacerWriteReply) :
    Option ReplacerWriteResult :=
  match input, mapping with
  | [], _ => some ⟨0, 0, []⟩
  | _ :: _, none => none
  | _ :: _, some lookup => some (applyByteStringWriteReplies script (byteStringWriteEvents lookup stringWriter [] input))

theorem byteStringReplacerWriteString_empty (mapping : Option (BitVec 8 → Option (List (BitVec 8))))
    (stringWriter : Bool) (script : List ReplacerWriteReply) :
    byteStringReplacerWriteString mapping [] stringWriter script = some ⟨0, 0, []⟩ := rfl

theorem byteStringReplacerWriteString_nil (head : BitVec 8) (tail : List (BitVec 8))
    (stringWriter : Bool) (script : List ReplacerWriteReply) :
    byteStringReplacerWriteString none (head :: tail) stringWriter script = none := rfl

theorem byteStringReplacerWriteString_changed (lookup : BitVec 8 → Option (List (BitVec 8)))
    (byte : BitVec 8) (replacement : List (BitVec 8)) (stringWriter : Bool)
    (present : lookup byte = some replacement) :
    byteStringReplacerWriteString (some lookup) [byte] stringWriter [] =
      some ⟨replacement.length, 0, [⟨false, replacement⟩]⟩ := by
  simp [byteStringReplacerWriteString, byteStringWriteEvents, unchangedEvent, applyByteStringWriteReplies, present]

theorem byteStringReplacerWriteString_firstError :
    let lookup := fun (byte : BitVec 8) => if byte = 98 then some [] else none
    byteStringReplacerWriteString (some lookup) [97, 98, 97] true [⟨none, 0⟩, ⟨none, 1⟩] =
      some ⟨1, 1, [⟨true, [97]⟩, ⟨false, []⟩]⟩ := by rfl

theorem byteStringReplacerWriteString_shortSuccess :
    let lookup := fun (byte : BitVec 8) => if byte = 98 then some [0xff, 0] else none
    byteStringReplacerWriteString (some lookup) [98, 98, 97] true [⟨some 0, 0⟩] =
      some ⟨3, 0, [⟨false, [0xff, 0]⟩, ⟨false, [0xff, 0]⟩, ⟨true, [97]⟩]⟩ := by rfl

end GoLean.Strings
