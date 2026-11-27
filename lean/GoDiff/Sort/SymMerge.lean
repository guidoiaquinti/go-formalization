import GoLean.Sort.SymMerge
import GoDiff.Math.Bits.Tuple

open GoLean.Sort

def main : IO Unit := runTuple fun row => do
  let variant :: a :: m :: b :: raw := row | throw (IO.userError "expected variant, bounds and keys")
  let data : MergeState := ⟨raw.zipIdx.toArray, []⟩
  let data := if variant = 0 then symMerge data a.toNat m.toNat b.toNat
    else symMergeFunc data a.toNat m.toNat b.toNat
  let values := String.intercalate "," (data.values.toList.map (toString ∘ Prod.snd))
  let calls := data.reversedCalls.reverse.map fun (swap, i, j) =>
    s!"{if swap then "S" else "L"}{i},{j}"
  return values ++ ":" ++ String.intercalate ";" calls
