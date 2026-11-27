import GoLean.Index.Suffixarray.Lookup
import GoDiff.Math.Bits.Tuple

open GoLean.Index.Suffixarray

private def wire (v : View (BitVec width)) : String :=
  match v.backing with
  | none => "nil"
  | some words => s!"{v.length}:{String.intercalate "," (words.map (toString ∘ BitVec.toInt))}"

def main : IO Unit := runTuple fun fields => do
  let nativeWidth :: profile :: size :: needleSize :: limit :: cells := fields |
    throw (IO.userError "expected widths, lengths and backing")
  let text ← (cells.take size.toNat).mapM (wordInput 8)
  let needle ← ((cells.drop size.toNat).take needleSize.toNat).mapM (wordInput 8)
  let offsets := cells.drop (size + needleSize).toNat
  let sa : Ints := match profile with
    | 1 => ⟨⟨some (offsets.map (BitVec.ofInt 32)), size.toNat⟩, ⟨none, 0⟩⟩
    | 2 => ⟨⟨none, 0⟩, ⟨some (offsets.map (BitVec.ofInt 64)), size.toNat⟩⟩
    | _ => ⟨⟨none, 0⟩, ⟨none, 0⟩⟩
  let s := if profile = 0 then none else some (⟨⟨some text, size.toNat⟩, sa⟩ : IndexState)
  let all := match lookupAll s needle with
    | .ok located => wire located.words32 ++ "/" ++ wire located.words64
    | .error _ => "panic"
  let result := match indexLookup nativeWidth.toNat s needle limit with
    | .ok none => "nil"
    | .ok (some offsets) => s!"{offsets.length}:{String.intercalate "," (offsets.map toString)}"
    | .error _ => "panic"
  return all ++ "|" ++ result
