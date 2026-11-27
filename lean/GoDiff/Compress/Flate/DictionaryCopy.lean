import GoLean.Compress.Flate.DictionaryCopy
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

private def wire (r : Nat × DictionaryState) :=
  let (count, s) := r
  let bytes := String.intercalate "," (s.backing.map (toString ∘ BitVec.toNat))
  s!"{count}:{s.read}:{s.written}:{s.full}:{bytes}"

def main : IO Unit := runTuple fun row => do
  let length :: read :: written :: full :: dist :: count :: raw := row |
    throw (IO.userError "expected dictionary and copy request")
  let backing ← raw.mapM (wordInput 8)
  let s := ⟨backing, length.toNat, read.toNat, written.toNat, full == 1, true⟩
  return wire (writeCopy s dist.toNat count.toNat) ++ "|" ++ wire (tryWriteCopy s dist.toNat count.toNat)
