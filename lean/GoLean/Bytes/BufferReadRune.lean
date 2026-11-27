import GoLean.Bytes.BufferReset
import GoLean.Unicode.Utf8.DecodeRune

namespace GoLean.Bytes

structure BufferReadRuneResult where
  rune : Nat
  width : Nat
  eof : Bool
  state : BufferReadState

/-- Go 1.27.0 src/bytes/buffer.go:387–406. Decode the unread bytes
    and record the consumed width as the read-operation tag.
    An exhausted buffer resets, and None denotes a nil receiver panic. -/
def bufferReadRune (receiver : Option BufferReadState) : Option BufferReadRuneResult :=
  match receiver with
  | none => none
  | some state =>
    let tail := state.tail
    let unread := (tail.stored.getD []).drop tail.offset
    match unread with
    | [] =>
      (bufferReset (some state)).map (fun reset => ⟨0, 0, true, reset⟩)
    | head :: rest =>
      let (rune, width) := GoLean.Unicode.Utf8.decodeRune ((head :: rest).map BitVec.toNat)
      some ⟨rune, width, false, ⟨{tail with offset := tail.offset + width}, Int.ofNat width⟩⟩

theorem bufferReadRune_nil : bufferReadRune none = none := rfl

theorem bufferReadRune_eof (state : BufferReadState)
    (exhausted : (state.tail.stored.getD []).drop state.tail.offset = []) :
    bufferReadRune (some state) =
      (bufferReset (some state)).map (fun reset => ⟨0, 0, true, reset⟩) := by
  dsimp only [bufferReadRune]
  rw [exhausted]

theorem bufferReadRune_live (state : BufferReadState) (head : BitVec 8) (rest : List (BitVec 8))
    (unread : (state.tail.stored.getD []).drop state.tail.offset = head :: rest) :
    bufferReadRune (some state) =
      let (rune, width) := GoLean.Unicode.Utf8.decodeRune ((head :: rest).map BitVec.toNat)
      some ⟨rune, width, false,
        ⟨{state.tail with offset := state.tail.offset + width}, Int.ofNat width⟩⟩ := by
  dsimp only [bufferReadRune]
  rw [unread]

end GoLean.Bytes
