import GoLean.Strconv.AppendBool
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let flag :: dst := fields | throw (IO.userError "expected boolean flag and destination bytes")
  unless (flag = 0 || flag = 1) && dst.all (fun x => 0 ≤ x && x < 256) do
    throw (IO.userError "invalid boolean flag or destination byte")
  let result := GoLean.Strconv.appendBool (dst.map fun x => x.toNat.toUInt8) (flag == 1)
  return String.intercalate " " (result.map toString)
