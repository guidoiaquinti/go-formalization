import GoLean.Net.Url.Values
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Url

private def readBytes : List Int → IO (URLBytes × List Int)
  | n :: raw => do
    let bytes ← (raw.take n.toNat).mapM (wordInput 8)
    return (bytes, raw.drop n.toNat)
  | _ => throw (IO.userError "expected a string length")

private def readWords : Nat → List Int → IO (List URLBytes × List Int)
  | 0, raw => pure ([], raw)
  | n + 1, raw => do
    let (word, raw) ← readBytes raw
    let (words, raw) ← readWords n raw
    return (word :: words, raw)

private def readEntries : Nat → List Int → IO (List ValuesEntry × List Int)
  | 0, raw => pure ([], raw)
  | n + 1, raw => do
    let (key, raw) ← readBytes raw
    let size :: raw := raw | throw (IO.userError "expected entry size")
    let (words, raw) ← readWords size.toNat raw
    let (entries, raw) ← readEntries n raw
    return ((key, if size < 0 then none else some words) :: entries, raw)

private def wire (raw : URLBytes) :=
  s!"{raw.length}:{String.intercalate "," (raw.map (toString ∘ BitVec.toNat))}"

private def state : Values → String
  | none => "nil"
  | some entries =>
    let showEntry (entry : ValuesEntry) :=
      wire entry.1 ++ "=" ++ (entry.2.map (fun words =>
        "list:" ++ String.intercalate ";" (words.map wire))).getD "nil"
    "map:" ++ String.intercalate "/" ((sortValuesEntries entries).map showEntry)

def main : IO Unit := runTuple fun row => do
  let tag :: size :: raw := row | throw (IO.userError "expected map tag and size")
  let (key, raw) ← readBytes raw
  let (word, raw) ← readBytes raw
  let (entries, _) ← readEntries size.toNat raw
  let v := if tag = 0 then none else some entries
  return String.intercalate "|" [wire (valuesGet v key), toString (valuesHas v key),
    (valuesSet v key word).map (fun entries => state (some entries)) |>.getD "panic",
    (valuesAdd v key word).map (fun entries => state (some entries)) |>.getD "panic",
    state (valuesDel v key), state (valuesClone v), wire (valuesEncode v)]
