import GoLean.Encoding.Binary.CoderValue
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Binary

private def fieldInfo (kind : Int) : ValueItem × Nat :=
  if kind == 13 then (.ignored, 0)
  else
    let width := if kind ≤ 2 then 1 else if kind ≤ 4 then 2
      else if kind ≤ 6 || kind == 9 || kind == 11 then 4 else 8
    let parts := if kind == 11 || kind == 12 then 2 else 1
    (.field width parts (kind == 0), parts)

private def wire (values : List Nat) : String :=
  String.intercalate "," (values.map toString)

def main : IO Unit := runTuple fun row => do
  let kind :: shape :: count :: order :: offset :: length :: raw := row |
    throw (IO.userError "expected kind, shape, count, order, offset, length and cells")
  let (field, wordsPerField) := fieldInfo kind
  let n := if shape == 0 then 1 else if shape == 3 then count.toNat + 1 else count.toNat
  let plan := if shape == 3 then field :: .blank 2 :: List.replicate count.toNat field
    else List.replicate n field
  let wordCount := n * wordsPerField
  let words := (raw.take wordCount).map Int.toNat
  let cells ← (raw.drop wordCount).mapM (wordInput 8)
  let (d, read, values) := decoderValue (order != 0) length.toNat cells plan offset.toNat words
  let (e, written, backing) := encoderValue (order != 0) length.toNat plan offset.toNat words cells
  return s!"{d}:{read}:{wire values}|{e}:{written}:{wire (backing.map BitVec.toNat)}"
