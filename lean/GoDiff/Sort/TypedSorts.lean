import GoLean.Sort.TypedSorts
import GoDiff.Math.Bits.Tuple

open GoLean.Sort

private def parseWords : Nat → List Int → List (List UInt8)
  | 0, _ | _, [] => []
  | fuel + 1, size :: raw =>
    ((raw.take size.toNat).map (UInt8.ofNat ∘ Int.toNat)) ::
      parseWords fuel (raw.drop size.toNat)

def main : IO Unit := runTuple fun row => do
  let kind :: width :: raw := row | throw (IO.userError "expected kind, width and data")
  let width := width.toNat
  if kind = 0 then
    return "values:" ++ String.intercalate "," ((sortInts width raw.toArray).toList.map toString)
  else if kind = 1 || kind = 2 then
    let floats ← raw.mapM (wordInput 64)
    let sorted := if kind = 1 then sortFloat64s width floats.toArray
      else float64SliceSort width floats.toArray
    return "values:" ++ String.intercalate "," (sorted.toList.map (toString ∘ BitVec.toNat))
  else if kind = 5 then
    let word ← wordInput 64 (raw.headD 0)
    return toString (sortIsNaN word)
  else
    let words := (parseWords raw.length raw).toArray
    let sorted := if kind = 3 then sortStrings width words else stringSliceSort width words
    let showWord := fun bs : List UInt8 => s!"{bs.length}:" ++
      String.intercalate "," (bs.map (toString ∘ UInt8.toNat))
    return "values:" ++ String.intercalate ";" (sorted.toList.map showWord)
