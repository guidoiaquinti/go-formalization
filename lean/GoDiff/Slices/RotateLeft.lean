import GoLean.Slices.RotateLeft
import GoDiff.Math.Bits.Tuple

open GoLean.Slices

def main : IO Unit := runTuple fun fields => do
  let rotation :: length :: cells := fields | throw (IO.userError "expected rotation, length and capacity cells")
  let (ok, out) := sliceRotateLeft cells length.toNat rotation
  return String.intercalate " " ((if ok then "ok" else "panic") :: out.map toString)
