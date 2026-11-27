import GoLean.Net.Url.Clone
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Url

private def wire (b : URLBytes) : String :=
  String.intercalate "," (b.map (toString ∘ BitVec.toNat))

private def unpack : List Int → List (BitVec 8) → List URLBytes
  | [], _ => []
  | size :: rest, raw => raw.take size.toNat :: unpack rest (raw.drop size.toNat)

def main : IO Unit := runTuple fun row => do
  let presence :: userTag :: force :: omitFlag :: fields := row |
    throw (IO.userError "expected URL profile and fields")
  let bytes ← (fields.drop 10).mapM (wordInput 8)
  let [scheme, opaqueValue, host, path, fragment, query, rawPath, rawFragment, username, password] :=
    unpack (fields.take 10) bytes | throw (IO.userError "expected ten URL strings")
  let user := if userTag == 0 then none else some (2, ⟨username, password, userTag == 2⟩)
  let original := ⟨⟨scheme, opaqueValue, host, path, fragment, query, rawPath, rawFragment,
    force == 1, omitFlag == 1⟩, 1, user⟩
  match urlClone (if presence == 0 then none else some original) 3 4 with
  | none => return "nil"
  | some result =>
    let f := result.fields
    let username := result.user.map (fun u => u.2.username) |>.getD []
    let password := result.user.map (fun u => u.2.password) |>.getD []
    let userTag := result.user.map (fun u => if u.2.passwordSet then 2 else 1) |>.getD 0
    let userFresh := result.user.all (fun u => original.user.all (fun old => u.1 != old.1))
    return String.intercalate "|" (([f.scheme, f.opaqueValue, f.host, f.path, f.fragment,
      f.rawQuery, f.rawPath, f.rawFragment, username, password].map wire) ++
      [toString f.forceQuery, toString f.omitHost, toString userTag,
        toString (result.identity != original.identity), toString userFresh])
