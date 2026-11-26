import GoLean.Slices.BinarySearchFunc
import GoDiff.Math.Bits.Tuple

structure SearchKey where
  value : Int
  descending : Bool

def compareElement (element : Int) (key : SearchKey) : Int :=
  if element = key.value then 0
  else if (decide (element < key.value) != key.descending) then -7 else 13

def main : IO Unit := runTuple fun fields => do
  let direction :: key :: xs := fields | throw (IO.userError "expected direction and key")
  let (index, found) :=
    GoLean.Slices.binarySearchFunc compareElement 0 xs ⟨key, direction != 0⟩
  return s!"{index} {found}"
