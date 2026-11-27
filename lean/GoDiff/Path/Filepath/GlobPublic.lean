import GoLean.Path.Filepath.GlobWithLimit
import GoDiff.Math.Bits.Tuple

open GoLean.Path.Filepath

private def raw (value : String) : List (BitVec 8) :=
  value.toUTF8.data.toList.map UInt8.toBitVec

private def files : List (String × List String) :=
  [(".", ["alpha", "beta", "plain"]), ("alpha", ["a.txt", "z.go", "tree"]), ("beta", ["b.txt"]), ("alpha/tree", ["deep.txt"])]

private def entryNames (windows : Bool) (path : List (BitVec 8)) : Option (List (List (BitVec 8))) :=
  (files.find? fun entry => fromSlash (if windows then 92 else 47) (raw entry.1) == cleanPath windows path).map fun entry => entry.2.map raw

private def traversable (windows : Bool) (prefixRev : List (BitVec 8)) :
    List (BitVec 8) → Bool
  | [] => true
  | byte :: rest =>
    if (byte == 47 || (windows && byte == 92)) &&
        (entryNames windows prefixRev.reverse).isNone then false
    else traversable windows (byte :: prefixRev) rest

private def lstat (windows : Bool) (path : List (BitVec 8)) : Bool :=
  !path.isEmpty && traversable windows [] path &&
    ((entryNames windows path).isSome ||
    (!(path.getLast?.any fun byte => byte == 47 || (windows && byte == 92)) &&
      files.any fun entry =>
        entry.2.any fun name => joinPath windows [raw entry.1, raw name] == cleanPath windows path))

def main : IO Unit := runTuple fun fields => do
  let profile :: depth :: fields := fields | throw (IO.userError "expected profile, depth and pattern")
  let windows := profile == 1
  let pattern ← fields.mapM (wordInput 8)
  let directory := fun path =>
    match if traversable windows [] path then entryNames windows path else none with
    | some names => GlobDirectory.opened names
    | none => GlobDirectory.statFailed
  let (result, invalid) := globWithLimit windows pattern depth (lstat windows) directory
  if invalid then return "bad"
  match result with
  | none => return "nil"
  | some paths =>
    let render := fun (path : List (BitVec 8)) =>
      String.intercalate "," (path.map (toString ∘ BitVec.toNat))
    return s!"{paths.length}:{String.intercalate ";" (paths.map render)}"
