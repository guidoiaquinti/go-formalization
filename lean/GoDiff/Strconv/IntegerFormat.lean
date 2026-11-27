import GoLean.Strconv.IntegerFormat
import GoDiff.Math.Bits.Tuple

open GoLean.Strconv

private def wire (values : Option (List Nat)) : String :=
  (values.map fun raw => "slice:" ++ String.intercalate "," (raw.map toString)).getD "nil"

private def formatted (values : Option (List Nat)) : String :=
  if values.isNone then "panic" else wire values

private def appended (out : FormatAppend) : String :=
  s!"{out.panicked}/{wire out.output}/{wire out.original}/{out.reused}"

def main : IO Unit := runTuple fun row => do
  let width :: unsigned :: base :: length :: present :: raw := row |
    throw (IO.userError "expected word width, word, radix and buffer fields")
  let word ← wordInput 64 unsigned
  let cells ← raw.mapM fun v => BitVec.toNat <$> wordInput 8 v
  let dst : FormatBuffer := ⟨if present = 0 then none else some cells, length.toNat⟩
  let internal := [formatted (integerFormatUint word base), formatted (integerFormatInt word base),
    formatted (integerItoa width.toNat word.toInt), appended (integerAppendUint dst word base),
    appended (integerAppendInt dst word base)]
  let publicResults := [formatted (publicFormatUint word base), formatted (publicFormatInt word base),
    formatted (publicItoa width.toNat word.toInt), appended (publicAppendUint dst word base),
    appended (publicAppendInt dst word base)]
  return String.intercalate "|" (internal ++ publicResults)
