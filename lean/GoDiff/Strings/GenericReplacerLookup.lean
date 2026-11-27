import GoLean.Strings.GenericReplacerLookup
import GoDiff.Math.Bits.Tuple

private def takeBytes (length : Int) (fields : List Int) :
    IO (List (BitVec 8) × List Int) := do
  if length < 0 || length.toNat > fields.length then
    throw (IO.userError "invalid byte count")
  let bytes ← (fields.take length.toNat).mapM (wordInput 8)
  return (bytes, fields.drop length.toNat)

private def parsePairs : Nat → List Int →
    IO (List (List (BitVec 8) × List (BitVec 8)) × List Int)
  | 0, fields => pure ([], fields)
  | count + 1, keyLength :: valueLength :: fields => do
    let (key, fields) ← takeBytes keyLength fields
    let (value, fields) ← takeBytes valueLength fields
    let (pairs, fields) ← parsePairs count fields
    return ((key, value) :: pairs, fields)
  | _, _ => throw (IO.userError "missing key and value lengths")

def main : IO Unit := runTuple fun row => do
  let flag :: count :: fields := row | throw (IO.userError "expected root flag and pair count")
  let (pairs, remaining) ← parsePairs count.toNat fields
  let input ← remaining.mapM (wordInput 8)
  match GoLean.Strings.genericReplacerLookup pairs input (flag != 0) with
  | none => return "miss:0:"
  | some result =>
    let bytes := String.intercalate "," (result.value.map fun byte => toString byte.toNat)
    return s!"hit:{result.keyLength}:{bytes}"
