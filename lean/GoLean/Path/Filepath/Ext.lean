import GoLean.Path.Filepath.ToSlash

namespace GoLean.Path.Filepath

private def isSeparator (windows : Bool) (byte : BitVec 8) : Bool :=
  (byte == slash) || (windows && byte == 92)

private def scanExt (windows : Bool) (reversed suffix : List (BitVec 8)) : List (BitVec 8) :=
  match reversed with
  | [] => []
  | byte :: rest =>
    if isSeparator windows byte then []
    else if byte = 46 then byte :: suffix
    else scanExt windows rest (byte :: suffix)

def ext (windows : Bool) (path : List (BitVec 8)) : List (BitVec 8) :=
  scanExt windows path.reverse []

theorem ext_empty (windows : Bool) : ext windows [] = [] := rfl

theorem ext_trailingSlash (windows : Bool) (path : List (BitVec 8)) :
    ext windows (path ++ [slash]) = [] := by
  simp [ext, List.reverse_append, scanExt, isSeparator]

theorem ext_trailingDot (windows : Bool) (path : List (BitVec 8)) :
    ext windows (path ++ [46]) = [46] := by
  cases windows <;> simp [ext, List.reverse_append, scanExt, isSeparator, slash]

theorem ext_trailingWindowsSeparator (path : List (BitVec 8)) :
    ext true (path ++ [92]) = [] := by
  simp [ext, List.reverse_append, scanExt, isSeparator]

end GoLean.Path.Filepath
