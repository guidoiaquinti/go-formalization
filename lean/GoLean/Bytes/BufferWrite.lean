import GoLean.Bytes.BufferWriteByte

namespace GoLean.Bytes

inductive BufferWriteInput where
  | detached (bytes : List (BitVec 8))
  | backing (start length : Nat)

private def BufferWriteInput.length : BufferWriteInput → Nat
  | .detached bytes => bytes.length
  | .backing _ length => length

private def BufferWriteInput.read (backing : List (BitVec 8)) : BufferWriteInput → List (BitVec 8)
  | .detached bytes => bytes
  | .backing start length => (backing.drop start).take length

private structure WriteSpace where
  retained : List (BitVec 8)
  backing : List (BitVec 8)
  offset : Nat
  reused : Bool

private def writeSpace (tail : BufferTailState) (count : Nat) : WriteSpace :=
  let stored := tail.stored.getD []
  let backing := stored ++ tail.spare
  if count ≤ tail.spare.length then ⟨stored, backing, tail.offset, true⟩
  else
    let unread := stored.drop tail.offset
    let reset := unread.isEmpty && tail.offset > 0
    let capacity := backing.length
    let length := if reset then 0 else stored.length
    let offset := if reset then 0 else tail.offset
    if reset && count ≤ capacity then ⟨[], backing, 0, true⟩
    else if tail.stored.isNone && count ≤ 64 then ⟨[], List.replicate 64 0, 0, false⟩
    else if unread.length + count ≤ capacity / 2 then
      ⟨unread, unread ++ backing.drop unread.length, 0, true⟩
    else
      let requested := max (length + count) (2 * (capacity - offset))
      ⟨unread, List.replicate (smallByteCapacity requested) 0, 0, false⟩

/-- Go 1.27.0 Buffer.Write. An aliased source is observed after buffer
    growth or compaction prepares space, then copied as an overlapping copy.
    The second component of a present result denotes successful buffer state. -/
def bufferWrite (receiver : Option BufferReadState) (source : BufferWriteInput) :
    Option (Nat × BufferWriteByteResult) :=
  receiver.map fun state =>
    let count := source.length
    let space := writeSpace state.tail count
    let original := state.tail.stored.getD [] ++ state.tail.spare
    let payload := source.read (if space.reused then space.backing else original)
    let stored := space.retained ++ payload
    let raw := if count == 0 then state.tail.stored else some stored
    let result : BufferReadState := ⟨⟨raw, space.backing.drop stored.length, space.offset⟩, 0⟩
    (payload.length, ⟨result, space.reused && !space.backing.isEmpty⟩)

theorem bufferWrite_nil (source : BufferWriteInput) : bufferWrite none source = none := rfl

theorem bufferWrite_empty (offset : Nat) (flag : Int) :
    bufferWrite (some ⟨⟨none, [], offset⟩, flag⟩) (.detached []) =
      some (0, ⟨⟨⟨none, [], offset⟩, 0⟩, false⟩) := rfl

theorem bufferWrite_compactAlias :
    bufferWrite (some ⟨⟨some [1, 2, 3, 4, 5, 6, 7, 8], [], 6⟩, 1⟩) (.backing 0 2) =
      some (2, ⟨⟨⟨some [7, 8, 7, 8], [5, 6, 7, 8], 0⟩, 0⟩, true⟩) := rfl

theorem bufferWrite_allocateAlias :
    bufferWrite (some ⟨⟨some [1, 2], [], 0⟩, 1⟩) (.backing 0 2) =
      some (2, ⟨⟨⟨some [1, 2, 1, 2], [0, 0, 0, 0], 0⟩, 0⟩, false⟩) := rfl

end GoLean.Bytes
