import GoLean.Path.Filepath.VolumeName
import Init.Data.List.Lemmas

namespace GoLean.Path.Filepath

private def separator (windows : Bool) (byte : BitVec 8) : Bool :=
  byte == slash || (windows && byte == 92)

def splitIndex (windows : Bool) (path : List (BitVec 8)) : Nat :=
  let volume := (volumeName windows path).length
  let tail := (path.reverse.takeWhile fun byte => !separator windows byte).length
  max volume (path.length - tail)

def splitPath (windows : Bool) (path : List (BitVec 8)) : List (BitVec 8) × List (BitVec 8) :=
  let index := splitIndex windows path
  (path.take index, path.drop index)

theorem splitPath_join (windows : Bool) (path : List (BitVec 8)) :
    (splitPath windows path).1 ++ (splitPath windows path).2 = path := by
  simp [splitPath]

theorem splitPath_empty (windows : Bool) : splitPath windows [] = ([], []) := by
  cases windows <;> decide

theorem splitPath_drive :
    splitPath true [67, 58, 120] = ([67, 58], [120]) ∧
    splitPath false [67, 58, 120] = ([], [67, 58, 120]) := by decide

theorem splitPath_uncRoot :
    splitPath true [47, 47, 104, 47, 115] = ([47, 47, 104, 47, 115], []) := by decide

end GoLean.Path.Filepath
