import GoLean.Strings.SingleStringReplacerWriteString

namespace GoLean.Strings

private def unchangedEvent (stringWriter : Bool) (reversed : List (BitVec 8)) :
    List ReplacerWriteEvent :=
  if reversed.isEmpty then [] else [⟨stringWriter, reversed.reverse⟩]

private def byteWriteEvents (lookup : BitVec 8 → BitVec 8) (stringWriter : Bool)
    (reversed : List (BitVec 8)) : List (BitVec 8) → List ReplacerWriteEvent
  | [] => unchangedEvent stringWriter reversed
  | byte :: tail =>
    let mapped := lookup byte
    if mapped = byte then byteWriteEvents lookup stringWriter (byte :: reversed) tail
    else unchangedEvent stringWriter reversed ++
      ⟨false, [mapped]⟩ :: byteWriteEvents lookup stringWriter [] tail

private def applyByteWriteReplies (script : List ReplacerWriteReply) :
    List ReplacerWriteEvent → ReplacerWriteResult
  | [] => ⟨0, 0, []⟩
  | event :: events =>
    let reply := script.headD ⟨none, 0⟩
    let size := event.bytes.length
    let written := min size (reply.limit.getD size)
    if reply.error != 0 then ⟨written, reply.error, [event]⟩
    else
      let rest := applyByteWriteReplies script.tail events
      ⟨written + rest.written, rest.error, event :: rest.calls⟩

/-- Go 1.27.0 src/strings/replace.go:457–487. Send only nonempty unchanged
    spans through the selected string writer. Send each changed byte through
    the underlying Write method directly, and stop at the first error. -/
def byteReplacerWriteString (mapping : Option (BitVec 8 → BitVec 8))
    (input : List (BitVec 8)) (stringWriter : Bool) (script : List ReplacerWriteReply) :
    Option ReplacerWriteResult :=
  match input, mapping with
  | [], _ => some ⟨0, 0, []⟩
  | _ :: _, none => none
  | _ :: _, some lookup => some (applyByteWriteReplies script (byteWriteEvents lookup stringWriter [] input))

theorem byteReplacerWriteString_empty (mapping : Option (BitVec 8 → BitVec 8))
    (stringWriter : Bool) (script : List ReplacerWriteReply) :
    byteReplacerWriteString mapping [] stringWriter script = some ⟨0, 0, []⟩ := rfl

theorem byteReplacerWriteString_nil (head : BitVec 8) (tail : List (BitVec 8))
    (stringWriter : Bool) (script : List ReplacerWriteReply) :
    byteReplacerWriteString none (head :: tail) stringWriter script = none := rfl

theorem byteReplacerWriteString_changed (lookup : BitVec 8 → BitVec 8)
    (byte : BitVec 8) (stringWriter : Bool) (changed : lookup byte ≠ byte) :
    byteReplacerWriteString (some lookup) [byte] stringWriter [] =
      some ⟨1, 0, [⟨false, [lookup byte]⟩]⟩ := by
  simp [byteReplacerWriteString, byteWriteEvents, unchangedEvent, applyByteWriteReplies, changed]

theorem byteReplacerWriteString_firstError :
    let lookup := fun (byte : BitVec 8) => if byte = 98 then 0xff else byte
    byteReplacerWriteString (some lookup) [97, 98, 97] true [⟨none, 0⟩, ⟨none, 1⟩] =
      some ⟨2, 1, [⟨true, [97]⟩, ⟨false, [0xff]⟩]⟩ := by rfl

theorem byteReplacerWriteString_shortSuccess :
    let lookup := fun (byte : BitVec 8) => if byte = 98 then 0xff else byte
    byteReplacerWriteString (some lookup) [98, 98, 97] true [⟨some 0, 0⟩] =
      some ⟨2, 0, [⟨false, [0xff]⟩, ⟨false, [0xff]⟩, ⟨true, [97]⟩]⟩ := by rfl

end GoLean.Strings
