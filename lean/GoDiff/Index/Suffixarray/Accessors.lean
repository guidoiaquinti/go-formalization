import GoLean.Index.Suffixarray.Accessors
import GoDiff.Math.Bits.Tuple

open GoLean.Index.Suffixarray

private def view (v : View α) (render : α → String) : String :=
  match v.backing with
  | none => "nil"
  | some backing => s!"{v.length}:{String.intercalate "," (backing.map render)}"

private def ints (s : Ints) : String :=
  view s.words32 (toString ∘ BitVec.toInt) ++ "/" ++ view s.words64 (toString ∘ BitVec.toInt)

private def result (r : Except AccessFailure α) (render : α → String) : String :=
  match r with
  | .ok value => "ok:" ++ render value
  | .error .nilReceiver => "nilReceiver"
  | .error .bounds => "bounds"

def main : IO Unit := runTuple fun fields => do
  let tag :: len32 :: cap32 :: len64 :: cap64 :: byteLen :: i :: j :: value :: raw := fields |
    throw (IO.userError "expected variant, view geometry, indices and backing")
  let words32 := (raw.take cap32.toNat).map (BitVec.ofInt 32)
  let words64 := ((raw.drop cap32.toNat).take cap64.toNat).map (BitVec.ofInt 64)
  let bytes ← (raw.drop (cap32 + cap64).toNat).mapM (wordInput 8)
  let s : Ints := ⟨
    ⟨if tag.toNat.testBit 0 then some words32 else none, len32.toNat⟩,
    ⟨if tag.toNat.testBit 1 then some words64 else none, len64.toNat⟩⟩
  let data : View (BitVec 8) := ⟨if tag.toNat.testBit 2 then some bytes else none, byteLen.toNat⟩
  let index := if tag.toNat.testBit 3 then none else some (⟨data, s⟩ : IndexState)
  let receiver := if tag.toNat.testBit 3 then none else some s
  let renderedBytes (v : View (BitVec 8)) := view v (toString ∘ BitVec.toNat)
  return String.intercalate "|" [
    result (intsLen receiver) toString, result (intsGet receiver i) toString,
    result (intsSet receiver i value) ints, result (intsSlice receiver i j) ints,
    result (indexBytes index) renderedBytes, result (indexAt index i) renderedBytes]
