import GoLean.Bytes.BufferWriteByte
import GoLean.Unicode.Utf8.AppendRune

namespace GoLean.Bytes

private def runeWriteState (kept encoded backing : List (BitVec 8)) (offset : Nat) (reused : Bool) :
    BufferWriteByteResult :=
  let stored := kept ++ encoded
  ⟨⟨⟨some stored, backing.drop stored.length, offset⟩, 0⟩, reused⟩

private def writeReservedRune (tail : BufferTailState) (encoded : List (BitVec 8)) :
    BufferWriteByteResult :=
  let stored := tail.stored.getD []
  let backing := stored ++ tail.spare
  if 4 ≤ tail.spare.length then
    runeWriteState stored encoded backing tail.offset true
  else
    let unread := stored.drop tail.offset
    let reset := unread.isEmpty && tail.offset > 0
    let capacity := backing.length
    let size := if reset then 0 else stored.length
    let offset := if reset then 0 else tail.offset
    if reset && 4 ≤ capacity then
      runeWriteState [] encoded backing 0 true
    else if tail.stored.isNone then
      runeWriteState [] encoded (List.replicate 64 0) 0 false
    else if unread.length + 4 ≤ capacity / 2 then
      runeWriteState unread encoded backing 0 true
    else
      let requested := max (size + 4) (2 * (capacity - offset))
      runeWriteState unread encoded (List.replicate (smallByteCapacity requested) 0) 0 false

/-- Pinned Go 1.27.0 Buffer.WriteRune. The non-ASCII path reserves four
    bytes even when the actual encoding is shorter. The result's width is the
    number of bytes written, and a present result denotes a nil error. -/
def bufferWriteRune (receiver : Option BufferReadState) (rune : Int) :
    Option (Nat × BufferWriteByteResult) :=
  if 0 ≤ rune && rune < 128 then
    (bufferWriteByte receiver (BitVec.ofInt 8 rune)).map (fun result => (1, result))
  else
    let encoded := (GoLean.Unicode.Utf8.appendRune [] rune).map (BitVec.ofNat 8)
    receiver.map fun state => (encoded.length, writeReservedRune state.tail encoded)

theorem bufferWriteRune_nil (rune : Int) : bufferWriteRune none rune = none := by
  simp [bufferWriteRune, bufferWriteByte]

theorem bufferWriteRune_ascii (receiver : Option BufferReadState) (rune : Int)
    (nonnegative : 0 ≤ rune) (ascii : rune < 128) :
    bufferWriteRune receiver rune =
      (bufferWriteByte receiver (BitVec.ofInt 8 rune)).map (fun result => (1, result)) := by
  simp [bufferWriteRune, nonnegative, ascii]

theorem bufferWriteRune_replacement :
    bufferWriteRune (some ⟨⟨none, [], 0⟩, 1⟩) (-1) =
      some (3, ⟨⟨⟨some [239, 191, 189], List.replicate 61 0, 0⟩, 0⟩, false⟩) := rfl

theorem bufferWriteRune_reserveFour :
    bufferWriteRune (some ⟨⟨some [1, 2], [3, 4], 0⟩, 1⟩) 233 =
      some (2, ⟨⟨⟨some [1, 2, 195, 169], [0, 0, 0, 0], 0⟩, 0⟩, false⟩) := rfl

theorem bufferWriteRune_reslice :
    bufferWriteRune (some ⟨⟨some [1, 2], [3, 4, 5, 6], 1⟩, 2⟩) 233 =
      some (2, ⟨⟨⟨some [1, 2, 195, 169], [5, 6], 1⟩, 0⟩, true⟩) := rfl

end GoLean.Bytes
