import GoLean.Compress.Flate.DictionaryInit
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

private def wire (cells : List (BitVec 8)) :=
  String.intercalate "," (cells.map (toString ∘ BitVec.toNat))

def main : IO Unit := runTuple fun row => do
  let present :: length :: capacity :: read :: written :: full :: size :: raw := row |
    throw (IO.userError "expected dictionary configuration")
  let values ← raw.mapM (wordInput 8)
  let state := ⟨values.take capacity.toNat, length.toNat, read.toNat, written.toNat,
    full == 1, present == 1⟩
  let r := dictionaryInit state size.toNat (values.drop capacity.toNat)
  let s := r.next
  return s!"{s.present}|{s.length}|{s.backing.length}|{s.read}|{s.written}|{s.full}|{r.reused}|{wire s.backing}|{wire r.oldBacking}"
