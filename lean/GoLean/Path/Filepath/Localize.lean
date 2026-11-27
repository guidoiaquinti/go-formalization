import GoLean.Path.Filepath.IsLocal
import GoLean.Unicode.Utf8.ValidString

namespace GoLean.Path.Filepath

private def slashParts (pending : List (BitVec 8)) :
    List (BitVec 8) → List (List (BitVec 8))
  | [] => [pending.reverse]
  | byte :: rest =>
    if byte == slash then pending.reverse :: slashParts [] rest
    else slashParts (byte :: pending) rest

def localizePath (profile : AbsProfile) (dosDevice : List (BitVec 8) → Bool)
    (path : List (BitVec 8)) : Option (List (BitVec 8)) :=
  let parts := slashParts [] path
  let validPart := fun (part : List (BitVec 8)) =>
    !part.isEmpty && part != [46] && part != [46, 46]
  let valid := GoLean.Unicode.Utf8.validString 64 (path.map BitVec.toNat) &&
    (path == [46] || parts.all validPart)
  if !valid || path.contains 0 then none
  else
    match profile with
    | .unix => some path
    | .plan9 => if path.head? == some 35 then none else some path
    | .windows =>
      if path.contains 58 || path.contains 92 ||
          !parts.all (isLocalPath .windows dosDevice) then none
      else some (fromSlash 92 path)

theorem localizePath_dot (profile : AbsProfile) :
    localizePath profile (fun _ => false) [46] = some [46] := by
  cases profile <;> decide

theorem localizePath_nonCanonical :
    localizePath .unix (fun _ => false) [97, 47, 46, 47, 98] = none := by decide

theorem localizePath_windowsSlash :
    localizePath .windows (fun _ => false) [97, 47, 98] = some [97, 92, 98] := by decide

theorem localizePath_plan9Device :
    localizePath .plan9 (fun _ => false) [35, 120] = none := by decide

end GoLean.Path.Filepath
