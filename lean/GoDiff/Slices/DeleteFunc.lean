import GoLean.Slices.DeleteFunc
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let mode :: xs := fields | throw (IO.userError "expected mode and elements")
  let predicate := fun (x : Int) =>
    if mode != 0 && x == 13 then .error () else .ok (x % 2 == 0)
  let (result, backing, trace) := GoLean.Slices.sliceDeleteFunc predicate 0 xs
  let header := match result with
    | none => "panic"
    | some kept => toString kept.length
  return String.intercalate " " (header :: backing.map toString ++ "|" :: trace.map toString)
