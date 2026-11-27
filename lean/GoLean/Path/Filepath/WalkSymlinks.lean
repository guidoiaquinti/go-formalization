import GoLean.Path.Filepath.Clean
import GoLean.Path.Filepath.IsAbs

namespace GoLean.Path.Filepath

inductive SymlinkEntry where
  | file
  | directory
  | symlink

inductive WalkSymlinkError (ε : Type) where
  | observed (error : ε)
  | notDirectory
  | tooManyLinks
  | exhaustedModel

private abbrev PathBytes := List (BitVec 8)

private def symlinkSeparator (windows : Bool) (byte : BitVec 8) : Bool :=
  byte == 47 || (windows && byte == 92)

private def volumeLength (windows : Bool) (path : PathBytes) : Nat :=
  (volumeName windows path).length

private def initialVolumeLength (windows : Bool) (path : PathBytes) : Nat :=
  let length := volumeLength windows path
  if path[length]?.any (symlinkSeparator windows) then length + 1 else length

private def lastSeparator (windows : Bool) (start : Nat) (path : PathBytes) : Option Nat :=
  ((path.zipIdx 0).filter fun pair => start ≤ pair.2 && symlinkSeparator windows pair.1).getLast?.map Prod.snd

private def addPathPart (windows : Bool) (dest part : PathBytes) : PathBytes :=
  let separator : BitVec 8 := if windows then 92 else 47
  if volumeLength windows dest < dest.length && !dest.getLast?.any (symlinkSeparator windows)
  then dest ++ (separator :: part)
  else dest ++ part

private structure SymlinkCursor where
  input : PathBytes
  offset : Nat
  volume : PathBytes
  volumeLength : Nat
  destination : PathBytes
  links : Nat

private def applyLink (windows : Bool) (cursor : SymlinkCursor) (link : PathBytes) :
    SymlinkCursor :=
  let input := link ++ cursor.input.drop cursor.offset
  let length := volumeLength windows link
  if length > 0 then
    let volume := link.take (initialVolumeLength windows link)
    { cursor with input, offset := volume.length, volume, destination := volume }
  else if link.head?.any (symlinkSeparator windows) then
    let volume := link.take 1
    { cursor with input, offset := 1, volume, volumeLength := 1, destination := volume }
  else
    let dest := match lastSeparator windows cursor.volumeLength cursor.destination with
      | none => cursor.volume
      | some offset => cursor.destination.take offset
    { cursor with input, offset := 0, destination := dest }

private def walkSymlinkCursor (windows : Bool)
    (stat : PathBytes → Except ε SymlinkEntry) (readLink : PathBytes → Except ε PathBytes) :
    Nat → SymlinkCursor → Except (WalkSymlinkError ε) PathBytes
  | 0, _ => .error .exhaustedModel
  | fuel + 1, cursor => do
    let isSep := symlinkSeparator windows
    let tail := cursor.input.drop cursor.offset
    let remaining := tail.dropWhile isSep
    let part := remaining.takeWhile (!isSep ·)
    if part.isEmpty then return cleanPath windows cursor.destination
    let finish := cursor.offset + (tail.length - remaining.length) + part.length
    let next := { cursor with offset := finish }
    let windowsDot := windows && cursor.input.drop (volumeLength windows cursor.input) == [46]
    if part == [46] && !windowsDot then
      walkSymlinkCursor windows stat readLink fuel next
    else if part == [46, 46] then
      let dest := cursor.destination
      let modified := match lastSeparator windows cursor.volumeLength dest with
        | some split =>
          if dest.drop (split + 1) == [46, 46] then addPathPart windows dest part else dest.take split
        | none =>
          if cursor.volumeLength < dest.length
          then dest ++ [if windows then 92 else 47, 46, 46]
          else dest ++ part
      walkSymlinkCursor windows stat readLink fuel { next with destination := modified }
    else
      let dest := addPathPart windows cursor.destination part
      let entry ← (stat dest).mapError WalkSymlinkError.observed
      let next := { next with destination := dest }
      match entry with
      | .file =>
        if finish < cursor.input.length then .error .notDirectory
        else walkSymlinkCursor windows stat readLink fuel next
      | .directory => walkSymlinkCursor windows stat readLink fuel next
      | .symlink =>
        if cursor.links ≥ 255 then .error .tooManyLinks
        else
          let link ← (readLink dest).mapError WalkSymlinkError.observed
          if windowsDot && !isAbs .windows link then return cleanPath windows dest
          walkSymlinkCursor windows stat readLink fuel
            (applyLink windows { next with links := cursor.links + 1 } link)

def walkSymlinks (windows : Bool) (stepBound : Nat) (path : PathBytes)
    (stat : PathBytes → Except ε SymlinkEntry) (readLink : PathBytes → Except ε PathBytes) :
    Except (WalkSymlinkError ε) PathBytes :=
  let length := initialVolumeLength windows path
  let volume := path.take length
  walkSymlinkCursor windows stat readLink (stepBound + 1) ⟨path, length, volume, length, volume, 0⟩

theorem walkSymlinks_unixEmpty (stepBound : Nat)
    (stat : PathBytes → Except ε SymlinkEntry) (readLink : PathBytes → Except ε PathBytes) :
    walkSymlinks false stepBound [] stat readLink = .ok [46] := rfl

end GoLean.Path.Filepath
