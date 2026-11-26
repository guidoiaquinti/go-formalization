import GoLean.Slices.Clone
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let nilFlag :: _spare :: elements := fields | throw (IO.userError "expected nil flag, spare capacity and elements")
  let input := if nilFlag = 1 then none else some elements
  match GoLean.Slices.sliceClone input with
  | none => return "true"
  | some xs => return String.intercalate " " ("false" :: xs.map toString)
