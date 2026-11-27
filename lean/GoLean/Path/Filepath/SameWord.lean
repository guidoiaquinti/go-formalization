import GoLean.Strings.EqualFold

namespace GoLean.Path.Filepath

def sameWord (windows : Bool) (left right : List Nat) : Bool :=
  if windows then GoLean.Strings.equalFold left right else left == right

theorem sameWord_unix (left right : List Nat) :
    sameWord false left right = (left == right) := rfl

theorem sameWord_windows (left right : List Nat) :
    sameWord true left right = GoLean.Strings.equalFold left right := rfl

theorem sameWord_asciiExample :
    sameWord true [65] [97] = true ∧ sameWord false [65] [97] = false := by decide

theorem sameWord_malformedExample :
    sameWord true [255] [239, 191, 189] = true ∧
    sameWord false [255] [239, 191, 189] = false := by decide

end GoLean.Path.Filepath
