import GoLean.Slices.Concat
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let limit :: values := fields | throw (IO.userError "expected length limit and slice parts")
  let parts ← if values.isEmpty then pure [] else do
    let [] :: parts := values.splitOn (-1) | throw (IO.userError "expected parts introduced by -1")
    pure parts
  match GoLean.Slices.sliceConcat limit.toNat parts with
  | .error msg => return msg
  | .ok none => return "nil"
  | .ok (some values) => return String.intercalate " " ("ok" :: values.map toString)
