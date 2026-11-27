import GoLean.Path.Filepath.Clean
import GoLean.Path.Filepath.Split

namespace GoLean.Path.Filepath

def dirName (windows : Bool) (path : List (BitVec 8)) : List (BitVec 8) :=
  let volume := volumeName windows path
  let rawDir := (splitPath windows path).1.drop volume.length
  let cleaned := cleanPath windows rawDir
  if cleaned == [46] && volume.length > 2 then volume else volume ++ cleaned

theorem dirName_unix (path : List (BitVec 8)) :
    dirName false path = cleanPath false (splitPath false path).1 := by
  simp [dirName, volumeName_unix]

theorem dirName_empty (windows : Bool) : dirName windows [] = [46] := by
  cases windows <;> decide

theorem dirName_windowsDrive :
    dirName true [67, 58, 120] = [67, 58, 46] ∧
    dirName true [67, 58, 47, 120] = [67, 58, 92] := by decide

theorem dirName_uncRoot :
    dirName true [47, 47, 104, 47, 115] = [92, 92, 104, 92, 115] := by decide

end GoLean.Path.Filepath
