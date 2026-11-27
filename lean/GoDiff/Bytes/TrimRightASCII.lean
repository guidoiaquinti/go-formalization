import GoLean.Bytes.TrimRightASCII
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let count :: capacity :: bytes := fields |
    throw (IO.userError "expected table size, capacity or nil marker, and bytes")
  let raw ← bytes.mapM (wordInput 8)
  let members := raw.take count.toNat
  let table := if count < 0 then none else some (fun byte => members.contains byte)
  let input := if capacity < 0 then none else some (raw.drop count.toNat, capacity.toNat)
  match GoLean.Bytes.trimRightASCII input table with
  | none => return "panic"
  | some none => return "nil"
  | some (some (trimmed, capacity)) =>
    return toString capacity ++ ":" ++ String.intercalate "," (trimmed.map (fun byte => toString byte.toNat))
