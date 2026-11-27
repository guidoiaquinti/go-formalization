import GoLean.Compress.Flate.LittleEndian
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

private def render (x : Option Nat) : String := x.map toString |>.getD "bounds"

def main : IO Unit := runTuple fun row => do
  let width :: signed :: start :: length :: value :: raw := row |
    throw (IO.userError "expected index and backing")
  let word ← wordInput 64 value
  let backing ← raw.mapM (wordInput 8)
  let (stored, updated) := storeLE64 backing length.toNat word
  let bytes := String.intercalate "," (updated.map (toString ∘ BitVec.toNat))
  return String.intercalate "|" [render (loadLE8 backing length.toNat start),
    render (loadLE32 width.toNat (signed == 1) backing start),
    render (loadLE64 width.toNat (signed == 1) backing start), toString stored, bytes]
