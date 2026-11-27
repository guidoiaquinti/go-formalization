import GoLean.Path.Filepath.Join
import GoLean.Path.Filepath.Match
import Init.Data.List.Sort

namespace GoLean.Path.Filepath

inductive GlobDirectory where
  | statFailed
  | notDirectory
  | openFailed
  | opened (names : List (List (BitVec 8)))

private def globNames (windows : Bool) (directory pattern : List (BitVec 8)) :
    List (List (BitVec 8)) → Option (List (List (BitVec 8))) →
      Option (List (List (BitVec 8))) × Bool
  | [], collected => (collected, false)
  | name :: rest, collected =>
    let (matched, invalid) := matchPath windows pattern name
    if invalid then (collected, true)
    else
      let next :=
        if matched then some (collected.getD [] ++ [joinPath windows [directory, name]])
        else collected
      globNames windows directory pattern rest next

def globDirectory (windows : Bool) (directory pattern : List (BitVec 8))
    (collected : Option (List (List (BitVec 8)))) (observed : GlobDirectory) :
    Option (List (List (BitVec 8))) × Bool :=
  match observed with
  | .opened names =>
    let ordered := names.mergeSort fun a b => !b.lex a (fun x y => decide (x.toNat < y.toNat))
    globNames windows directory pattern ordered collected
  | _ => (collected, false)

theorem globDirectory_statFailed (windows : Bool) (directory pattern : List (BitVec 8))
    (collected : Option (List (List (BitVec 8)))) :
    globDirectory windows directory pattern collected .statFailed = (collected, false) := rfl

theorem globDirectory_openFailed (windows : Bool) (directory pattern : List (BitVec 8))
    (collected : Option (List (List (BitVec 8)))) :
    globDirectory windows directory pattern collected .openFailed = (collected, false) := rfl

theorem globDirectory_emptyDirectory (windows : Bool) (directory pattern : List (BitVec 8))
    (collected : Option (List (List (BitVec 8)))) :
    globDirectory windows directory pattern collected (.opened []) = (collected, false) := by simp [globDirectory, globNames]

theorem globDirectory_badPattern :
    globDirectory false [97] [91] none (.opened [[98]]) = (none, true) := by
  simp only [globDirectory, List.mergeSort_singleton]
  decide

end GoLean.Path.Filepath
