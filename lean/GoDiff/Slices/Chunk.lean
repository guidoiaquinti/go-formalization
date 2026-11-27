import GoLean.Slices.Chunk
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

private def cells (values : List Int) := String.intercalate "," (values.map toString)

private def wire (result : ChunkRun Int) : String :=
  let fragments := result.yields.map fun chunk => s!"{chunk.capacity}:{cells chunk.values}"
  s!"{String.intercalate ";" fragments}/{cells result.backing}"

def main : IO Unit := runTuple fun fields => do
  let length :: size :: first :: second :: write :: replacement :: backing := fields |
    throw (IO.userError "expected length, chunk size, budgets and backing")
  let replacement := if write = 0 then none else some replacement
  match sliceChunk backing length.toNat size first.toNat replacement with
  | .error failure => return failure
  | .ok result =>
    match sliceChunk result.backing length.toNat size second.toNat none with
    | .error failure => return failure
    | .ok replay => return wire result ++ "|" ++ wire replay
