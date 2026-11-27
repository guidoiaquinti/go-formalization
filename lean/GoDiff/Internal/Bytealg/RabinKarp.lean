import GoLean.Internal.Bytealg.RabinKarp
import GoDiff.Math.Bits.Tuple

open GoLean.Internal.Bytealg

def main : IO Unit := runTuple fun row => do
  let size :: raw := row | throw (IO.userError "expected separator length")
  if size < 0 || size.toNat > raw.length then throw (IO.userError "invalid separator length")
  let data ← raw.mapM (wordInput 8)
  let separator := data.take size.toNat
  let input := data.drop size.toNat
  let (forward, power) := hashStr separator
  let (backward, reversePower) := hashStrRev separator
  let first := (indexRabinKarp input separator).map toString |>.getD "panic"
  let last := (lastIndexRabinKarp input separator).map toString |>.getD "panic"
  return s!"{forward.toNat},{power.toNat},{backward.toNat},{reversePower.toNat},{first},{last}"
