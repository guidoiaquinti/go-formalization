import GoLean.Path.Filepath.Walk

namespace GoLean.Path.Filepath

private def normalizeWalkError : Option WalkDirError → Option WalkDirError
  | some .skipDir | some .skipAll => none
  | error => error

def walkPublic (windows : Bool) (maximumDepth : Nat) (path : List (BitVec 8))
    (root : Except Nat WalkNode) (callback : WalkDirPublicVisit → Option WalkDirError) :
    WalkResult :=
  match root with
  | .error failure =>
    let event : WalkDirPublicVisit := ⟨path, none, some failure⟩
    ([event], normalizeWalkError (callback event))
  | .ok root =>
    let (trace, error) := walkTree windows maximumDepth path root callback
    (trace, normalizeWalkError error)

theorem walkPublic_missingSkipDir (windows : Bool) (maximumDepth : Nat)
    (path : List (BitVec 8)) (failure : Nat) :
    walkPublic windows maximumDepth path (.error failure) (fun _ => some .skipDir) =
      ([⟨path, none, some failure⟩], none) := rfl

theorem walkPublic_missingError (windows : Bool) (maximumDepth : Nat)
    (path : List (BitVec 8)) (failure response : Nat) :
    walkPublic windows maximumDepth path (.error failure) (fun _ => some (.other response)) =
      ([⟨path, none, some failure⟩], some (.other response)) := rfl

end GoLean.Path.Filepath
