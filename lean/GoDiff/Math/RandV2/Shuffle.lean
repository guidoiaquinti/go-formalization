import GoLean.Math.RandV2.Shuffle
import GoDiff.Math.Bits.Tuple

open GoLean.Math.RandV2

def main : IO Unit := runTuple fun row => do
  let width :: n :: mode :: failAt :: raw := row |
    throw (IO.userError "expected word width, length, script and source words")
  let words ← raw.mapM fun w => BitVec.toNat <$> wordInput 64 w
  let script := if mode == 1 then SwapScript.failAt failAt.toNat
    else if mode == 2 then SwapScript.missing else SwapScript.normal
  let r := randShuffle width.toNat n words script
  let (s, pulls, perm) := randPerm width.toNat n words
  let swapWire := (r.calls.map fun (i, j) => s!"{i},{j}")
  let values := (perm.map (fun values => "slice:" ++ String.intercalate "," (values.map toString))).getD "none"
  return s!"{r.status}:{r.pulls}:{String.intercalate ";" swapWire}|{s}:{pulls}:{values}"
