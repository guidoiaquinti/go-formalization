import GoLean.Path.Filepath.UnixAbs

namespace GoLean.Path.Filepath

def absPath (profile : AbsProfile) (path : List (BitVec 8))
    (workingDirectory : Except ε (List (BitVec 8)))
    (fullPath : List (BitVec 8) → Except ε (List (BitVec 8))) : Except ε (List (BitVec 8)) :=
  match profile with
  | .windows =>
    let query := if path.isEmpty then [46] else path
    (fullPath query).map (cleanPath true)
  | profile => unixAbsPath profile path workingDirectory

theorem absPath_unix (path : List (BitVec 8)) (cwd : Except ε (List (BitVec 8)))
    (fullPath : List (BitVec 8) → Except ε (List (BitVec 8))) :
    absPath .unix path cwd fullPath = unixAbsPath .unix path cwd := rfl

theorem absPath_plan9 (path : List (BitVec 8)) (cwd : Except ε (List (BitVec 8)))
    (fullPath : List (BitVec 8) → Except ε (List (BitVec 8))) :
    absPath .plan9 path cwd fullPath = unixAbsPath .plan9 path cwd := rfl

theorem absPath_windowsEmpty (cwd : Except ε (List (BitVec 8)))
    (fullPath : List (BitVec 8) → Except ε (List (BitVec 8))) :
    absPath .windows [] cwd fullPath = (fullPath [46]).map (cleanPath true) := rfl

theorem absPath_windowsError (path : List (BitVec 8)) (cwd : Except ε (List (BitVec 8)))
    (fullPath : List (BitVec 8) → Except ε (List (BitVec 8))) (failure : ε)
    (failed : fullPath (if path.isEmpty then [46] else path) = .error failure) :
    absPath .windows path cwd fullPath = .error failure := by
  change (fullPath (if path.isEmpty then [46] else path)).map (cleanPath true) = .error failure
  rw [failed]
  rfl

end GoLean.Path.Filepath
