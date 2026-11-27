import GoLean.Bytes.TrimLeftASCII
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let count :: spare :: bytes := fields |
    throw (IO.userError "expected table size, spare capacity, and bytes")
  let raw ← bytes.mapM (wordInput 8)
  let members := raw.take count.toNat
  let table := if count < 0 then none else some (fun byte => members.contains byte)
  match GoLean.Bytes.trimLeftASCII (raw.drop count.toNat) table spare.toNat with
  | none => return "panic"
  | some none => return "nil"
  | some (some (trimmed, capacity)) =>
    return toString capacity ++ ":" ++ String.intercalate "," (trimmed.map (fun byte => toString byte.toNat))
