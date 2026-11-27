import Init.Data.BitVec.Basic
import Init.Data.List.Lemmas

namespace GoLean.Path.Filepath

def slash : BitVec 8 := 47

def toSlash (separator : BitVec 8) (path : List (BitVec 8)) : List (BitVec 8) :=
  path.map fun byte => if byte = separator then slash else byte

theorem toSlash_nativeSlash (path : List (BitVec 8)) : toSlash slash path = path := by
  induction path with
  | nil => rfl
  | cons byte rest ih =>
    simp only [toSlash, List.map_cons] at ih ⊢
    rw [ih]
    split <;> simp_all

theorem toSlash_cons (separator byte : BitVec 8) (path : List (BitVec 8)) :
    toSlash separator (byte :: path) =
      (if byte = separator then slash else byte) :: toSlash separator path := rfl

theorem toSlash_length (separator : BitVec 8) (path : List (BitVec 8)) :
    (toSlash separator path).length = path.length := by
  simp [toSlash]

end GoLean.Path.Filepath
