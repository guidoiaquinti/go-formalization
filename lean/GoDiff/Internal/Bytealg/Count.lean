import GoLean.Internal.Bytealg.Count
import GoDiff.Math.Bits.Tuple

open GoLean.Internal.Bytealg

def main : IO Unit := runTuple fun row => do
  let target :: values := row | throw (IO.userError "expected target byte")
  let target ← wordInput 8 target
  let bytes ← values.mapM (wordInput 8)
  return String.intercalate "," ([count bytes target, countString bytes target,
    countGeneric bytes target, countGenericString bytes target].map toString)
