import GoLean.Strings.TrimRightASCII
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let count :: bytes := fields | throw (IO.userError "expected the table size and bytes")
  let raw ← bytes.mapM (wordInput 8)
  let members := raw.take count.toNat
  let table := if count < 0 then none else some (fun byte => members.contains byte)
  match GoLean.Strings.trimRightASCII (raw.drop count.toNat) table with
  | none => return "panic"
  | some trimmed =>
    return "bytes:" ++ String.intercalate "," (trimmed.map (fun byte => toString byte.toNat))
