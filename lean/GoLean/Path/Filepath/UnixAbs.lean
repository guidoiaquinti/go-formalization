import GoLean.Path.Filepath.IsAbs
import GoLean.Path.Filepath.Join

namespace GoLean.Path.Filepath

private def usesWindowsSeparators : AbsProfile → Bool
  | .windows => true
  | _ => false

def unixAbsPath (profile : AbsProfile) (path : List (BitVec 8))
    (workingDirectory : Except ε (List (BitVec 8))) : Except ε (List (BitVec 8)) :=
  let windows := usesWindowsSeparators profile
  if isAbs profile path then .ok (cleanPath windows path)
  else workingDirectory.map fun cwd => joinPath windows [cwd, path]

theorem unixAbsPath_absolute (profile : AbsProfile) (path : List (BitVec 8))
    (failure : ε) (absolute : isAbs profile path = true) :
    unixAbsPath profile path (.error failure) =
      .ok (cleanPath (usesWindowsSeparators profile) path) := by
  unfold unixAbsPath
  rw [absolute]
  rfl

theorem unixAbsPath_relativeError (profile : AbsProfile) (path : List (BitVec 8))
    (failure : ε) (relative : isAbs profile path = false) :
    unixAbsPath profile path (.error failure) = .error failure := by
  unfold unixAbsPath
  rw [relative]
  rfl

theorem unixAbsPath_relativeJoin (profile : AbsProfile) (path cwd : List (BitVec 8))
    (relative : isAbs profile path = false) :
    unixAbsPath (ε := Unit) profile path (.ok cwd) =
      .ok (joinPath (usesWindowsSeparators profile) [cwd, path]) := by
  unfold unixAbsPath
  rw [relative]
  rfl

end GoLean.Path.Filepath
