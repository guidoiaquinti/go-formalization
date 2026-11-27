import GoLean.Strings.ReplacerReplace
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

private def wire (output : List (BitVec 8)) : String :=
  String.intercalate "," (output.map fun byte => toString byte.toNat)

def main : IO Unit := runTuple fun row => do
  let tag :: count :: fields := row | throw (IO.userError "expected receiver and pair count")
  let (pairs, remaining) ← parsePairs count.toNat fields
  let firstWidth :: remaining := remaining | throw (IO.userError "expected input byte count")
  let (first, rest) ← takeBytes firstWidth remaining
  let second ← rest.mapM (wordInput 8)
  let receiver := if tag == 0 then none else some (ReplacerState.pending pairs)
  let some (state, first) := replacerReplace receiver first | return "panic"
  let some (_, second) := replacerReplace (some state) second | return "panic"
  return wire first ++ "|" ++ wire second
