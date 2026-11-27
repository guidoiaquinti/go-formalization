import GoLean.Strings.StringWriterWriteString
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let bounded :: count :: error :: fields := fields
    | throw (IO.userError "expected response fields and input bytes")
  let input ← fields.mapM (wordInput 8)
  let write := fun (payload : List (BitVec 8)) =>
    (if bounded == 1 then min count (Int.ofNat payload.length) else count, error.toNat)
  let result := GoLean.Strings.stringWriterWriteString write input
  let calls := result.writes.map fun payload =>
    "B" ++ String.join (payload.map (fun byte => "," ++ toString byte.toNat))
  return s!"{result.count}:{result.error}:" ++ String.intercalate ";" calls
