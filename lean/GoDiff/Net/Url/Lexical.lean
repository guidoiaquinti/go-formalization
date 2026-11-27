import GoLean.Net.Url.Lexical
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Url

private def wire (raw : URLBytes) :=
  s!"{raw.length}:{String.intercalate "," (raw.map (toString ∘ BitVec.toNat))}"

def main : IO Unit := runTuple fun row => do
  let tag :: raw := row | throw (IO.userError "expected receiver tag and bytes")
  let data ← raw.mapM (wordInput 8)
  let live := if tag == 0 then none else some data
  let (host, port) := splitHostPort data
  let scheme := match getScheme data with
    | none => "missing"
    | some (name, suffix) => wire name ++ "/" ++ wire suffix
  return String.intercalate "|" [toString (validOptionalPort data), wire host ++ "/" ++ wire port,
    (urlHostname live).map wire |>.getD "panic", (urlPort live).map wire |>.getD "panic",
    (urlIsAbs live).map toString |>.getD "panic", scheme, toString (stringContainsCTLByte data)]
