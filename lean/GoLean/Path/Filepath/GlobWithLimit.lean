import GoLean.Path.Filepath.CleanGlobPath
import GoLean.Path.Filepath.CleanGlobPathWindows
import GoLean.Path.Filepath.Glob
import GoLean.Path.Filepath.HasMeta
import GoLean.Path.Filepath.Split

namespace GoLean.Path.Filepath

abbrev GlobResult := Option (List (List (BitVec 8))) × Bool

private def descendGlob (windows : Bool)
    (lstat : List (BitVec 8) → Bool) (directory : List (BitVec 8) → GlobDirectory) :
    Nat → Int → List (BitVec 8) → GlobResult
  | 0, _, _ => (none, true)
  | fuel + 1, depth, pattern =>
    if depth == 10000 || (matchPath windows pattern []).2 then (none, true)
    else if !hasMeta windows pattern then
      (if lstat pattern then some [pattern] else none, false)
    else
      let (rawDirectory, file) := splitPath windows pattern
      let (volumeLength, dir) :=
        if windows then cleanGlobPathWindows true rawDirectory
        else (0, cleanGlobPath 47 rawDirectory)
      if !hasMeta windows (dir.drop volumeLength) then
        globDirectory windows dir file none (directory dir)
      else if dir == pattern then (none, true)
      else
        let (candidates, invalid) := descendGlob windows lstat directory fuel (depth + 1) dir
        if invalid then (none, true)
        else (candidates.getD []).foldl
          (fun result candidate =>
            if result.2 then result else globDirectory windows candidate file result.1 (directory candidate))
          (none, false)

def globWithLimit (windows : Bool) (pattern : List (BitVec 8)) (depth : Int)
    (lstat : List (BitVec 8) → Bool) (directory : List (BitVec 8) → GlobDirectory) :
    GlobResult :=
  descendGlob windows lstat directory (pattern.length + 1) depth pattern

theorem globWithLimit_depthLimit (windows : Bool) (pattern : List (BitVec 8))
    (lstat : List (BitVec 8) → Bool) (directory : List (BitVec 8) → GlobDirectory) :
    globWithLimit windows pattern 10000 lstat directory = (none, true) := rfl

theorem globWithLimit_existingLiteral (windows : Bool) (pattern : List (BitVec 8))
    (depth : Int) (lstat : List (BitVec 8) → Bool) (directory : List (BitVec 8) → GlobDirectory)
    (within : (depth == 10000) = false) (valid : (matchPath windows pattern []).2 = false)
    (literal : hasMeta windows pattern = false) (present : lstat pattern = true) :
    globWithLimit windows pattern depth lstat directory = (some [pattern], false) := by
  simp [globWithLimit, descendGlob, within, valid, literal, present]

end GoLean.Path.Filepath
