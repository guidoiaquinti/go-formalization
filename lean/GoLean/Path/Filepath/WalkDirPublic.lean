import GoLean.Path.Filepath.WalkDir

namespace GoLean.Path.Filepath

structure WalkDirPublicVisit where
  path : List (BitVec 8)
  entryIsDirectory : Option Bool
  readError : Option Nat

private def walkDirPublicError : Option WalkDirError → Option WalkDirError
  | some .skipDir | some .skipAll => none
  | result => result

def walkDirPublic (windows : Bool) (maximumDepth : Nat) (path : List (BitVec 8))
    (root : Except Nat WalkDirNode) (visit : WalkDirPublicVisit → Option WalkDirError) :
    List WalkDirPublicVisit × Option WalkDirError :=
  match root with
  | .error error =>
    let event : WalkDirPublicVisit := ⟨path, none, some error⟩
    ([event], walkDirPublicError (visit event))
  | .ok root =>
    let withEntry := fun event : WalkDirVisit =>
      WalkDirPublicVisit.mk event.path (some event.directory) event.readError
    let (trace, error) := walkDirTree windows maximumDepth path root (visit ∘ withEntry)
    (trace.map withEntry, walkDirPublicError error)

theorem walkDirPublic_missingSkipAll (windows : Bool) (maximumDepth : Nat)
    (path : List (BitVec 8)) (error : Nat) :
    walkDirPublic windows maximumDepth path (.error error) (fun _ => some .skipAll) =
      ([⟨path, none, some error⟩], none) := rfl

theorem walkDirPublic_missingOther (windows : Bool) (maximumDepth : Nat)
    (path : List (BitVec 8)) (error response : Nat) :
    walkDirPublic windows maximumDepth path (.error error) (fun _ => some (.other response)) =
      ([⟨path, none, some error⟩], some (.other response)) := rfl

end GoLean.Path.Filepath
