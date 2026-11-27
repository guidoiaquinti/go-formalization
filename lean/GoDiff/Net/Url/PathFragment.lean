import GoLean.Net.Url.PathFragment
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Url

private def bytes (xs : List (BitVec 8)) : String :=
  String.intercalate "," (xs.map (toString ∘ BitVec.toNat))

private def setter : Option (Except UnescapeError Unit × Option URLPart) → String
  | none => "panic"
  | some (result, part) =>
    let status := match result with
      | .ok _ => "ok"
      | .error (.escape raw) => "escape:" ++ bytes raw
      | .error (.host b) => s!"host:{b.toNat}"
    let fields := part.map (fun p => bytes p.decoded ++ "/" ++ bytes p.raw) |>.getD "nil"
    status ++ "/" ++ fields

def main : IO Unit := runTuple fun row => do
  let tag :: mode :: size :: rawSize :: fields := row | throw (IO.userError "expected URL fields")
  let mode ← wordInput 8 mode
  let data ← fields.mapM (wordInput 8)
  let raw := (data.drop size.toNat).take rawSize.toNat
  let update := data.drop (size.toNat + rawSize.toNat)
  let prior : Option URLPart :=
    if tag = 0 then none else some ⟨data.take size.toNat, raw⟩
  let escaped (x : Option (List (BitVec 8))) := x.map (fun xs => "ok:" ++ bytes xs) |>.getD "panic"
  return String.intercalate "|" [toString (validEncoded raw mode), escaped (escapedPath prior),
    escaped (escapedFragment prior), setter (setPath prior update), setter (setFragment prior update)]
