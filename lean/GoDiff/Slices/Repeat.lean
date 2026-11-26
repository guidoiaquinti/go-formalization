import GoLean.Slices.Repeat
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let limit :: count :: values := fields | throw (IO.userError "expected length limit, count and elements")
  match GoLean.Slices.sliceRepeat limit.toNat values count with
  | .error message => return message
  | .ok result => return String.intercalate " " ("ok" :: result.map toString)
