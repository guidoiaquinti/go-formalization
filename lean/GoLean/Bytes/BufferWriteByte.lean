import GoLean.Bytes.BufferReset

namespace GoLean.Bytes

structure BufferWriteByteResult where
  state : BufferReadState
  reused : Bool

private def writeByteState (stored spare : List (BitVec 8)) (offset : Nat)
    (reused : Bool) : BufferWriteByteResult :=
  ⟨⟨⟨some stored, spare, offset⟩, 0⟩, reused⟩

/-- Pinned noscan classes through the upper bound of the illustrated domain. -/
def smallByteCapacity (requested : Nat) : Nat :=
  let classes := [8, 16, 24, 32, 48, 64, 80, 96, 112, 128, 144, 160, 176, 192,
    208, 224, 240, 256, 288, 320, 352, 384, 416, 448, 480, 512]
  (classes.find? (requested ≤ ·)).getD requested

/-- Go 1.27.0 src/bytes/buffer.go:303–311 and the called growth paths.
    Results model successful allocation through raw request 512.
    None represents a nil receiver; a result represents a nil error. -/
def bufferWriteByte (receiver : Option BufferReadState) (byte : BitVec 8) :
    Option BufferWriteByteResult :=
  receiver.map fun state =>
    let tail := state.tail
    let stored := tail.stored.getD []
    match tail.spare with
    | _ :: spare => writeByteState (stored ++ [byte]) spare tail.offset true
    | [] =>
      if tail.stored.isNone then writeByteState [byte] (List.replicate 63 0) 0 false
      else
        let unread := stored.drop tail.offset
        let length := unread.length
        let capacity := stored.length
        if unread.isEmpty && tail.offset > 0 then
          writeByteState [byte] (stored.drop 1) 0 true
        else if length + 1 ≤ capacity / 2 then
          writeByteState (unread ++ [byte]) (stored.drop (length + 1)) 0 true
        else
          let requested := max (capacity + 1) (2 * (capacity - tail.offset))
          writeByteState (unread ++ [byte]) (List.replicate (smallByteCapacity requested - length - 1) 0) 0 false

theorem bufferWriteByte_nil (byte : BitVec 8) : bufferWriteByte none byte = none := rfl

theorem bufferWriteByte_reslice (stored spare : List (BitVec 8)) (discard byte : BitVec 8)
    (offset : Nat) (flag : Int) :
    bufferWriteByte (some ⟨⟨some stored, discard :: spare, offset⟩, flag⟩) byte =
      some ⟨⟨⟨some (stored ++ [byte]), spare, offset⟩, 0⟩, true⟩ := rfl

theorem bufferWriteByte_nilStorage (byte : BitVec 8) (flag : Int) :
    bufferWriteByte (some ⟨⟨none, [], 0⟩, flag⟩) byte =
      some ⟨⟨⟨some [byte], List.replicate 63 0, 0⟩, 0⟩, false⟩ := rfl

theorem bufferWriteByte_compactExample :
    bufferWriteByte (some ⟨⟨some [1, 2, 3, 4, 5, 6], [], 4⟩, 1⟩) 9 =
      some ⟨⟨⟨some [5, 6, 9], [4, 5, 6], 0⟩, 0⟩, true⟩ := rfl

theorem bufferWriteByte_allocateExample :
    bufferWriteByte (some ⟨⟨some [1, 2], [], 1⟩, 1⟩) 9 =
      some ⟨⟨⟨some [2, 9], [0, 0, 0, 0, 0, 0], 0⟩, 0⟩, false⟩ := rfl

end GoLean.Bytes
