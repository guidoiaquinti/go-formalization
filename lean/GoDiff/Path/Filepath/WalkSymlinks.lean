import GoLean.Path.Filepath.WalkSymlinks
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def raw (path : String) : List (BitVec 8) :=
  path.toUTF8.data.toList.map UInt8.toBitVec

private def key (windows : Bool) (path : List (BitVec 8)) : List (BitVec 8) :=
  let path := cleanPath windows path
  if windows && path.take 3 == raw "C:\\" then path.drop 3 else path

private def stat (windows : Bool) (path : List (BitVec 8)) : Except String SymlinkEntry :=
  let name := key windows path
  if path.contains 0 then .error "invalid"
  else if windows && name == raw "." then .ok .symlink
  else if [raw ".", raw "dir", raw "..", raw "/base", raw "/base/dir"].contains name then .ok .directory
  else if [raw "dir/file", raw "dir\\file", raw "/base/dir/file"].contains name then .ok .file
  else if [raw "filelink", raw "link", raw "chain", raw "broken", raw "cycle", raw "up", raw "abs"].contains name then .ok .symlink
  else .error "missing"

private def readLink (windows : Bool) (path : List (BitVec 8)) : Except String (List (BitVec 8)) :=
  let name := key windows path
  match ([(".", "dir"), ("filelink", "dir/file"), ("link", "dir"), ("chain", "link"),
      ("broken", "missing"), ("cycle", "cycle"), ("up", "dir/../dir"), ("abs", "/base/dir")].find?
        fun pair => raw pair.1 == name) with
  | some entry => .ok (raw entry.2)
  | none => .error "invalid"

def main : IO Unit := runTuple fun row => do
  let platform :: fields := row | throw (IO.userError "expected profile and path")
  let windows := platform == 1
  let path ← fields.mapM (wordInput 8)
  match walkSymlinks windows 1024 path (stat windows) (readLink windows) with
  | .ok bytes => return "ok:" ++ String.intercalate "," (bytes.map (toString ∘ BitVec.toNat))
  | .error (.observed error) => return error
  | .error .notDirectory => return "notDirectory"
  | .error .tooManyLinks => return "tooManyLinks"
  | .error .exhaustedModel => return "exhaustedModel"
