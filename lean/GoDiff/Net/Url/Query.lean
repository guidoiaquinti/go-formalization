import GoLean.Net.Url.Query
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Url

private def wire (raw : URLBytes) :=
  s!"{raw.length}:{String.intercalate "," (raw.map (toString ∘ BitVec.toNat))}"

private def state : Values → String
  | none => "nil"
  | some entries =>
    let showEntry (entry : ValuesEntry) :=
      wire entry.1 ++ "=" ++ (entry.2.map (fun words =>
        "list:" ++ String.intercalate ";" (words.map wire))).getD "nil"
    "map:" ++ String.intercalate "/" ((sortValuesEntries entries).map showEntry)

private def report (v : Values × Option ParseQueryError) : String :=
  let status := match v.2 with
    | none => "ok"
    | some .limit => "limit"
    | some .semicolon => "semicolon"
    | some (.unescape (.escape bytes)) => "escape:" ++ wire bytes
    | some (.unescape (.host byte)) => "host:" ++ wire [byte]
  status ++ "/" ++ state v.1

def main : IO Unit := runTuple fun row => do
  let width :: params :: start :: present :: configLen :: raw := row |
    throw (IO.userError "expected width, count, map, receiver and config")
  let data ← raw.mapM (wordInput 8)
  let setting := data.take configLen.toNat
  let query := data.drop configLen.toNat
  let v : Values := if start == 0 then none else if start == 1 then some []
    else some [([97], some [[111, 108, 100]]), ([122], none)]
  let (allowed, changed) := urlParamsWithinMax params setting width.toNat
  let native := parseQueryInto v query setting width.toNat
  let parsed := parseQuery query setting width.toNat
  let url := urlQuery (if present == 0 then none else some query) setting width.toNat
  return String.intercalate "|" [s!"{allowed}:{changed}", native.map report |>.getD "panic",
    report parsed, url.map state |>.getD "panic"]
