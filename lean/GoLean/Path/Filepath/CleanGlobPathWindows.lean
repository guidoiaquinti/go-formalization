import GoLean.Path.Filepath.VolumeName

namespace GoLean.Path.Filepath

def cleanGlobPathWindows (windows : Bool) (path : List (BitVec 8)) : Nat × List (BitVec 8) :=
  let length := path.length
  let volume := (volumeName windows path).length
  let separator byte := byte == slash || (windows && byte == 92)
  if path.isEmpty then (0, [46])
  else if volume + 1 = length && path.getLast?.any separator then (volume + 1, path)
  else if volume = length && length = 2 then (volume, path ++ [46])
  else (min volume (length - 1), path.take (length - 1))

theorem cleanGlobPathWindows_empty (windows : Bool) :
    cleanGlobPathWindows windows [] = (0, [46]) := by
  cases windows <;> decide

theorem cleanGlobPathWindows_root : cleanGlobPathWindows true [67, 58, 47] = (3, [67, 58, 47]) := by
  decide

theorem cleanGlobPathWindows_drive : cleanGlobPathWindows true [67, 58] = (2, [67, 58, 46]) := by
  decide

theorem cleanGlobPathWindows_uncClamp :
    cleanGlobPathWindows true [47, 47, 104, 47, 115] = (4, [47, 47, 104, 47]) := by decide

end GoLean.Path.Filepath
