import GoLean.Path.Filepath.VolumeName

namespace GoLean.Path.Filepath

private def isSeparator (windows : Bool) (byte : BitVec 8) : Bool :=
  byte == slash || (windows && byte == 92)

private def components (windows : Bool) (pending : List (BitVec 8)) :
    List (BitVec 8) → List (List (BitVec 8))
  | [] => [pending.reverse]
  | byte :: rest =>
    if isSeparator windows byte then pending.reverse :: components windows [] rest
    else components windows (byte :: pending) rest

private structure CleanState where
  stack : List (List (BitVec 8))
  backing : Option (List (BitVec 8))

private def render (separator : BitVec 8) (rooted : Bool) (stack : List (List (BitVec 8))) :
    List (BitVec 8) :=
  (if rooted then [separator] else []) ++ List.intercalate [separator] stack.reverse

private def writeBacking (initial : List (BitVec 8)) :
    Option (List (BitVec 8)) → List (BitVec 8) → Option (List (BitVec 8))
  | some bytes, written => some (written ++ bytes.drop written.length)
  | none, written =>
    if initial.take written.length == written then none
    else some (written ++ List.replicate (initial.length - written.length) 0)

private def addComponent (separator : BitVec 8) (rooted : Bool)
    (initial : List (BitVec 8)) (state : CleanState) (part : List (BitVec 8)) : CleanState :=
  let stack := part :: state.stack
  let written := render separator rooted stack
  ⟨stack, writeBacking initial state.backing written⟩

private def cleanStep (separator : BitVec 8) (rooted : Bool)
    (initial : List (BitVec 8)) (state : CleanState) (part : List (BitVec 8)) : CleanState :=
  if part.isEmpty || part == [46] then state
  else if part == [46, 46] then
    match state.stack with
    | top :: rest =>
      if top == [46, 46] then addComponent separator rooted initial state part
      else { state with stack := rest }
    | [] => if rooted then state else addComponent separator rooted initial state part
  else addComponent separator rooted initial state part

private def protectRelative (backing path : List (BitVec 8)) : List (BitVec 8) :=
  if (backing.takeWhile fun byte => !isSeparator true byte).contains 58 then [46, 92] ++ path
  else if (backing.head?.any (isSeparator true)) && backing[1]? == some 63 && backing[2]? == some 63 then
    [92, 46] ++ path
  else path

def cleanPath (windows : Bool) (path : List (BitVec 8)) : List (BitVec 8) :=
  let separator := if windows then 92 else slash
  let volume := (volumeName windows path).length
  let body := path.drop volume
  if body.isEmpty then
    if volume > 1 && path.head?.any (isSeparator windows) &&
        (path[1]?).any (isSeparator windows) then fromSlash separator path
    else path ++ [46]
  else
    let rooted := body.head?.any (isSeparator windows)
    let initial := if rooted then [separator] else []
    let backing := writeBacking body none initial
    let state := (components windows [] body).foldl (cleanStep separator rooted body) ⟨[], backing⟩
    let written := render separator rooted state.stack
    let final := if written.isEmpty then [46] else written
    let backing := if written.isEmpty then writeBacking body state.backing final else state.backing
    let safe :=
      if windows && volume == 0 then
        match backing with
        | none => final
        | some storage => protectRelative storage final
      else final
    fromSlash separator (path.take volume ++ safe)

theorem cleanPath_empty (windows : Bool) : cleanPath windows [] = [46] := by
  cases windows <;> decide

theorem cleanPath_rootParent : cleanPath false [47, 46, 46] = [47] := by decide

theorem cleanPath_windowsColon :
    cleanPath true [97, 47, 46, 46, 47, 99, 58] = [46, 92, 99, 58] := by decide

theorem cleanPath_windowsDevice :
    cleanPath true [47, 97, 47, 46, 46, 47, 63, 63, 47, 97] =
      [92, 46, 92, 63, 63, 92, 97] := by decide

theorem cleanPath_windowsStaleColon :
    cleanPath true [46, 47, 102, 111, 111, 58, 98, 97, 114, 47, 46, 46] = [46, 92, 46] := by
  decide

theorem cleanPath_windowsStaleDevice :
    cleanPath true [92, 46, 92, 63, 63, 92, 46, 46] = [92, 46, 92] := by decide

end GoLean.Path.Filepath
