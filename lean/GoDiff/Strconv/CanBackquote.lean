import GoLean.Strconv.CanBackquote
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: input := fields | throw (IO.userError "expected marker followed by input bytes")
  unless input.all (fun x => 0 ≤ x && x < 256) do
    throw (IO.userError "invalid input byte")
  return toString (GoLean.Strconv.canBackquote ⟨input.toArray.map fun x => x.toNat.toUInt8⟩)
