import GoLean.Path.Filepath.Clean
import GoLean.Path.Filepath.IsAbs

namespace GoLean.Path.Filepath

inductive EvalSymlinkError (ε : Type) where
  | observed (identity : ε) (message : List (BitVec 8))
  | notDirectory

abbrev EvalSymlinkResult (ε : Type) :=
  Except (EvalSymlinkError ε) (List (BitVec 8))

private def containsWord (needle : List (BitVec 8)) : List (BitVec 8) → Bool
  | [] => needle.isEmpty
  | byte :: rest =>
    (byte :: rest).take needle.length == needle || containsWord needle rest

def evalSymlinks (profile : AbsProfile) (path : List (BitVec 8))
    (walk : List (BitVec 8) → EvalSymlinkResult ε)
    (normalize : List (BitVec 8) → EvalSymlinkResult ε)
    (stat : List (BitVec 8) → Option (ε × List (BitVec 8))) :
    EvalSymlinkResult ε :=
  match profile with
  | .unix => walk path
  | .windows => (walk path).bind normalize
  | .plan9 =>
    let failure := if path.isEmpty then none else stat path
    match failure with
    | none => .ok (cleanPath false path)
    | some (identity, message) =>
      if containsWord ("not a directory".toUTF8.data.toList.map UInt8.toBitVec) message
      then .error .notDirectory
      else .error (.observed identity message)

theorem evalSymlinks_unix (path : List (BitVec 8))
    (walk normalize : List (BitVec 8) → EvalSymlinkResult ε)
    (stat : List (BitVec 8) → Option (ε × List (BitVec 8))) :
    evalSymlinks .unix path walk normalize stat = walk path := rfl

theorem evalSymlinks_windowsWalkError (path : List (BitVec 8))
    (walk normalize : List (BitVec 8) → EvalSymlinkResult ε)
    (stat : List (BitVec 8) → Option (ε × List (BitVec 8)))
    (failure : EvalSymlinkError ε) (observed : walk path = .error failure) :
    evalSymlinks .windows path walk normalize stat = .error failure := by
  change (walk path).bind normalize = .error failure
  rw [observed]
  rfl

theorem evalSymlinks_plan9Empty
    (walk normalize : List (BitVec 8) → EvalSymlinkResult ε)
    (stat : List (BitVec 8) → Option (ε × List (BitVec 8))) :
    evalSymlinks .plan9 [] walk normalize stat = .ok [46] := rfl

end GoLean.Path.Filepath
