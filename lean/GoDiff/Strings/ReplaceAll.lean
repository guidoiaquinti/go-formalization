import GoLean.Strings.ReplaceAll
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let oldWidth :: newWidth :: fields := fields
    | throw (IO.userError "expected old length, new length and bytes")
  let raw ← fields.mapM (wordInput 8)
  let raw := raw.map BitVec.toNat
  let old := raw.take oldWidth.toNat
  let rest := raw.drop oldWidth.toNat
  let replacement := rest.take newWidth.toNat
  let input := rest.drop newWidth.toNat
  let result := GoLean.Strings.replaceAll input old replacement
  return "bytes:" ++ String.intercalate "," (result.map toString)
