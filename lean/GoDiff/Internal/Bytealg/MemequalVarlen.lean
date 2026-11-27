import GoLean.Internal.Bytealg.MemequalVarlen
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let size :: leftLength :: raw := row | throw (IO.userError "expected closure byte count and left length")
  let data ← raw.mapM (wordInput 8)
  let a := data.take leftLength.toNat
  let b := data.drop leftLength.toNat
  return toString (GoLean.Internal.Bytealg.runtimeMemequalVarlen size.toNat a b)
