import GoLean.Slices.Delete
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let i :: j :: xs := fields | throw (IO.userError "expected bounds and elements")
  let (result, backing) := GoLean.Slices.sliceDelete (0 : Int) xs i j
  let header := match result with
    | none => "panic"
    | some kept => toString kept.length
  return String.intercalate " " (header :: backing.map toString)
