import GoLean.Strings.ReplacerBuild
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

open GoLean.Strings

private def wire (bytes : List (BitVec 8)) : String :=
  String.intercalate "," (bytes.map fun byte => toString byte.toNat)

def main : IO Unit := runTuple fun row => do
  let count :: fields := row | throw (IO.userError "expected pair count")
  let (pairs, remaining) ← parsePairs count.toNat fields
  if !remaining.isEmpty then throw (IO.userError "unexpected extra bytes")
  let alphabet := (List.range 256).map (BitVec.ofNat 8)
  match replacerBuild pairs with
  | .single state => return "single:" ++ wire state.finder.pattern ++ "|" ++ wire state.value
  | .generic _ state => return s!"generic:{state.tableSize}:" ++ wire state.mapping
  | .bytes lookup => return "bytes:" ++ wire (alphabet.map lookup)
  | .byteStrings state =>
    let values := alphabet.map fun byte =>
      match state.lookup byte with
      | none => "-"
      | some value => "+" ++ wire value
    return "byteStrings:" ++ wire state.keys ++ "|" ++ String.intercalate "|" values
