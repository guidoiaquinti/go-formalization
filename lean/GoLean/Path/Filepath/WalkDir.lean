import GoLean.Path.Filepath.Join

namespace GoLean.Path.Filepath

inductive WalkDirError where
  | skipDir
  | skipAll
  | other (code : Nat)
  deriving DecidableEq, BEq

inductive WalkDirNode where
  | entry (name : List (BitVec 8))
  | directory (name : List (BitVec 8)) (readError : Option Nat) (children : List WalkDirNode)

def WalkDirNode.name : WalkDirNode → List (BitVec 8)
  | .entry name | .directory name _ _ => name

structure WalkDirVisit where
  path : List (BitVec 8)
  directory : Bool
  readError : Option Nat

abbrev WalkDirResult := List WalkDirVisit × Option WalkDirError

private def visitTree (windows : Bool) (visit : WalkDirVisit → Option WalkDirError) :
    Nat → List (BitVec 8) → WalkDirNode → WalkDirResult
  | 0, _, _ => ([], some (.other 0))
  | fuel + 1, path, node =>
    let directory := match node with | .entry _ => false | .directory .. => true
    let first : WalkDirVisit := ⟨path, directory, none⟩
    let response := visit first
    match node with
    | .entry _ => ([first], response)
    | .directory _ readError children =>
      let (visits, answer) :=
        if response.isSome then ([first], response)
        else
          match readError with
          | none => ([first], none)
          | some code =>
            let second : WalkDirVisit := ⟨path, true, some code⟩
            ([first, second], visit second)
      match answer with
      | some .skipDir => (visits, none)
      | some error => (visits, some error)
      | none =>
        let (accumulated, error, _) := children.foldl
          (fun (state : List WalkDirVisit × Option WalkDirError × Bool) child =>
            if state.2.2 then state
            else
              let fullPath := joinPath windows [path, child.name]
              let (childVisits, result) := visitTree windows visit fuel fullPath child
              let visited := state.1 ++ childVisits
              match result with
              | some .skipDir => (visited, none, true)
              | some error => (visited, some error, true)
              | none => (visited, none, false))
          (visits, none, false)
        (accumulated, error)

def walkDirTree (windows : Bool) (maximumDepth : Nat) (path : List (BitVec 8)) (root : WalkDirNode)
    (visit : WalkDirVisit → Option WalkDirError) : WalkDirResult :=
  visitTree windows visit (maximumDepth + 1) path root

theorem walkDirTree_entry (windows : Bool) (maximumDepth : Nat) (path name : List (BitVec 8))
    (visit : WalkDirVisit → Option WalkDirError) :
    walkDirTree windows maximumDepth path (.entry name) visit =
      ([⟨path, false, none⟩], visit ⟨path, false, none⟩) := by
  simp [walkDirTree, visitTree]

end GoLean.Path.Filepath
