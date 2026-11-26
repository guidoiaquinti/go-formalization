import GoLean.Slices.BinarySearch
import Init.Data.Float.Float
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let kind :: values := fields | throw (IO.userError "expected element type and search data")
  let (index, found) ← if kind = 0 then do
    let key :: xs := values | throw (IO.userError "expected integer key")
    pure (GoLean.Slices.binarySearch (fun a b => decide (a < b)) (· == ·) 0 xs key)
  else do
    let words ← values.mapM fun n => do
      return Float.ofBits (UInt64.ofBitVec (← wordInput 64 n))
    let key :: xs := words | throw (IO.userError "expected binary64 key")
    let less := fun (a b : Float) => (a != a && !(b != b)) || a.lt b
    let equal := fun (a b : Float) => a == b || (a != a && b != b)
    pure (GoLean.Slices.binarySearch less equal 0 xs key)
  return s!"{index} {found}"
