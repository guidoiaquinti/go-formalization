import GoLean.Slices.RotateRight
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

def main : IO Unit := runTuple fun fields => do
  let width :: rotation :: length :: cells := fields |
    throw (IO.userError "expected native width, rotation, length and capacity cells")
  let (ok, out) := sliceRotateRight cells length.toNat width.toNat rotation
  return String.intercalate " " ((if ok then "ok" else "panic") :: out.map toString)
