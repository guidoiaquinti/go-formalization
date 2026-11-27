import GoLean.Path.Filepath.ToSlash

namespace GoLean.Path.Filepath

def fromSlash (separator : BitVec 8) (path : List (BitVec 8)) : List (BitVec 8) :=
  path.map fun byte => if byte = slash then separator else byte

theorem fromSlash_nativeSlash (path : List (BitVec 8)) : fromSlash slash path = path :=
  toSlash_nativeSlash path

theorem fromSlash_cons (separator byte : BitVec 8) (path : List (BitVec 8)) :
    fromSlash separator (byte :: path) =
      (if byte = slash then separator else byte) :: fromSlash separator path := rfl

theorem fromSlash_length (separator : BitVec 8) (path : List (BitVec 8)) :
    (fromSlash separator path).length = path.length := by
  simp [fromSlash]

end GoLean.Path.Filepath
