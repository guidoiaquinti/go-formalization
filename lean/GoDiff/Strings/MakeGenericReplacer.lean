import GoLean.Strings.MakeGenericReplacer
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

private def terminalWire : GenericTerminal → String
  | none => "0:"
  | some (priority, value) => s!"{priority}:{wire value}"

private def trieWire : Nat → GenericTrie → String
  | 0, _ => "!"
  | fuel + 1, trie =>
    match trie with
    | .leaf terminal => "L" ++ terminalWire terminal
    | .edge terminal prefixBytes child =>
      "E" ++ terminalWire terminal ++ "[" ++ wire prefixBytes ++ "|" ++ trieWire fuel child ++ "]"
    | .table terminal children =>
      "T" ++ terminalWire terminal ++ "{" ++
        String.intercalate "|" (children.map fun (byte, child) =>
          toString byte.toNat ++ "=" ++ trieWire fuel child) ++ "}"

def main : IO Unit := runTuple fun row => do
  let count :: fields := row | throw (IO.userError "expected pair count")
  let (pairs, remaining) ← parsePairs count.toNat fields
  if !remaining.isEmpty then throw (IO.userError "unexpected extra bytes")
  let result := makeGenericReplacer pairs
  let depth := pairs.foldl (fun n pair => max n pair.1.length) 0
  return s!"{result.tableSize}:{wire result.mapping}|" ++ trieWire (depth + 1) result.root
