import GoLean.Strings.MakeSingleStringReplacer
import GoLean.Strings.StringFinderNext

namespace GoLean.Strings

structure ReplacerWriteReply where
  limit : Option Nat
  error : Nat

structure ReplacerWriteEvent where
  stringWriter : Bool
  bytes : List (BitVec 8)

structure ReplacerWriteResult where
  written : Nat
  error : Nat
  calls : List ReplacerWriteEvent

private def writeSegments (replacer : SingleStringReplacer) (input : List (BitVec 8)) :
    Nat → List (List (BitVec 8))
  | 0 => [input]
  | fuel + 1 =>
    let matchIndex := replacer.finder.next input
    if matchIndex < 0 then [input]
    else
      let start := matchIndex.toNat
      input.take start :: replacer.value ::
        writeSegments replacer (input.drop (start + replacer.finder.pattern.length)) fuel

private def writeReplies (stringWriter : Bool) (script : List ReplacerWriteReply) :
    List (List (BitVec 8)) → ReplacerWriteResult
  | [] => ⟨0, 0, []⟩
  | payload :: segments =>
    let reply := script.headD ⟨none, 0⟩
    let written := min payload.length (reply.limit.getD payload.length)
    let event := ⟨stringWriter, payload⟩
    if reply.error != 0 then ⟨written, reply.error, [event]⟩
    else
      let rest := writeReplies stringWriter script.tail segments
      ⟨written + rest.written, rest.error, event :: rest.calls⟩

/-- Go 1.27.0 src/strings/replace.go:314–329, 410–433.
    The finder supplies unmatched spans, replacements and one final tail
    write. A deterministic writer selects WriteString or the Write adapter;
    each reply controls the returned count and whether later calls occur. -/
def singleStringReplacerWriteString (pattern value input : List (BitVec 8))
    (stringWriter : Bool) (script : List ReplacerWriteReply) : ReplacerWriteResult :=
  writeReplies stringWriter script
    (writeSegments (makeSingleStringReplacer pattern value) input (input.length + 1))

theorem singleStringReplacerWriteString_noMatch (pattern value input : List (BitVec 8))
    (stringWriter : Bool) (script : List ReplacerWriteReply)
    (missing : stringFinderNext pattern input = -1) :
    (singleStringReplacerWriteString pattern value input stringWriter script).calls =
      [⟨stringWriter, input⟩] := by
  simp only [singleStringReplacerWriteString, writeSegments, makeSingleStringReplacer,
    FinderTables.next_constructed, missing]
  simp only [Int.reduceNeg, Int.reduceLT, ite_true, writeReplies]
  split <;> rfl

theorem singleStringReplacerWriteString_shortSuccess (value : List (BitVec 8))
    (stringWriter : Bool) :
    let actual := singleStringReplacerWriteString [0xff] value [0x7f, 0xff] stringWriter
      [⟨some 0, 0⟩]
    actual.written = value.length ∧ actual.error = 0 ∧
      actual.calls = [⟨stringWriter, [0x7f]⟩, ⟨stringWriter, value⟩, ⟨stringWriter, []⟩] := by
  change let result := writeReplies stringWriter [⟨some 0, 0⟩] [[0x7f], value, []]
         result.written = value.length ∧ result.error = 0 ∧
           result.calls = [⟨stringWriter, [0x7f]⟩, ⟨stringWriter, value⟩, ⟨stringWriter, []⟩]
  simp [writeReplies]

theorem singleStringReplacerWriteString_firstError (value : List (BitVec 8))
    (stringWriter : Bool) (reply : ReplacerWriteReply) (script : List ReplacerWriteReply)
    (failed : reply.error ≠ 0) :
    let actual := singleStringReplacerWriteString [0xff] value [0xff] stringWriter (reply :: script)
    actual.written = 0 ∧ actual.error = reply.error ∧ actual.calls = [⟨stringWriter, []⟩] := by
  change let result := writeReplies stringWriter (reply :: script) [[], value, []]
         result.written = 0 ∧ result.error = reply.error ∧ result.calls = [⟨stringWriter, []⟩]
  simp [writeReplies, failed]

end GoLean.Strings
