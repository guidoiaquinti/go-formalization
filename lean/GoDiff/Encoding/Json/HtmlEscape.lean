import GoLean.Encoding.Json.HtmlEscape
import GoDiff.Math.Bits.Tuple

open GoLean.Encoding.Json

private def wire (data : Option (List Nat)) : String :=
  (data.map fun b => "slice:" ++ String.intercalate "," (b.map toString)).getD "nil"

def main : IO Unit := runTuple fun row => do
  let present :: dst :: length :: cap :: raw := row | throw (IO.userError "expected buffer and slice state")
  let words ← raw.mapM (wordInput 8)
  let data := words.map BitVec.toNat
  let backing := if dst = 0 then none else some (data.take cap.toNat)
  let src := data.drop cap.toNat
  let result := appendHTMLEscape backing length.toNat src
  let receiver := if present = 0 then none else some ((backing.getD []).take length.toNat)
  return s!"{wire result.output}/{wire result.original}/{result.reused}|{(htmlEscape receiver src).map (fun values => wire (some values)) |>.getD "panic"}"
