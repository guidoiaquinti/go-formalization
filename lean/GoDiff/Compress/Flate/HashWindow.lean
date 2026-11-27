import GoLean.Compress.Flate.HashWindow
import GoDiff.Math.Bits.Tuple

open GoLean.Compress.Flate

def main : IO Unit := runTuple fun row => do
  let size :: cap :: dstSize :: values := row | throw (IO.userError "expected slice geometries")
  let source ← (values.take cap.toNat).mapM (wordInput 8)
  let dst ← (values.drop cap.toNat).mapM (wordInput 32)
  let single := (hash4 source).map (toString ∘ BitVec.toNat) |>.getD "bounds"
  let (failure, backing) := bulkHash4 (source.take size.toNat) dstSize.toNat dst
  let cells := String.intercalate "," (backing.map (toString ∘ BitVec.toNat))
  return s!"{single}|{failure}|{cells}"
