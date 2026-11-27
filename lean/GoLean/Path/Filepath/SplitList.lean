import Init.Data.BitVec.Basic

namespace GoLean.Path.Filepath

private def scanList (separator : BitVec 8) (quoted : Bool) (pending : List (BitVec 8)) :
    List (BitVec 8) → List (List (BitVec 8))
  | [] => [pending.reverse]
  | byte :: rest =>
    if separator == 59 && byte == 34 then scanList separator (!quoted) pending rest
    else if byte == separator && !quoted then
      pending.reverse :: scanList separator quoted [] rest
    else scanList separator quoted (byte :: pending) rest

def splitList (separator : BitVec 8) (path : List (BitVec 8)) : List (List (BitVec 8)) :=
  match path with
  | [] => []
  | _ :: _ => scanList separator false [] path

theorem splitList_empty (separator : BitVec 8) : splitList separator [] = [] := by
  rfl

theorem splitList_unixSeparator : splitList 58 [58] = [[], []] := by
  decide

theorem splitList_windowsSeparator : splitList 59 [59] = [[], []] := by
  decide

theorem splitList_plan9Separator : splitList 0 [0] = [[], []] := by
  decide

theorem splitList_windowsQuoted : splitList 59 [34, 59, 34, 59, 120] = [[59], [120]] := by
  decide

theorem splitList_windowsUnclosed : splitList 59 [34, 59, 120] = [[59, 120]] := by
  decide

end GoLean.Path.Filepath
