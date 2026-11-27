import GoLean.Net.Url.Userinfo
import GoDiff.Math.Bits.Tuple

open GoLean.Net.Url

private def wire (raw : URLBytes) :=
  s!"{raw.length}:{String.intercalate "," (raw.map (toString ∘ BitVec.toNat))}"

private def joined (password : URLBytes) (set : Bool) : String :=
  wire password ++ ":" ++ toString set

private def state (u : Userinfo) : String :=
  wire u.username ++ "/" ++ joined u.password u.passwordSet

def main : IO Unit := runTuple fun row => do
  let tag :: nameLen :: passwordLen :: raw := row | throw (IO.userError "expected user fields")
  let bytes ← raw.mapM (wordInput 8)
  let name := bytes.take nameLen.toNat
  let password := (bytes.drop nameLen.toNat).take passwordLen.toNat
  let validate := bytes.drop (nameLen.toNat + passwordLen.toNat)
  let receiver : Option Userinfo := if tag == 0 then none else some ⟨name, password, tag == 2⟩
  let (pw, set) := userinfoPassword receiver
  return String.intercalate "|" [state (user name), state (userPassword name password),
    wire (userinfoUsername receiver), joined pw set, wire (userinfoString receiver),
    toString (validUserinfo validate)]
