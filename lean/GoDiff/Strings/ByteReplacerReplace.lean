import GoLean.Strings.ByteReplacerReplace
import GoDiff.Math.Bits.Tuple

def main : IO Unit := runTuple fun fields => do
  let mode :: fields := fields | throw (IO.userError "expected table marker and bytes")
  let raw ← fields.mapM (wordInput 8)
  if mode == 1 && raw.length < 256 then
    throw (IO.userError "expected the full byte table")
  let table := (raw.take 256).toArray
  let lookup := if mode == 1 then some (fun (byte : BitVec 8) => table[byte.toNat]!) else none
  let input := if mode == 1 then raw.drop 256 else raw
  let some result := GoLean.Strings.byteReplacerReplace lookup input | return "panic"
  return "bytes:" ++ String.intercalate "," (result.map (fun byte => toString byte.toNat))
