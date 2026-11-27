import GoLean.Path.Filepath.VolumeName

namespace GoLean.Path.Filepath

private def separator (windows : Bool) (byte : BitVec 8) : Bool :=
  byte == slash || (windows && byte == 92)

def baseName (windows : Bool) (path : List (BitVec 8)) : List (BitVec 8) :=
  if path.isEmpty then [46]
  else
    let ends := path.reverse.dropWhile (separator windows)
    let trimmed := ends.reverse
    let body := trimmed.drop (volumeName windows trimmed).length
    let leaf := (body.reverse.takeWhile (fun byte => !(separator windows byte))).reverse
    if leaf.isEmpty then [if windows then 92 else slash] else leaf

theorem baseName_empty (windows : Bool) : baseName windows [] = [46] := rfl

theorem baseName_roots :
    baseName false [47, 47] = [47] ∧ baseName true [47, 92] = [92] := by decide

theorem baseName_windowsDrive :
    baseName true [67, 58] = [92] ∧ baseName true [67, 58, 120] = [120] := by decide

theorem baseName_uncParent :
    baseName true [47, 47, 104, 47, 46, 46, 47] = [46, 46] := by decide

end GoLean.Path.Filepath
