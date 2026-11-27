import GoLean.Compress.Flate.DictionaryState
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

private def wire (bytes : List (BitVec 8)) :=
  String.intercalate "," (bytes.map (toString ∘ BitVec.toNat))

private def viewWire (v : DictionaryView) :=
  s!"{v.present}:{v.capacity}:{v.origin}:{wire v.bytes}"

private def stateWire (s : DictionaryState) :=
  s!"{s.read}:{s.written}:{s.full}:{wire s.backing}"

def main : IO Unit := runTuple fun row => do
  let present :: length :: read :: written :: full :: count :: value :: raw := row |
    throw (IO.userError "expected dictionary state")
  let byte ← wordInput 8 value
  let backing ← raw.mapM (wordInput 8)
  let state := ⟨backing, length.toNat, read.toNat, written.toNat, full == 1, present == 1⟩
  let (flushed, after) := readFlush state
  let stored := (writeByte state byte).map stateWire |>.getD "bounds"
  return String.intercalate "|" [toString (histSize state), toString (availRead state),
    toString (availWrite state), viewWire (writeSlice state), stateWire (writeMark state count.toNat),
    stored, viewWire flushed, stateWire after]
