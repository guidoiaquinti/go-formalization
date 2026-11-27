import GoLean.Bytes.GrowSlice
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun row => do
  let cap :: count :: maxInt :: raw := row | throw (IO.userError "expected capacity, count and source bytes")
  let input ← raw.mapM (wordInput 8)
  match GoLean.Bytes.growSlice input cap.toNat count maxInt.toNat with
  | none => return "nil"
  | some (payload, spare) =>
    let wire := fun (xs : List (BitVec 8)) => String.intercalate "," (xs.map (toString ∘ BitVec.toNat))
    return s!"{payload.length + spare.length}:{wire payload}:{wire spare}"
