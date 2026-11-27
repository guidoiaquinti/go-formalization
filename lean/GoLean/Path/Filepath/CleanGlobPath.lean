import Init.Data.BitVec.Basic
import Init.Data.List.Lemmas

namespace GoLean.Path.Filepath

def cleanGlobPath (separator : BitVec 8) (path : List (BitVec 8)) : List (BitVec 8) :=
  match path with
  | [] => [46]
  | [byte] => if byte = separator then [byte] else []
  | _ :: _ :: _ => path.take (path.length - 1)

theorem cleanGlobPath_empty (separator : BitVec 8) :
    cleanGlobPath separator [] = [46] := rfl

theorem cleanGlobPath_root (separator : BitVec 8) :
    cleanGlobPath separator [separator] = [separator] := by
  simp [cleanGlobPath]

theorem cleanGlobPath_single (separator byte : BitVec 8) (other : byte ≠ separator) :
    cleanGlobPath separator [byte] = [] := by
  simp [cleanGlobPath, other]

theorem cleanGlobPath_many (separator first second : BitVec 8) (rest : List (BitVec 8)) :
    cleanGlobPath separator (first :: second :: rest) =
      (first :: second :: rest).take ((first :: second :: rest).length - 1) := rfl

end GoLean.Path.Filepath
