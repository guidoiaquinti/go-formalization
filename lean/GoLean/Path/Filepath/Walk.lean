import GoLean.Path.Filepath.WalkDirPublic

namespace GoLean.Path.Filepath

inductive WalkNode where
  | file (name : List (BitVec 8))
  | statFailed (name : List (BitVec 8)) (error : Nat)
  | directory (name : List (BitVec 8)) (listing : Except Nat (List WalkNode))

def WalkNode.name : WalkNode → List (BitVec 8)
  | .file name | .statFailed name _ | .directory name _ => name

private def WalkNode.swallowsSkipDir : WalkNode → Bool
  | .file _ => false
  | _ => true

abbrev WalkResult := List WalkDirPublicVisit × Option WalkDirError

private def visitWalk (windows : Bool) (callback : WalkDirPublicVisit → Option WalkDirError) :
    Nat → List (BitVec 8) → WalkNode → WalkResult
  | 0, _, _ => ([], some (.other 0))
  | fuel + 1, path, node =>
    match node with
    | .file _ =>
      let event : WalkDirPublicVisit := ⟨path, some false, none⟩
      ([event], callback event)
    | .statFailed _ error =>
      let event : WalkDirPublicVisit := ⟨path, none, some error⟩
      ([event], callback event)
    | .directory _ (.error error) =>
      let event : WalkDirPublicVisit := ⟨path, some true, some error⟩
      ([event], callback event)
    | .directory _ (.ok children) =>
      let event : WalkDirPublicVisit := ⟨path, some true, none⟩
      let response := callback event
      if response.isSome then ([event], response)
      else
        children.foldl
          (fun (state : WalkResult) child =>
            if state.2.isSome then state
            else
              let (trace, error) := visitWalk windows callback fuel
                (joinPath windows [path, child.name]) child
              let error := if child.swallowsSkipDir && error == some .skipDir then none else error
              (state.1 ++ trace, error))
          ([event], none)

def walkTree (windows : Bool) (maximumDepth : Nat) (path : List (BitVec 8))
    (root : WalkNode) (callback : WalkDirPublicVisit → Option WalkDirError) : WalkResult :=
  visitWalk windows callback (maximumDepth + 1) path root

theorem walkTree_file (windows : Bool) (maximumDepth : Nat)
    (path name : List (BitVec 8)) (callback : WalkDirPublicVisit → Option WalkDirError) :
    walkTree windows maximumDepth path (.file name) callback =
      ([⟨path, some false, none⟩], callback ⟨path, some false, none⟩) := by
  simp [walkTree, visitWalk]

theorem walkTree_readFailure (windows : Bool) (maximumDepth : Nat)
    (path name : List (BitVec 8)) (callback : WalkDirPublicVisit → Option WalkDirError)
    (error : Nat) :
    walkTree windows maximumDepth path (.directory name (.error error)) callback =
      ([⟨path, some true, some error⟩], callback ⟨path, some true, some error⟩) := by
  simp [walkTree, visitWalk]

end GoLean.Path.Filepath
