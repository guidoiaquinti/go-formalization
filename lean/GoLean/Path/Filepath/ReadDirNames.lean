import Init.Data.List.Sort

namespace GoLean.Path.Filepath

inductive DirNameOperation where
  | openFile
  | readAllNames
  | closeFile
  deriving DecidableEq

inductive DirNameRead (ε : Type) where
  | openError (error : ε)
  | opened (names : List (List UInt8)) (readError closeError : Option ε)

def readDirNames (observed : DirNameRead ε) :
    List DirNameOperation × Except ε (List (List UInt8)) :=
  match observed with
  | .openError error => ([.openFile], .error error)
  | .opened names error _ =>
    let trace := [DirNameOperation.openFile, .readAllNames, .closeFile]
    match error with
    | some error => (trace, .error error)
    | none => (trace, .ok (names.mergeSort (fun a b => !(b.lex a))))

theorem readDirNames_openError (error : ε) :
    readDirNames (.openError error) = ([.openFile], .error error) := rfl

theorem readDirNames_readError (names : List (List UInt8)) (error : ε) (close : Option ε) :
    readDirNames (.opened names (some error) close) =
      ([.openFile, .readAllNames, .closeFile], .error error) := rfl

theorem readDirNames_ignoresCloseError (names : List (List UInt8)) (close : ε) :
    readDirNames (.opened names none (some close)) = readDirNames (.opened names none none) := rfl

theorem readDirNames_preservesNames (names : List (List UInt8)) :
    (match (readDirNames (ε := Unit) (.opened names none none)).2 with
     | .ok sorted => sorted
     | .error _ => []).Perm names := by
  exact List.mergeSort_perm _ _

end GoLean.Path.Filepath
