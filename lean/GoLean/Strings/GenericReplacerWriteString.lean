import GoLean.Strings.GenericReplacerLookup
import GoLean.Strings.SingleStringReplacerWriteString

namespace GoLean.Strings

private def genericWriteSegments
    (pairs : List (List (BitVec 8) × List (BitVec 8)))
    (pending : List (BitVec 8)) (skipRoot : Bool) (remaining : List (BitVec 8)) :
    Nat → List (List (BitVec 8))
  | 0 => []
  | fuel + 1 =>
    match genericReplacerLookup pairs remaining skipRoot with
    | some found =>
      pending.reverse :: found.value ::
        genericWriteSegments pairs [] (found.keyLength == 0) (remaining.drop found.keyLength) fuel
    | none =>
      match remaining with
      | [] => if pending.isEmpty then [] else [pending.reverse]
      | byte :: rest => genericWriteSegments pairs (byte :: pending) false rest fuel

/-- Go 1.27.0 src/strings/replace.go:336–374. Every match issues an
    unchanged-span write and a replacement write, including empty payloads.
    An empty match can cause one additional lookup at the same byte index. -/
def genericReplacerWriteSegments (pairs : List (List (BitVec 8) × List (BitVec 8)))
    (input : List (BitVec 8)) : List (List (BitVec 8)) :=
  genericWriteSegments pairs [] false input (2 * (input.length + 1))

private def genericWriteReplies (stringWriter : Bool) (script : List ReplacerWriteReply) :
    List (List (BitVec 8)) → ReplacerWriteResult
  | [] => ⟨0, 0, []⟩
  | payload :: remaining =>
    let event := ⟨stringWriter, payload⟩
    let reply := script.headD ⟨none, 0⟩
    let written := min payload.length (reply.limit.getD payload.length)
    if reply.error != 0 then ⟨written, reply.error, [event]⟩
    else
      let rest := genericWriteReplies stringWriter script.tail remaining
      ⟨written + rest.written, rest.error, event :: rest.calls⟩

/-- Observe the dispatch and count from the selected string writer or
    its byte adapter. A short successful reply does not stop the scan. -/
def genericReplacerWriteString (pairs : List (List (BitVec 8) × List (BitVec 8)))
    (input : List (BitVec 8)) (stringWriter : Bool)
    (script : List ReplacerWriteReply) : ReplacerWriteResult :=
  genericWriteReplies stringWriter script (genericReplacerWriteSegments pairs input)

theorem genericReplacerWriteString_noPairsEmpty (stringWriter : Bool)
    (script : List ReplacerWriteReply) :
    genericReplacerWriteString [] [] stringWriter script = ⟨0, 0, []⟩ := rfl

theorem genericReplacerWriteString_rootThenKey :
    genericReplacerWriteString [([], [1]), ([0xff], [2])] [0xff] true [] =
      ⟨3, 0, [⟨true, []⟩, ⟨true, [1]⟩, ⟨true, []⟩, ⟨true, [2]⟩,
        ⟨true, []⟩, ⟨true, [1]⟩]⟩ := rfl

theorem genericReplacerWriteString_emptyReplacementError :
    genericReplacerWriteString [([97], [])] [97, 97] false [⟨none, 0⟩, ⟨none, 2⟩] =
      ⟨0, 2, [⟨false, []⟩, ⟨false, []⟩]⟩ := rfl

theorem genericReplacerWriteString_shortSuccess :
    genericReplacerWriteString [([98], [0, 255])] [97, 98, 99] true [⟨some 0, 0⟩] =
      ⟨3, 0, [⟨true, [97]⟩, ⟨true, [0, 255]⟩, ⟨true, [99]⟩]⟩ := rfl

end GoLean.Strings
