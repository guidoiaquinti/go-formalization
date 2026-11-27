import GoLean.Bytes.BufferReset

namespace GoLean.Bytes

inductive BufferWriteToStatus where
  | ok
  | shortWrite
  | writerError
  | invalidCount

structure BufferWriteToResult where
  status : BufferWriteToStatus
  count : Nat
  call : Option BufferTailView
  state : BufferReadState

/-- Go 1.27.0 Buffer.WriteTo. A scripted writer reports a nonnegative
    count and either nil or an observed error. -/
def bufferWriteTo (receiver : Option BufferReadState) (accepted : Nat) (failed : Bool) :
    Option BufferWriteToResult :=
  receiver.map fun original =>
    let state := {original with lastRead := 0}
    let tail := state.tail
    let unread := (tail.stored.getD []).drop tail.offset
    let reset := (bufferReset (some state)).getD state
    if unread.isEmpty then ⟨.ok, 0, none, reset⟩
    else
      let call := some ⟨unread, tail.spare⟩
      let advanced := {state with tail := {tail with offset := tail.offset + accepted}}
      if unread.length < accepted then ⟨.invalidCount, 0, call, state⟩
      else if failed then ⟨.writerError, accepted, call, advanced⟩
      else if accepted < unread.length then ⟨.shortWrite, accepted, call, advanced⟩
      else ⟨.ok, accepted, call, reset⟩

theorem bufferWriteTo_nil (accepted : Nat) (failed : Bool) :
    bufferWriteTo none accepted failed = none := rfl

theorem bufferWriteTo_invalid :
    bufferWriteTo (some ⟨⟨some [1, 2], [3], 1⟩, 4⟩) 2 true =
      some ⟨.invalidCount, 0, some ⟨[2], [3]⟩, ⟨⟨some [1, 2], [3], 1⟩, 0⟩⟩ := rfl

theorem bufferWriteTo_error :
    bufferWriteTo (some ⟨⟨some [1, 2], [3], 1⟩, 4⟩) 1 true =
      some ⟨.writerError, 1, some ⟨[2], [3]⟩, ⟨⟨some [1, 2], [3], 2⟩, 0⟩⟩ := rfl

theorem bufferWriteTo_reset :
    bufferWriteTo (some ⟨⟨some [1, 2], [3], 1⟩, 4⟩) 1 false =
      some ⟨.ok, 1, some ⟨[2], [3]⟩, ⟨⟨some [], [1, 2, 3], 0⟩, 0⟩⟩ := rfl

end GoLean.Bytes
