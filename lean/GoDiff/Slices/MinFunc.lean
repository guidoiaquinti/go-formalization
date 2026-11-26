import GoLean.Slices.MinFunc
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let 0 :: keys := fields | throw (IO.userError "expected marker and keys")
  let entries := keys.zipIdx
  let compareKey := fun (a b : Int × Nat) => if a.1 < b.1 then -7 else if a.1 = b.1 then 0 else 13
  match GoLean.Slices.sliceMinFunc compareKey entries with
  | .ok entry => return s!"{entry.1} {entry.2}"
  | .error message => return message
