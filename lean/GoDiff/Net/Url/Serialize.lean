import GoLean.Net.Url.Serialize
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Url

private def wire (raw : URLBytes) := String.intercalate "," (raw.map (toString ∘ BitVec.toNat))

private def splitParts (sizes : List Int) (cells : URLBytes) : List URLBytes :=
  match sizes with
  | [] => [cells]
  | n :: rest => cells.take n.toNat :: splitParts rest (cells.drop n.toNat)

def main : IO Unit := runTuple fun row => do
  let present :: flags :: kind :: length :: fields := row |
    throw (IO.userError "expected URL, append and flags")
  let sizes := fields.take 10
  let cells ← (fields.drop 10).mapM (wordInput 8)
  let [scheme, opaquePart, user, password, host, path, rawPath, query, fragment, rawFragment, backing] :=
    splitParts sizes cells | throw (IO.userError "expected ten URL field widths")
  let credentials := if kind == 0 then none else some ⟨user, password, kind == 2⟩
  let u : URLRender := ⟨scheme, opaquePart, credentials, host, ⟨path, rawPath⟩,
    flags.toNat % 2 == 1, flags.toNat / 2 == 1, query, ⟨fragment, rawFragment⟩⟩
  let receiver := if present == 1 then some u else none
  let optionalWire (out : Option URLBytes) := out.map wire |>.getD "nil"
  let appended :=
    match urlAppendBinary receiver backing length.toNat with
    | none => "nil:" ++ wire backing
    | some a => s!"{wire a.output}:{wire a.original}:{a.reused}"
  return String.intercalate "|" [optionalWire (urlString receiver), optionalWire (urlRequestURI receiver),
    wire (urlRedacted receiver), optionalWire (urlMarshalBinary receiver), appended]
