import GoLean.Path.Filepath.EvalSymlinks
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def readStrings : Nat → List Int → IO (List (List (BitVec 8)))
  | 0, [] => pure []
  | count + 1, length :: fields => do
    if length < 0 || length.toNat > fields.length then throw (IO.userError "invalid field length")
    let str ← (fields.take length.toNat).mapM (wordInput 8)
    return str :: (← readStrings count (fields.drop length.toNat))
  | _, _ => throw (IO.userError "invalid fields")

def main : IO Unit := runTuple fun row => do
  let platform :: walkError :: normError :: statError :: raw := row |
    throw (IO.userError "expected profile and delegate outcomes")
  let [path, walked, normalized, message] ← readStrings 4 raw |
    throw (IO.userError "expected path and response fields")
  let outcome := fun (flag : Int) bytes =>
    if flag == 0 then Except.ok bytes else Except.error (EvalSymlinkError.observed flag bytes)
  let profile := if platform == 1 then AbsProfile.windows else if platform == 2 then AbsProfile.plan9
    else AbsProfile.unix
  let result := evalSymlinks profile path (fun _ => outcome walkError walked)
    (fun query => if query == walked then outcome normError normalized
      else .error (.observed (-1 : Int) []))
    (fun _ => if statError == 0 then none else some (statError, message))
  let render := fun (bytes : List (BitVec 8)) => String.intercalate "," (bytes.map (toString ∘ BitVec.toNat))
  match result with
  | .ok bytes => return "ok:" ++ render bytes
  | .error .notDirectory => return "notDirectory"
  | .error (.observed identity message) => return s!"err:{identity}:" ++ render message
