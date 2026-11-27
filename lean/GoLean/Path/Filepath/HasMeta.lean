import Init.Data.BitVec.Basic
import Init.Data.List.Lemmas

namespace GoLean.Path.Filepath

def metaByte (windows : Bool) (byte : BitVec 8) : Bool :=
  byte == 42 || byte == 63 || byte == 91 || (!windows && byte == 92)

def hasMeta (windows : Bool) (path : List (BitVec 8)) : Bool :=
  path.any (metaByte windows)

theorem hasMeta_empty (windows : Bool) : hasMeta windows [] = false := rfl

theorem hasMeta_cons (windows : Bool) (byte : BitVec 8) (rest : List (BitVec 8)) :
    hasMeta windows (byte :: rest) = (metaByte windows byte || hasMeta windows rest) := rfl

theorem hasMeta_backslash : hasMeta true [92] = false ∧ hasMeta false [92] = true := by
  decide

theorem hasMeta_append (windows : Bool) (left right : List (BitVec 8)) :
    hasMeta windows (left ++ right) = (hasMeta windows left || hasMeta windows right) := by
  simp [hasMeta]

end GoLean.Path.Filepath
