import GoLean.Internal.Dag.Topo
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let size :: codes := fields | throw (IO.userError "expected size followed by matrix")
  let n := size.toNat
  let roots := List.range n
  if size < 0 || codes.length != n * n then
    throw (IO.userError "invalid graph matrix")
  let next := fun src => roots.filter fun dest => codes[src * n + dest]?.getD 0 != 0
  return String.intercalate " " ("ok" :: (GoLean.Internal.Dag.topo roots next).map toString)
