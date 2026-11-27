import GoLean.Path.Filepath.Clean
import GoLean.Path.Filepath.IsAbs

namespace GoLean.Path.Filepath

private def localPathParts (isSeparator : BitVec 8 → Bool) (pending : List (BitVec 8)) :
    List (BitVec 8) → List (List (BitVec 8))
  | [] => [pending.reverse]
  | byte :: rest =>
    if isSeparator byte then pending.reverse :: localPathParts isSeparator [] rest
    else localPathParts isSeparator (byte :: pending) rest

private def upperPathByte (byte : BitVec 8) : BitVec 8 :=
  if byte.toNat ≥ 97 && byte.toNat ≤ 122 then byte - 32 else byte

private def reservedBase (name : List (BitVec 8)) : Bool :=
  let name := name.map upperPathByte
  [[67, 79, 78], [80, 82, 78], [65, 85, 88], [78, 85, 76],
   [67, 79, 78, 73, 78, 36], [67, 79, 78, 79, 85, 84, 36]].contains name ||
  ((name.take 3 == [67, 79, 77] || name.take 3 == [76, 80, 84]) &&
    match name.drop 3 with
    | [digit] => digit.toNat ≥ 49 && digit.toNat ≤ 57
    | [194, superscript] => [178, 179, 185].contains superscript
    | _ => false)

private def reservedPathPart (dosDevice : List (BitVec 8) → Bool) (part : List (BitVec 8)) : Bool :=
  let base := part.takeWhile fun byte => byte != 46 && byte != 58
  let base := (base.reverse.dropWhile (· == 32)).reverse
  reservedBase base && (base.length == part.length || (!part.contains 0 && dosDevice part))

private def startsUp (separator : BitVec 8) (path : List (BitVec 8)) : Bool :=
  path == [46, 46] || path.take 3 == [46, 46, separator]

def isLocalPath (profile : AbsProfile) (dosDevice : List (BitVec 8) → Bool)
    (path : List (BitVec 8)) : Bool :=
  match profile with
  | .windows =>
    let isSeparator := fun (byte : BitVec 8) => byte == slash || byte == 92
    let parts := localPathParts isSeparator [] path
    let dots := parts.any fun part => part == [46] || part == [46, 46]
    !path.isEmpty && !path.head?.any isSeparator && !path.contains 58 &&
      !parts.any (reservedPathPart dosDevice) &&
      !startsUp 92 (if dots then cleanPath true path else path)
  | profile =>
    let parts := localPathParts (· == slash) [] path
    let dots := parts.any fun part => part == [46] || part == [46, 46]
    !path.isEmpty && !isAbs profile path &&
      !startsUp slash (if dots then cleanPath false path else path)

theorem isLocalPath_plan9Device : isLocalPath .plan9 (fun _ => false) [35, 120] = false := by
  decide

theorem isLocalPath_windowsDevice : isLocalPath .windows (fun _ => false) [99, 111, 110] = false :=
  by decide

theorem isLocalPath_windowsDeviceExtension :
    isLocalPath .windows (fun _ => false) [99, 111, 110, 46, 120] = true ∧
    isLocalPath .windows (fun _ => true) [99, 111, 110, 46, 120] = false := by decide

end GoLean.Path.Filepath
