import GoLean.Internal.Bytealg.Index
import GoDiff.Math.Bits.Tuple

open GoLean.Internal.Bytealg

def main : IO Unit := runTuple fun row => do
  let size :: words := row | throw (IO.userError "expected needle size")
  if size < 0 || size.toNat > words.length then throw (IO.userError "invalid needle size")
  let data := (← words.mapM (wordInput 8)).map BitVec.toNat
  let needle := data.take size.toNat
  let haystack := data.drop size.toNat
  let first := (index haystack needle).map toString |>.getD "-1"
  let second := (indexString haystack needle).map toString |>.getD "-1"
  return first ++ "," ++ second
